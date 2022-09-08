package ofd

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/JohanVong/fiscal_v3_sample/internal/db"
	"github.com/JohanVong/fiscal_v3_sample/internal/metrics"
	"github.com/JohanVong/fiscal_v3_sample/internal/models"
	"github.com/JohanVong/fiscal_v3_sample/internal/services/calculations"
	"github.com/JohanVong/fiscal_v3_sample/internal/services/graylog"
	"github.com/JohanVong/fiscal_v3_sample/internal/transport/protobuf/ofd"
	"github.com/JohanVong/fiscal_v3_sample/internal/transport/protobuf/ofd/cpcr"
	"github.com/JohanVong/fiscal_v3_sample/internal/transport/protobuf/ofd/kkm_proto"
)

type Pool map[int]ofd.HandlerOFD

var OfdPool Pool

var mutex sync.Mutex

func Register(handler ofd.HandlerOFD, id int) {
	OfdPool[id] = handler
}

func RetrieveInstance(id int) ofd.HandlerOFD {
	mutex.Lock()
	defer mutex.Unlock()

	h, ok := OfdPool[id]
	if ok == true {
		return h
	}
	return nil
}

type OfdCommand uint32

const (
	OfdCommandSystem OfdCommand = iota
	OfdCommandInfo
	OfdCommandSale
	OfdCommandSaleReturn
	OfdCommandPurchase
	OfdCommandPurchaseReturn
	OfdCommandRollback
	OfdCommandCloseShift
	OfdCommandReportX
	OfdCommandReportZ
	OfdCommandDeposit
	OfdCommandWithdraw
)

type offlineData struct {
	document models.Document
	idKkm    int
}

type kkmState struct {
	data             chan offlineData
	offlineModeStart time.Time
	stop             chan bool
	idKkm            int
}

func newKkmState(start time.Time, idKkm int) kkmState {
	//defer recoverpanic.RecoverPanic()
	return kkmState{
		data:             make(chan offlineData),
		offlineModeStart: start,
		stop:             make(chan bool),
		idKkm:            idKkm,
	}
}

var offlineQueue map[int]kkmState
var excludeQueue map[int]bool

func ProcessOfflineQueueV2(interval time.Duration) {
	var (
		kkm_ids []int
		wg      sync.WaitGroup
	)

	for {
		db.Orm.Table("KKM").Where("OfflineQueue > ?", 0).Pluck("idKKM", &kkm_ids)

		if len(kkm_ids) == 0 {
			time.Sleep(interval)
		} else {
			for _, kkm_id := range kkm_ids {
				wg.Add(1)

				go processKkmQueueV2(kkm_id, &wg)
			}
			wg.Wait()
		}
	}
}

func processKkmQueueV2(kkm_id int, wg *sync.WaitGroup) {
	var (
		documents []models.Document
		kkm       models.Kkm
		err       error
	)

	defer wg.Done()

	err = db.Orm.Order("idDocuments asc").Preload("Shift").Where("Offline=? AND idKKM = ?", true,
		kkm_id).Find(&documents).Error
	if err != nil {
		return
	}
	if len(documents) == 0 {
		db.Orm.Model(&kkm).Where("idKKM=?", kkm_id).Updates(map[string]interface{}{"OfflineQueue": 0})
		return
	}

	systemPart := &cpcr.OfdSystemPart{OfflinePeriod: &cpcr.OfflinePeriod{
		BeginTime: documents[0].DateDocument,
		EndTime:   time.Now(),
	}}
	db.Orm.Where("idKKM=?", kkm_id).Find(&kkm)
	//err = OfdPool.PerformEndOfflineModeRequest(&kkm, offlineModeStart)
	//if err != nil {
	//	return
	//}
	for _, document := range documents {
		var txErr error
		if time.Now().Add(time.Hour * (-72)).After(document.DateDocument) {
			db.Orm.Model(&models.Kkm{}).Where("idKKM=?", kkm_id).Update("idStatusKKM", 3)
		}
		tx := db.Orm.Begin()
		defer func() {
			tx.RollbackUnlessCommitted()
		}()
		var positions []models.Position
		var header *cpcr.OfdHeader
		var fp *cpcr.TicketNumber

		switch document.IdTypedocument {
		case 1:
			txErr = tx.Preload("Section").Where("idDocuments=?", document.Id).Find(&positions).Error
			if txErr != nil {
				return
			}
			header, fp, err = OfdPool.PerformTicketRequest(OfdCommandSale, &kkm, &document, positions, systemPart)
		case 4:
			txErr = tx.Preload("Section").Where("idDocuments=?", document.Id).Find(&positions).Error
			if txErr != nil {
				return
			}
			header, fp, err = OfdPool.PerformTicketRequest(OfdCommandPurchase, &kkm, &document, positions, systemPart)
		case 5:
			txErr = tx.Preload("Section").Where("idDocuments=?", document.Id).Find(&positions).Error
			if txErr != nil {
				return
			}
			header, fp, err = OfdPool.PerformTicketRequest(OfdCommandSaleReturn, &kkm, &document, positions, systemPart)
		case 2:
			header, err = OfdPool.PerformMoneyPlacementRequest(OfdCommandDeposit, &kkm, &document, systemPart)
		case 3:
			header, err = OfdPool.PerformMoneyPlacementRequest(OfdCommandWithdraw, &kkm, &document, systemPart)
		case 7:
			header, _, err = OfdPool.PerformReportRequest(OfdCommandCloseShift, &kkm, &document, systemPart)
		case 8:
			header, _, err = OfdPool.PerformReportRequest(OfdCommandReportX, &kkm, &document, systemPart)
		}
		if txErr != nil {
			return
		}
		if header != nil {
			txErr = tx.Model(&kkm).Where("idKKM=?", kkm.Id).Updates(map[string]interface{}{"tokenCPCR": header.Token,
				"reqnumCPCR": header.ReqNum}).Error //, "OfflineQueue": gorm.Expr("OfflineQueue - ?", 1)})
			if txErr != nil {
				return
			}
		}
		if fp != nil {
			txErr = tx.Model(&models.Document{}).Where("idDocuments=?",
				document.Id).Updates(map[string]interface{}{"FiscalNumber": fp.Number}).Error
			if txErr != nil {
				return
			}
		}
		if err == nil {
			txErr = tx.Model(&kkm).Where("idKKM=?", kkm.Id).Updates(map[string]interface{}{"OfflineQueue": gorm.Expr(
				"OfflineQueue - ?", 1)}).Error
			if txErr != nil {
				return
			}
			txErr = tx.Model(&models.Document{}).Where("idDocuments=?",
				document.Id).Updates(map[string]interface{}{"Offline": 0}).Error
			if txErr != nil {
				return
			}
			txErr = tx.Model(&models.Document{}).Where("idDocuments=?",
				document.Id).Update(map[string]interface{}{"token": kkm.TokenCPCR, "reqNum": kkm.ReqNumCPCR}).Error
			if txErr != nil {
				return
			}
		}
		tx.Commit()
		if err != nil {
			//RepairToken(err, &kkm)
			return
		}
	}
}

func ProcessOfflineQueue(wg *sync.WaitGroup, interval time.Duration) {
	if offlineQueue == nil {
		offlineQueue = make(map[int]kkmState)
	}
	handlerIsRunning := false
	for {
		select {
		case <-time.Tick(interval):
			wg.Add(1)
			excludeQueue = make(map[int]bool)
			var documents []models.Document
			db.Orm.Order("idDocuments asc").Preload("Shift").Where("Offline=?", true).Find(&documents)
			if handlerIsRunning == false {
				handlerIsRunning = true
				for _, doc := range documents {

					if excludeQueue[doc.IdKkm] == true {
						continue
					}
					_, ok := offlineQueue[doc.IdKkm]
					if ok == false {
						offlineQueue[doc.IdKkm] = newKkmState(doc.DateDocument, doc.IdKkm)
						go processKkmQueue(offlineQueue[doc.IdKkm])
					}

					d := offlineData{
						document: doc,
						idKkm:    doc.IdKkm,
					}

					stop := <-offlineQueue[doc.IdKkm].stop

					if stop == false {
						offlineQueue[doc.IdKkm].data <- d
					} else {
						close(offlineQueue[doc.IdKkm].data)
						delete(offlineQueue, doc.IdKkm)
						excludeQueue[doc.IdKkm] = true
					}

				}
				handlerIsRunning = false
			}

			wg.Done()

		}
	}
}

func processKkmQueue(state kkmState) {
	var kkm models.Kkm
	errkkm := db.Orm.Model(&models.Kkm{}).Where("idKKM=?", state.idKkm).Find(&kkm).Error
	if errkkm != nil {
		state.stop <- true
		return
	}
	if time.Now().Add(time.Hour*(-72)).After(state.offlineModeStart) && kkm.IdStatusKkm == 2 {
		db.Orm.Model(&models.Kkm{}).Where("idKKM=?", state.idKkm).Update("idStatusKKM", 3)
	}

	//err := OfdPool.PerformEndOfflineModeRequest(&kkm, state.offlineModeStart)
	systemPart := &cpcr.OfdSystemPart{
		OfflinePeriod: &cpcr.OfflinePeriod{
			BeginTime: state.offlineModeStart,
			EndTime:   time.Now(),
		},
	}

	//if err != nil {
	//	state.stop <- true
	//	log.Println("state stop and return", kkm.Id)
	//	return
	//}
	state.stop <- false
	for {
		select {
		case d, open := <-state.data:
			if open == false {
				return
			}
			var positions []models.Position
			var err, txErr error
			var header *cpcr.OfdHeader
			var fp *cpcr.TicketNumber
			tx := db.Orm.Begin()
			txErr = tx.Model(&models.Kkm{}).Where("idKKM=?", state.idKkm).Find(&kkm).Error
			switch d.document.IdTypedocument {
			case 1:
				txErr = tx.Preload("Section").Where("idDocuments=?", d.document.Id).Find(&positions).Error
				if txErr != nil {
					tx.Rollback()
					state.stop <- true
					break
				}
				header, fp, err = OfdPool.PerformTicketRequest(OfdCommandSale, &kkm, &d.document, positions, systemPart)
			case 4:
				txErr = tx.Preload("Section").Where("idDocuments=?", d.document.Id).Find(&positions).Error
				if txErr != nil {
					tx.Rollback()
					state.stop <- true
					break
				}
				header, fp, err = OfdPool.PerformTicketRequest(OfdCommandPurchase, &kkm, &d.document, positions, systemPart)
			case 5:
				txErr = tx.Preload("Section").Where("idDocuments=?", d.document.Id).Find(&positions).Error
				if txErr != nil {
					tx.Rollback()
					state.stop <- true
					break
				}
				header, fp, err = OfdPool.PerformTicketRequest(OfdCommandSaleReturn, &kkm, &d.document, positions, systemPart)
			case 9:
				txErr = tx.Preload("Section").Where("idDocuments=?", d.document.Id).Find(&positions).Error
				if txErr != nil {
					tx.Rollback()
					state.stop <- true
					break
				}
				header, fp, err = OfdPool.PerformTicketRequest(OfdCommandPurchaseReturn, &kkm, &d.document, positions, systemPart)
			case 2:
				header, err = OfdPool.PerformMoneyPlacementRequest(OfdCommandDeposit, &kkm, &d.document, systemPart)
			case 3:
				header, err = OfdPool.PerformMoneyPlacementRequest(OfdCommandWithdraw, &kkm, &d.document, systemPart)
			case 7:
				header, _, err = OfdPool.PerformReportRequest(OfdCommandCloseShift, &kkm, &d.document, systemPart)
			case 8:
				header, _, err = OfdPool.PerformReportRequest(OfdCommandReportX, &kkm, &d.document, systemPart)
			}
			if txErr != nil {
				tx.Rollback()
				state.stop <- true
				break
			}
			if header != nil {
				txErr = tx.Model(&models.Kkm{}).Where("idKKM=?",
					kkm.Id).Updates(map[string]interface{}{"tokenCPCR": header.Token,
					"reqnumCPCR": header.ReqNum}).Error //, "OfflineQueue": gorm.Expr("OfflineQueue - ?", 1)})
				if txErr != nil {
					tx.Rollback()
					state.stop <- true
					break
				}
			}
			if fp != nil {
				txErr = tx.Model(&models.Document{}).Where("idDocuments=?",
					d.document.Id).Updates(map[string]interface{}{"FiscalNumber": fp.Number}).Error
				if txErr != nil {
					tx.Rollback()
					state.stop <- true
					break
				}
			}
			if err == nil {
				tx.Model(&models.Kkm{}).Where("idkkm=?", kkm.Id).Update("ofd_code", 0)
				txErr = tx.Model(&models.Kkm{}).Where("idKKM=?",
					kkm.Id).Updates(map[string]interface{}{"OfflineQueue": gorm.Expr("OfflineQueue - ?", 1)}).Error
				if txErr != nil {
					tx.Rollback()
					state.stop <- true
					break
				}
				txErr = tx.Model(&models.Document{}).Where("idDocuments=?",
					d.document.Id).Updates(map[string]interface{}{"Offline": 0}).Error
				if txErr != nil {
					tx.Rollback()
					state.stop <- true
					break
				}
				txErr = tx.Model(&models.Document{}).Where("idDocuments=?", d.document.Id).Update(map[string]interface{}{"token": kkm.TokenCPCR,
					"reqNum": kkm.ReqNumCPCR}).Error
				if txErr != nil {
					tx.Rollback()
					state.stop <- true
					break
				}
			} else {
				ofde, ok := err.(*cpcr.OfdError)
				if ok == true {
					tx.Model(&models.Kkm{}).Where("idkkm=?", kkm.Id).Update("ofd_code", ofde.Code)
				} else {
					tx.Model(&models.Kkm{}).Where("idkkm=?", kkm.Id).Update("ofd_code", 500)
				}
			}
			tx.Commit()
			// Новый функционал (начало)
			db.Orm.Model(&models.Kkm{}).Where("idKKM=?", state.idKkm).Find(&kkm)
			if kkm.OfflineQueue == 0 && kkm.IdStatusKkm == 3 {
				db.Orm.Model(&models.Kkm{}).Where("idKKM=?", state.idKkm).Update("idStatusKKM", 2)
			}
			// Новый функционал (конец)
			if err != nil {
				//RepairToken(err, &kkm)
				state.stop <- true
				break
			} else {
				state.stop <- false
			}
		}
	}
}

func (pool Pool) ProcessTicketRequest(action OfdCommand, kkm *models.Kkm, document *models.Document, items []models.Position, tx *gorm.DB) (header *cpcr.OfdHeader, ticketNumber *cpcr.TicketNumber, err error) {

	//defer recoverpanic.RecoverPanic()

	if kkm.OfflineQueue > 0 {
		if document.AutonomousNumber == 0 {
			ticketNumber, err = PerformOfflineTicketOperation(kkm, document, tx)
			return header, ticketNumber, err
		} else {
			return
		}
	}

	header, ticketNumber, err = pool.PerformTicketRequest(action, kkm, document, items, nil)

	if err != nil {
		ticketNumber, err = PerformOfflineTicketOperation(kkm, document, tx)
		return header, ticketNumber, err
	}

	return
}

func (pool Pool) ProcessMoneyPlacementRequest(action OfdCommand, kkm *models.Kkm, document *models.Document, tx *gorm.DB) (header *cpcr.OfdHeader, err error) {
	//defer recoverpanic.RecoverPanic()

	if kkm.OfflineQueue > 0 {
		if document.Offline == false {
			return header, PerformOfflineMoneyPlacementOperation(kkm, document, tx)
		} else {
			return
		}

	}

	header, err = pool.PerformMoneyPlacementRequest(action, kkm, document, nil)

	if err != nil {
		return header, PerformOfflineMoneyPlacementOperation(kkm, document, tx)
	}

	return
}

func (pool Pool) ProcessReportRequest(action OfdCommand, kkm *models.Kkm, document *models.Document, tx *gorm.DB) (response *kkm_proto.ReportResponse, header *cpcr.OfdHeader, err error) {
	//defer recoverpanic.RecoverPanic()

	if kkm.OfflineQueue > 0 {
		if document.Offline == false {
			return nil, header, PerformOfflineReportRequest(kkm, document, tx)
		} else {
			return
		}

	}

	header, response, err = pool.PerformReportRequest(action, kkm, document, nil)

	if err != nil {
		return nil, header, PerformOfflineReportRequest(kkm, document, tx)
	}
	return
}

func PerformOfflineTicketOperation(kkm *models.Kkm, document *models.Document, tx *gorm.DB) (ticketNumber *cpcr.TicketNumber, err error) {
	//defer recoverpanic.RecoverPanic()

	ticketNumber = new(cpcr.TicketNumber)

	ap := calculations.GetAfp(kkm.Id)
	ticketNumber.Number = uint64(ap)
	ticketNumber.Offline = true
	if document.AutonomousNumber == 0 {
		document.AutonomousNumber = ap
	}
	//tx := models.Orm.Begin()
	tx.Model(&models.Kkm{}).Where("idKKM=?", kkm.Id).Updates(map[string]interface{}{"OfflineQueue": gorm.Expr("OfflineQueue + ?", 1)})
	tx.Model(&models.Document{}).Where("idDocuments=?", document.Id).Updates(map[string]interface{}{"AutonomousNumber": document.AutonomousNumber})
	tx.Model(&models.Document{}).Where("idDocuments=?", document.Id).Updates(map[string]interface{}{"Offline": true})
	//tx.Commit()
	return ticketNumber, errors.New("операция проведена в автономном режиме")
}

func PerformOfflineMoneyPlacementOperation(kkm *models.Kkm, document *models.Document, tx *gorm.DB) (err error) {
	//tx := models.Orm.Begin()
	//defer recoverpanic.RecoverPanic()
	tx.Model(&models.Kkm{}).Where("idKKM=?", kkm.Id).Updates(map[string]interface{}{"OfflineQueue": gorm.Expr("OfflineQueue + ?", 1)})
	tx.Model(&models.Document{}).Where("idDocuments=?", document.Id).Updates(map[string]interface{}{"Offline": true})
	//tx.Commit()
	return errors.New("операция проведена в автономном режиме")
}

func PerformOfflineReportRequest(kkm *models.Kkm, document *models.Document, tx *gorm.DB) (err error) {
	//tx := models.Orm.Begin()
	//defer recoverpanic.RecoverPanic()

	tx.Model(&models.Kkm{}).Where("idKKM=?", kkm.Id).Updates(map[string]interface{}{"OfflineQueue": gorm.Expr("OfflineQueue + ?", 1)})
	tx.Model(&models.Document{}).Where("idDocuments=?", document.Id).Updates(map[string]interface{}{"Offline": true})
	//tx.Commit()
	//return errors.New("операция проведена в автономном режиме")
	return nil
}

func (pool Pool) PerformTicketRequest(action OfdCommand, kkm *models.Kkm, document *models.Document,
	items []models.Position, systemPart *cpcr.OfdSystemPart) (header *cpcr.OfdHeader, ticketNumber *cpcr.TicketNumber, err error) {
	//defer recoverpanic.RecoverPanic()
	var arg interface{}
	//var documentID int
	var ofdResult cpcr.OfdResult
	if document == nil || items == nil {
		return header, ticketNumber, errors.New(fmt.Sprintf("given command requires 2 non-nil argument of type %s and %s", reflect.TypeOf(&models.Document{}).String(), reflect.TypeOf([]models.Position{}).String()))
	}

	//documentID = document.Id
	ofdHeader := &cpcr.OfdHeader{
		Id:     kkm.IdCPCR,
		Token:  kkm.TokenCPCR,
		ReqNum: kkm.ReqNumCPCR,
	}
	handler := RetrieveInstance(kkm.IdOfd)

	if handler == nil {
		errofd := errors.New("попытка отправки даных в неинициализировнный ОФД")
		m := graylog.NewGELFMessage("Error")
		m.Kkm = kkm.Id
		*m.SendDate = time.Now()
		m.SendMessage = errofd.Error()

		return nil, nil, errofd
	}

	switch action {
	case OfdCommandSale:
		arg, err = handler.SendOperationSale(document, items, ofdHeader, systemPart)

	case OfdCommandSaleReturn:
		arg, err = handler.SendOperationSaleReturn(document, items, ofdHeader, systemPart)

	case OfdCommandPurchase:
		arg, err = handler.SendOperationPurchase(document, items, ofdHeader, systemPart)

	case OfdCommandPurchaseReturn:
		arg, err = handler.SendOperationPurchaseReturn(document, items, ofdHeader, systemPart)
	default:
		return header, ticketNumber, errors.New("unknown ofd command")
	}
	var status *models.KKMOFDStatus
	mutex.Lock()
	status, _ = models.KKMOFDStatusStorage[kkm.Id]
	mutex.Unlock()
	if status == nil {
		status = new(models.KKMOFDStatus)
		status.IDKKM = kkm.Id
		status.ZNM = kkm.Znm
	}
	status.Time = time.Now()
	if err != nil {
		ofde, ok := err.(*cpcr.OfdError)
		if ok == true {
			m := graylog.NewGELFMessage("Error")
			m.Kkm = kkm.Id
			m.Shift = document.IdShift
			*m.SendDate = time.Now()
			m.SendMessage = ofde.Error()
			_, _ = graylog.SendMessage(m)
			//errorMessage := fmt.Sprintf("<!channel> Ошибка ОФД при проведении чековой операции, ККМ: %d, "+
			//	"№ смены: %d", kkm.Id, document.IdShift)
			//if notifier.Notifier != nil {
			//	go notifier.Notifier.SendNotification(errorMessage)
			//}
			status.StatusOFD = ofde.Error()
			status.Code = ofde.Code
			metrics.SetKkmOfdStatusMetric(kkm.Id, int(ofde.Code))
		} else {
			status.StatusOFD = err.Error()
			status.Code = 500
			metrics.SetKkmOfdStatusMetric(kkm.Id, 500)
		}
		metrics.MetricMutex.Lock()
		metrics.OperationOFDSuccessMetric.ErroneousResponse++
		metrics.OperationOFDSuccessMetric.TotalSent++
		metrics.MetricMutex.Unlock()
		mutex.Lock()
		models.KKMOFDStatusStorage[kkm.Id] = status
		mutex.Unlock()
		return nil, ticketNumber, err
	}

	if arg != nil {
		ofdResult = arg.(cpcr.OfdResult)
		for _, entry := range ofdResult.Response {
			switch response := entry.(type) {
			case *kkm_proto.TicketResponse:
				if response != nil {
					if response.TicketNumber != nil {
						ticketNumber = new(cpcr.TicketNumber)
						fp, _ := strconv.ParseUint(*response.TicketNumber, 10, 64)
						ticketNumber.Number = fp
						ticketNumber.Offline = false
						ticketNumber.OFDQR = string(response.QrCode)
						//tx.Model(&models.Document{}).Where("idDocuments=?", documentID).Updates(map[string]interface{}{"FiscalNumber": fp})
						document.FiscalNumber = fp
					}
				} else {
					return nil, nil, errors.New("ответ офд не содержит фискльный признак")
				}
			}
		}
		metrics.MetricMutex.Lock()
		metrics.OperationOFDSuccessMetric.SuccessfulResponse++
		metrics.OperationOFDSuccessMetric.TotalSent++
		metrics.MetricMutex.Unlock()
		//tx.Commit()
		header = ofdResult.Header
		status.StatusOFD = "OK"
		status.Code = 0
		metrics.SetKkmOfdStatusMetric(kkm.Id, 0)
		mutex.Lock()
		models.KKMOFDStatusStorage[kkm.Id] = status
		mutex.Unlock()
	}
	return
}

func (pool Pool) PerformMoneyPlacementRequest(action OfdCommand, kkm *models.Kkm, document *models.Document,
	systemPart *cpcr.OfdSystemPart) (header *cpcr.OfdHeader, err error) {
	//defer recoverpanic.RecoverPanic()
	var arg interface{}
	var ofdResult cpcr.OfdResult
	if document == nil {
		return header, errors.New(fmt.Sprintf("given command requires an argument of type %s ", reflect.TypeOf(&models.Document{}).String()))
	}

	ofdHeader := &cpcr.OfdHeader{
		Id:     kkm.IdCPCR,
		Token:  kkm.TokenCPCR,
		ReqNum: kkm.ReqNumCPCR,
	}
	handler := RetrieveInstance(kkm.IdOfd)
	if handler == nil {
		errofd := errors.New("попытка отправки даных в неинициализировнный ОФД")
		m := graylog.NewGELFMessage("Error")
		m.Kkm = kkm.Id
		*m.SendDate = time.Now()
		m.SendMessage = errofd.Error()

		return nil, errofd
	}

	switch action {
	case OfdCommandDeposit:
		arg, err = handler.DepositMoney(document, ofdHeader, systemPart)

	case OfdCommandWithdraw:
		arg, err = handler.WithdrawMoney(document, ofdHeader, systemPart)

	default:
		return header, errors.New("unknown ofd command")
	}
	var status *models.KKMOFDStatus
	mutex.Lock()
	status, _ = models.KKMOFDStatusStorage[kkm.Id]
	mutex.Unlock()
	if status == nil {
		status = new(models.KKMOFDStatus)
		status.IDKKM = kkm.Id
		status.ZNM = kkm.Znm
	}
	status.Time = time.Now()
	if err != nil {
		ofde, ok := err.(*cpcr.OfdError)
		if ok == true {
			m := graylog.NewGELFMessage("Error")
			m.Kkm = kkm.Id
			m.Shift = document.IdShift
			*m.SendDate = time.Now()
			m.SendMessage = ofde.Error()
			_, _ = graylog.SendMessage(m)
			//errorMessage := fmt.Sprintf("<!channel> Ошибка ОФД при проведении операции внечения/снятия денег, "+
			//	"ККМ: %d, № смены: %d", kkm.Id, document.IdShift)
			//if notifier.Notifier != nil {
			//	go notifier.Notifier.SendNotification(errorMessage)
			//}
			status.StatusOFD = ofde.Error()
			status.Code = ofde.Code
			metrics.SetKkmOfdStatusMetric(kkm.Id, int(ofde.Code))
		} else {
			status.StatusOFD = err.Error()
			status.Code = 500
			metrics.SetKkmOfdStatusMetric(kkm.Id, 500)
		}
		metrics.MetricMutex.Lock()
		metrics.OperationOFDSuccessMetric.ErroneousResponse++
		metrics.OperationOFDSuccessMetric.TotalSent++
		metrics.MetricMutex.Unlock()
		mutex.Lock()
		models.KKMOFDStatusStorage[kkm.Id] = status
		mutex.Unlock()
		return header, err
	}

	if arg != nil {
		ofdResult = arg.(cpcr.OfdResult)
		header = ofdResult.Header
		metrics.MetricMutex.Lock()
		metrics.OperationOFDSuccessMetric.SuccessfulResponse++
		metrics.OperationOFDSuccessMetric.TotalSent++
		metrics.MetricMutex.Unlock()
		status.StatusOFD = "OK"
		status.Code = 0
		metrics.SetKkmOfdStatusMetric(kkm.Id, 0)
		mutex.Lock()
		models.KKMOFDStatusStorage[kkm.Id] = status
		mutex.Unlock()
	}
	return
}

func (pool Pool) PerformReportRequest(action OfdCommand, kkm *models.Kkm, document *models.Document, systemPart *cpcr.OfdSystemPart) (header *cpcr.OfdHeader, response *kkm_proto.ReportResponse, err error) {
	//defer recoverpanic.RecoverPanic()
	var arg interface{}
	var ofdResult cpcr.OfdResult
	if document == nil {
		return header, nil, errors.New(fmt.Sprintf("given command requires an argument of type %s ", reflect.TypeOf(&models.Document{}).String()))
	}
	ofdHeader := &cpcr.OfdHeader{
		Id:     kkm.IdCPCR,
		Token:  kkm.TokenCPCR,
		ReqNum: kkm.ReqNumCPCR,
	}
	handler := RetrieveInstance(kkm.IdOfd)

	if handler == nil {
		errofd := errors.New("попытка отправки даных в неинициализировнный ОФД")
		m := graylog.NewGELFMessage("Error")
		m.Kkm = kkm.Id
		*m.SendDate = time.Now()
		m.SendMessage = errofd.Error()

		return nil, nil, errofd
	}

	switch action {
	case OfdCommandReportX:
		arg, err = handler.RequestXReport(document, ofdHeader, systemPart)

	case OfdCommandReportZ:
		arg, err = handler.RequestZReport(document, ofdHeader, systemPart)

	case OfdCommandCloseShift:
		arg, err = handler.CloseShift(document, ofdHeader, systemPart)

	default:
		return header, nil, errors.New("unknown ofd command")
	}
	var status *models.KKMOFDStatus
	mutex.Lock()
	status, _ = models.KKMOFDStatusStorage[kkm.Id]
	mutex.Unlock()
	if status == nil {
		status = new(models.KKMOFDStatus)
		status.IDKKM = kkm.Id
		status.ZNM = kkm.Znm
	}
	status.Time = time.Now()
	if err != nil {
		ofde, ok := err.(*cpcr.OfdError)
		if ok == true {
			m := graylog.NewGELFMessage("Error")
			m.Kkm = kkm.Id
			m.Shift = document.IdShift
			*m.SendDate = time.Now()
			m.SendMessage = ofde.Error()
			_, _ = graylog.SendMessage(m)
			//errorMessage := fmt.Sprintf("<!channel> Ошибка ОФД при запросе отчета, ККМ: %d, № смены: %d", kkm.Id,
			//	document.IdShift)
			//if notifier.Notifier != nil {
			//	go notifier.Notifier.SendNotification(errorMessage)
			//}
			status.StatusOFD = ofde.Error()
			status.Code = ofde.Code
			metrics.SetKkmOfdStatusMetric(kkm.Id, int(ofde.Code))
		} else {
			status.StatusOFD = err.Error()
			status.Code = 500
			metrics.SetKkmOfdStatusMetric(kkm.Id, 500)
		}
		metrics.MetricMutex.Lock()
		metrics.OperationOFDSuccessMetric.ErroneousResponse++
		metrics.OperationOFDSuccessMetric.TotalSent++
		metrics.MetricMutex.Unlock()
		mutex.Lock()
		models.KKMOFDStatusStorage[kkm.Id] = status
		mutex.Unlock()
		return header, nil, err
	}

	if arg != nil {
		ofdResult = arg.(cpcr.OfdResult)
		for _, entry := range ofdResult.Response {
			switch resp := entry.(type) {
			case *kkm_proto.ReportResponse:
				response = resp
			}
		}
		header = ofdResult.Header
		metrics.MetricMutex.Lock()
		metrics.OperationOFDSuccessMetric.SuccessfulResponse++
		metrics.OperationOFDSuccessMetric.TotalSent++
		metrics.MetricMutex.Unlock()
		status.StatusOFD = "OK"
		status.Code = 0
		metrics.SetKkmOfdStatusMetric(kkm.Id, 0)
		mutex.Lock()
		models.KKMOFDStatusStorage[kkm.Id] = status
		mutex.Unlock()
	}
	return
}

func (pool Pool) PerformEndOfflineModeRequest(kkm *models.Kkm, endTime time.Time) (err error) {
	//defer recoverpanic.RecoverPanic()
	ofdHeader := &cpcr.OfdHeader{
		Id:     kkm.IdCPCR,
		Token:  kkm.TokenCPCR,
		ReqNum: kkm.ReqNumCPCR,
	}
	var status *models.KKMOFDStatus
	mutex.Lock()
	status, _ = models.KKMOFDStatusStorage[kkm.Id]
	mutex.Unlock()
	if status == nil {
		status = new(models.KKMOFDStatus)
		status.IDKKM = kkm.Id
		status.ZNM = kkm.Znm
	}
	status.Time = time.Now()
	handler := RetrieveInstance(kkm.IdOfd)
	if handler == nil {
		errofd := errors.New("попытка отправки даных в неинициализировнный ОФД")
		m := graylog.NewGELFMessage("Error")
		m.Kkm = kkm.Id
		*m.SendDate = time.Now()
		m.SendMessage = errofd.Error()
		return errofd
	}
	_, err = handler.SendOfflineModeInformation(endTime, ofdHeader)
	if err != nil {
		ofde, ok := err.(*cpcr.OfdError)
		if ok == true {
			status.StatusOFD = ofde.Error()
			status.Code = ofde.Code
			metrics.SetKkmOfdStatusMetric(kkm.Id, int(ofde.Code))
		} else {
			status.StatusOFD = err.Error()
			status.Code = 500
			metrics.SetKkmOfdStatusMetric(kkm.Id, 500)
		}
	} else {
		status.StatusOFD = "OK"
		status.Code = 0
		metrics.SetKkmOfdStatusMetric(kkm.Id, 0)
		//if arg != nil {
		//	switch h := arg.(type) {
		//	case cpcr.OfdHeader:
		//		models.Orm.Model(&models.Kkm{}).Where("idKKM=?", kkm.Id).Updates(map[string]interface{}{"reqnumCPCR": h.ReqNum})
		//	case *cpcr.OfdHeader:
		//		models.Orm.Model(&models.Kkm{}).Where("idKKM=?", kkm.Id).Updates(map[string]interface{}{"reqnumCPCR": h.ReqNum})
		//	}
		//}
	}
	mutex.Lock()
	models.KKMOFDStatusStorage[kkm.Id] = status
	mutex.Unlock()
	return
}

func RepairToken(incErr error, kkm *models.Kkm) {
	if incErr == nil {
		return
	}
	if ofde, ok := incErr.(*cpcr.OfdError); ok == true {
		if ofde.Code != 2 {
			return
		}
	}

	lastSuccessfulDoc := models.Document{}
	tx := db.Orm.Begin()
	defer func() {
		tx.RollbackUnlessCommitted()
	}()
	err := tx.Where("idKKM=? AND isnull(FiscalNumber, '')<>'' and isnull(token, 0)<>0", kkm.Id).Last(&lastSuccessfulDoc).Error
	if err != nil {
		return
	}

	var items []models.Position
	err = tx.Preload("Section").Where("idDocuments=?", lastSuccessfulDoc.Id).Find(&items).Error
	if err != nil {
		return
	}

	ofdHeader := cpcr.OfdHeader{
		Id:     kkm.IdCPCR,
		Token:  lastSuccessfulDoc.Token,
		ReqNum: lastSuccessfulDoc.ReqNum,
	}

	handler := RetrieveInstance(kkm.IdOfd)
	arg, err := handler.SendOperationSale(&lastSuccessfulDoc, items, &ofdHeader)
	if err != nil {
		return
	}

	if arg != nil {
		ofdResult := arg.(cpcr.OfdResult)
		err = tx.Model(&models.Kkm{}).Where("idKKM=?", kkm.Id).Updates(map[string]interface{}{"tokenCPCR": ofdResult.
			Header.Token,
			"reqnumCPCR": ofdResult.Header.ReqNum}).Error
		if err != nil {
			return
		}
		tx.Commit()

	}

}
