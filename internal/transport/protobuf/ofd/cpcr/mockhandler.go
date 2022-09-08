package cpcr

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/JohanVong/fiscal_v3_sample/internal/models"
	"github.com/JohanVong/fiscal_v3_sample/internal/services/calculations"
	"github.com/JohanVong/fiscal_v3_sample/internal/transport/protobuf/ofd/kkm_proto"
)

const (
	OperationSale = iota
	OperationPurchase
	OperationReturn
	OperationDeposit
	OperationWithdrawal
	OperationCloseShift
	OperationReportX
	OperationPurchaseReturn
)

type MockOFD struct {
	Down         bool
	NoConnection bool
	MinDelay     time.Duration
	MaxDelay     time.Duration
}

var OFDStorage *MockOFDStorage

type MockOFDStorage struct {
	KKMs map[uint32]KKM
}

func (o *MockOFDStorage) String() string {
	builder := strings.Builder{}
	for k, v := range o.KKMs {
		builder.WriteString(fmt.Sprintf("======================= КАССА №%d =======================\n", k))
		builder.WriteString(fmt.Sprintf("%s\n", v.String()))
	}

	return builder.String()
}

type KKM struct {
	ID     uint32
	Shifts []Shift
}

func (k *KKM) String() string {
	builder := strings.Builder{}
	for _, s := range k.Shifts {
		builder.WriteString(fmt.Sprintf("ОПЕРАЦИИ ПО СМЕНЕ №%d ----------------------------\n", s.ID))
		builder.WriteString(fmt.Sprintf("%s\n", s.String()))
	}

	return builder.String()
}

type Shift struct {
	ID         uint32
	IdShift    int
	Operations []OperationInterface
}

func (s *Shift) String() string {
	builder := strings.Builder{}
	for _, o := range s.Operations {
		builder.WriteString(fmt.Sprintf("%s\n", o.String()))
	}

	return builder.String()
}

type Operation struct {
	date    time.Time
	name    string
	opType  int
	idShift int
}

func (o *Operation) GetDate() time.Time {
	return o.date
}

func (o *Operation) GetName() string {
	return o.name
}

func (o Operation) GetType() int {
	return o.opType
}

type OperationInterface interface {
	String() string
	GetDate() time.Time
	GetName() string
	GetType() int
}

func NewOperation(opType int, date time.Time) (op *Operation) {
	op = new(Operation)
	op.opType = opType
	op.date = date
	switch opType {
	case 0:
		op.name = "Offline Mode Tag"
	case 1:
		op.name = "Sale"
	case 2:
		op.name = "Deposit Money"
	case 3:
		op.name = "Withdraw Money"
	case 4:
		op.name = "Purchase"
	case 5:
		op.name = "Sale Return"
	case 7:
		op.name = "Close Shift"
	case 8:
		op.name = "X Report"
	default:
		return nil
	}

	return
}

func (o *Operation) String() string {
	return fmt.Sprintf("| %d | %s | %s |", o.idShift, o.date.String(), o.name)

}

var mutex sync.Mutex

func (m *MockOFD) SendOfflineModeInformation(start time.Time, args ...interface{}) (interface{}, error) {
	header, err := findHeader(args)
	if err != nil {
		return nil, err
	}

	err = m.mockProcess()

	if err != nil {
		return nil, err
	}
	result := OfdResult{
		Header: header,
	}
	op := NewOperation(0, start)
	mutex.Lock()
	v, ok := OFDStorage.KKMs[header.Id]

	if ok == true {
		v.ID = header.Id
		v.Shifts[len(v.Shifts)-1].Operations = append(v.Shifts[len(v.Shifts)-1].Operations, op)
		//v.Shifts[len(v.Shifts)-1].ID = 1
	} else {
		OFDStorage.KKMs[header.Id] = KKM{
			ID:     header.Id,
			Shifts: make([]Shift, 1),
		}

		OFDStorage.KKMs[header.Id].Shifts[0].Operations = append(OFDStorage.KKMs[header.Id].Shifts[0].Operations, op)
		OFDStorage.KKMs[header.Id].Shifts[0].ID = 1
	}
	mutex.Unlock()
	return result, nil
}

func (m *MockOFD) SendSystemMessage(args ...interface{}) (interface{}, error) {
	header, err := findHeader(args)
	if err != nil {
		return nil, err
	}

	err = m.mockProcess()

	if err != nil {
		return nil, err
	}
	result := OfdResult{
		Header: header,
	}

	return result, nil
}

func (m *MockOFD) GetKKMInfo(args ...interface{}) (interface{}, error) {
	header, err := findHeader(args)
	if err != nil {
		return nil, err
	}

	err = m.mockProcess()

	if err != nil {
		return nil, err
	}
	result := OfdResult{
		Header: header,
	}

	return result, nil
}

func (m *MockOFD) SendOperationSale(document *models.Document, items []models.Position, args ...interface{}) (interface{}, error) {
	header, err := findHeader(args)
	if err != nil {
		return nil, err
	}

	err = m.mockProcess()

	if err != nil {
		return nil, err
	}
	result := OfdResult{
		Header: header,
	}
	addToOFD(document, header)
	fp := calculations.GetAfp(document.IdKkm)
	strFp := strconv.FormatUint(uint64(fp), 10)
	result.Response = append(result.Response, &kkm_proto.TicketResponse{
		TicketNumber: proto.String(strFp),
		QrCode:       nil,
	})
	return result, nil
}

func (m *MockOFD) SendOperationSaleReturn(document *models.Document, items []models.Position, args ...interface{}) (interface{}, error) {
	header, err := findHeader(args)
	if err != nil {
		return nil, err
	}

	err = m.mockProcess()

	if err != nil {
		return nil, err
	}
	result := OfdResult{
		Header: header,
	}
	addToOFD(document, header)
	fp := calculations.GetAfp(document.IdKkm)
	strFp := strconv.FormatUint(uint64(fp), 10)
	result.Response = append(result.Response, &kkm_proto.TicketResponse{
		TicketNumber: proto.String(strFp),
		QrCode:       nil,
	})
	return result, nil
}

func (m *MockOFD) SendOperationPurchase(document *models.Document, items []models.Position, args ...interface{}) (interface{}, error) {
	header, err := findHeader(args)
	if err != nil {
		return nil, err
	}

	err = m.mockProcess()

	if err != nil {
		return nil, err
	}
	result := OfdResult{
		Header: header,
	}
	addToOFD(document, header)
	fp := calculations.GetAfp(document.IdKkm)
	strFp := strconv.FormatUint(uint64(fp), 10)
	result.Response = append(result.Response, &kkm_proto.TicketResponse{
		TicketNumber: proto.String(strFp),
		QrCode:       nil,
	})
	return result, nil
}

func (m *MockOFD) SendOperationPurchaseReturn(document *models.Document, items []models.Position, args ...interface{}) (interface{}, error) {
	header, err := findHeader(args)
	if err != nil {
		return nil, err
	}

	err = m.mockProcess()

	if err != nil {
		return nil, err
	}
	result := OfdResult{
		Header: header,
	}
	fp := calculations.GetAfp(document.IdKkm)
	strFp := strconv.FormatUint(uint64(fp), 10)
	result.Response = append(result.Response, &kkm_proto.TicketResponse{
		TicketNumber: proto.String(strFp),
		QrCode:       nil,
	})
	return result, nil
}

func (m *MockOFD) SendTicketRollback(args ...interface{}) (interface{}, error) {
	header, err := findHeader(args)
	if err != nil {
		return nil, err
	}

	err = m.mockProcess()

	if err != nil {
		return nil, err
	}
	result := OfdResult{
		Header: header,
	}

	return result, nil
}

func (m *MockOFD) CloseShift(document *models.Document, args ...interface{}) (interface{}, error) {
	header, err := findHeader(args)
	if err != nil {
		return nil, err
	}

	err = m.mockProcess()

	if err != nil {
		return nil, err
	}
	result := OfdResult{
		Header: header,
	}
	addToOFD(document, header)
	return result, nil
}

func (m *MockOFD) RequestZReport(document *models.Document, args ...interface{}) (interface{}, error) {
	header, err := findHeader(args)
	if err != nil {
		return nil, err
	}

	err = m.mockProcess()

	if err != nil {
		return nil, err
	}
	result := OfdResult{
		Header: header,
	}

	return result, nil
}

func (m *MockOFD) RequestXReport(document *models.Document, args ...interface{}) (interface{}, error) {
	header, err := findHeader(args)
	if err != nil {
		return nil, err
	}

	err = m.mockProcess()

	if err != nil {
		return nil, err
	}
	result := OfdResult{
		Header: header,
	}
	addToOFD(document, header)
	return result, nil
}

func (m *MockOFD) DepositMoney(document *models.Document, args ...interface{}) (interface{}, error) {
	header, err := findHeader(args)
	if err != nil {
		return nil, err
	}

	err = m.mockProcess()

	if err != nil {
		return nil, err
	}
	result := OfdResult{
		Header: header,
	}
	addToOFD(document, header)

	return result, nil
}

func (m *MockOFD) WithdrawMoney(document *models.Document, args ...interface{}) (interface{}, error) {
	header, err := findHeader(args)
	if err != nil {
		return nil, err
	}

	err = m.mockProcess()

	if err != nil {
		return nil, err
	}
	result := OfdResult{
		Header: header,
	}
	addToOFD(document, header)

	return result, nil
}

func addToOFD(document *models.Document, header *OfdHeader) {
	mutex.Lock()
	defer mutex.Unlock()
	op := NewOperation(document.IdTypedocument, document.DateDocument)

	if document == nil {
		return
	}

	op.idShift = document.IdShift
	if OFDStorage != nil {
		v, ok := OFDStorage.KKMs[header.Id]
		if ok == true {
			v.ID = header.Id
			if v.Shifts[len(v.Shifts)-1].ID == uint32(document.Shift.ShiftIndex) {
				v.Shifts[len(v.Shifts)-1].Operations = append(v.Shifts[len(v.Shifts)-1].Operations, op)
				//v.Shifts[len(v.Shifts)-1].ID = uint32(document.Shift.ShiftIndex)
			} else {
				newShift := Shift{
					ID:      uint32(document.Shift.ShiftIndex),
					IdShift: document.IdShift,
				}
				newShift.Operations = append(newShift.Operations, op)
				v.Shifts = append(v.Shifts, newShift)
				OFDStorage.KKMs[header.Id] = v
			}
		} else {
			OFDStorage.KKMs[header.Id] = KKM{
				ID:     header.Id,
				Shifts: make([]Shift, 1),
			}

			OFDStorage.KKMs[header.Id].Shifts[0].Operations = append(OFDStorage.KKMs[header.Id].Shifts[0].Operations, op)
			OFDStorage.KKMs[header.Id].Shifts[0].ID = uint32(document.Shift.ShiftIndex)
		}
	}

}

func (m *MockOFD) SetTimeout(timeout time.Duration) {

}

func (m *MockOFD) mockProcess() (err error) {
	if m.NoConnection == true {
		return errors.New("could not connect to OFD server")
	} else if m.Down == true {
		return NewOfdError(uint32(kkm_proto.ResultTypeEnum_RESULT_TYPE_SERVICE_TEMPORARILY_UNAVAILABLE), "")
	}

	//failChance := rand.Intn(10)
	//
	//if failChance == 1 {
	//
	//	return NewOfdError(uint32(kkm_proto.ResultTypeEnum_RESULT_TYPE_INVALID_TOKEN), "")
	//}

	delay := time.Duration(rand.Intn(int(m.MaxDelay-m.MinDelay)) + int(m.MinDelay))
	time.Sleep(delay)

	return nil
}
