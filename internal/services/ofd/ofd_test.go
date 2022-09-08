package ofd_test

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"

	"github.com/JohanVong/fiscal_v3_sample/internal/db"
	"github.com/JohanVong/fiscal_v3_sample/internal/models"
	"github.com/JohanVong/fiscal_v3_sample/internal/services/graylog"
	"github.com/JohanVong/fiscal_v3_sample/internal/services/ofd"
	"github.com/JohanVong/fiscal_v3_sample/internal/transport/protobuf/ofd/cpcr"
)

const (
	NumberOfOperations = 700
	NumberOfKKMs       = 20
)

var controlOFDStorage *cpcr.MockOFDStorage

func InitDB(db *gorm.DB) {
	f, err := os.OpenFile("./testData.sql", os.O_RDONLY, 777)
	if err != nil {
		panic(err)
	}
	stat, err := f.Stat()
	log.Println(err)
	b := make([]byte, stat.Size())
	log.Println(f.Read(b))
	err = db.Exec(string(b)).Error
	if err != nil {
		panic(err)
	}

	td := []models.Typedocument{
		{
			Id:   1,
			Name: "Продажа",
		},
		{
			Id:   2,
			Name: "Служебный приход",
		},
		{
			Id:   3,
			Name: "Служебный расход",
		},
		{
			Id:   4,
			Name: "Покупка",
		},
		{
			Id:   5,
			Name: "Возврат продажи",
		},
		{
			Id:   6,
			Name: "Открытие смены",
		},
		{
			Id:   7,
			Name: "Закрытие смены",
		},
		{
			Id:   8,
			Name: "x-отчет",
		},
	}

	for _, t := range td {
		err = db.Create(&t).Error
		if err != nil {
			panic(err)
		}
	}
}

func ReturnRandomTypeDocument() int {
	rand.Seed(time.Now().UTC().UnixNano())
	typeDoc := rand.Intn(8) + 1
	if typeDoc == 6 {
		typeDoc = ReturnRandomTypeDocument()
	}
	return typeDoc
}

func PopulateDBGeneral(db *gorm.DB) {

	for i := 1; i <= NumberOfKKMs; i++ {
		kkm := models.Kkm{
			IdCompany:   1,
			IdModeKkm:   1,
			IdOfd:       1,
			IdNk:        1,
			IdPlaceUsed: 1,
			IdAddress:   i,
			Name:        "KKM " + strconv.Itoa(i),
			IdStatusKkm: 1,
			Rnm:         "303457803247",
			Lock:        false,
			IdCPCR:      uint32(i),
		}

		shift := models.Shift{
			IdUser:        1,
			IdKkm:         i,
			IdStatusShift: 1,
			DateOpen:      time.Now().Add(time.Hour * -11),
			ShiftIndex:    1,
		}

		db.Create(&kkm)
		db.Create(&shift)

		controlOFDStorage.KKMs[uint32(i)] = cpcr.KKM{
			ID:     uint32(i),
			Shifts: make([]cpcr.Shift, 1),
		}

		controlOFDStorage.KKMs[uint32(i)].Shifts[0].ID = 1

	}

	startTime := time.Now().Add(time.Hour * -10)
	timeInc := time.Duration(0)
	for i := 1; i <= NumberOfOperations; i++ {
		rand.Seed(time.Now().UTC().UnixNano())
		typeDoc := ReturnRandomTypeDocument()
		kkm := rand.Intn(NumberOfKKMs) + 1
		shift := models.Shift{}
		db.Model(&models.Shift{}).Where("idStatusShift=?", 1).Where("idKKM=?", kkm).Take(&shift)
		if typeDoc == 7 {
			db.Model(&models.Shift{}).Where("idKKM=?", kkm).Where("idStatusShift=?", 1).Updates(map[string]interface{}{"IdStatusShift": 2, "DateClose": startTime.Add(timeInc)})

			db.Create(&models.Shift{
				IdUser:        1,
				IdKkm:         shift.IdKkm,
				IdStatusShift: 1,
				DateOpen:      startTime.Add(timeInc),
				ShiftIndex:    shift.ShiftIndex + 1,
			})
		}
		doc := models.Document{
			IdShift:          shift.Id,
			IdUser:           1,
			IdTypedocument:   typeDoc,
			IdKkm:            kkm,
			DateDocument:     startTime.Add(timeInc),
			NumberDoc:        "",
			Checksum:         "",
			DocChain:         "",
			IdDomain:         0,
			Value:            *new(models.Decimal),
			Cash:             *new(models.Decimal),
			NonCash:          *new(models.Decimal),
			Coins:            0,
			Change:           *new(models.Decimal),
			FiscalNumber:     0,
			AutonomousNumber: 0,
			CheckLink:        "",
			Offline:          true,
		}
		db.Create(&doc)

		db.Model(&models.Kkm{}).Where("idKKM=?",
			kkm).Updates(map[string]interface{}{"OfflineQueue": gorm.Expr("OfflineQueue + ?", 1)})
		op := cpcr.NewOperation(doc.IdTypedocument, doc.DateDocument)

		if len(controlOFDStorage.KKMs[uint32(kkm)].Shifts[0].Operations) == 0 {
			controlOFDStorage.KKMs[uint32(kkm)].Shifts[0].Operations = append(controlOFDStorage.KKMs[uint32(kkm)].Shifts[0].Operations, cpcr.NewOperation(0, doc.DateDocument))
		}

		v, ok := controlOFDStorage.KKMs[uint32(kkm)]
		if ok == true {
			v.ID = uint32(kkm)
			if v.Shifts[len(v.Shifts)-1].ID == uint32(shift.ShiftIndex) {
				v.Shifts[len(v.Shifts)-1].Operations = append(v.Shifts[len(v.Shifts)-1].Operations, op)
				v.Shifts[len(v.Shifts)-1].ID = uint32(shift.ShiftIndex)
			} else {
				newShift := cpcr.Shift{
					ID: uint32(shift.ShiftIndex),
				}
				newShift.Operations = append(newShift.Operations, op)
				v.Shifts = append(v.Shifts, newShift)
				controlOFDStorage.KKMs[uint32(kkm)] = v
			}
		} else {
			controlOFDStorage.KKMs[uint32(kkm)] = cpcr.KKM{
				ID:     uint32(kkm),
				Shifts: make([]cpcr.Shift, 1),
			}

			controlOFDStorage.KKMs[uint32(kkm)].Shifts[0].Operations = append(controlOFDStorage.KKMs[uint32(kkm)].Shifts[0].Operations, op)
			controlOFDStorage.KKMs[uint32(kkm)].Shifts[0].ID = uint32(shift.ShiftIndex)
		}

		timeInc += time.Duration(rand.Intn(90)+30) * time.Second

	}
}

func TestProcessOfflineQueue(t *testing.T) {
	controlOFDStorage = &cpcr.MockOFDStorage{
		KKMs: make(map[uint32]cpcr.KKM),
	}
	cpcr.OFDStorage = &cpcr.MockOFDStorage{
		KKMs: make(map[uint32]cpcr.KKM),
	}
	var err error
	graylog.GraylogClient, err = graylog.NewClient("192.168.151.110:12222")
	if err != nil {
		panic(err)
	}

	mockOFD := cpcr.MockOFD{
		Down:     false,
		MinDelay: time.Millisecond,
		MaxDelay: time.Millisecond * 2,
	}

	ofd.OfdPool = make(ofd.Pool)

	ofd.OfdPool[1] = &mockOFD
	rand.Seed(int64(time.Now().Unix()))
	db.Orm, err = gorm.Open("sqlite3", "./testDB.db")
	db.Orm.DB().SetMaxOpenConns(1)
	if err != nil {
		panic(err)
	}
	defer func() {
		log.Println(os.Remove("./testDB.db"))
	}()
	InitDB(db.Orm)
	PopulateDBGeneral(db.Orm)

	db.RedisCl = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, redis_err := db.RedisCl.Ping().Result()
	if redis_err != nil {
		panic(redis_err)
	}

	models.KKMOFDStatusStorage = make(map[int]*models.KKMOFDStatus)

	go ofd.ProcessOfflineQueueV2(time.Second * 3)

	time.Sleep(time.Second * 7)

	//var wg sync.WaitGroup
	//go ofd.ProcessOfflineQueue(&wg, time.Second*3)
	//
	//time.Sleep(time.Second * 45)
	//wg.Wait()

	log.Println(controlOFDStorage.String())
	log.Println("********************************************************")
	log.Println(cpcr.OFDStorage.String())
	for k, v := range cpcr.OFDStorage.KKMs {
		//var kkms []models.Kkm
		var documents []models.Document

		for i, s := range v.Shifts {
			db.Orm.Joins("INNER JOIN Shift ON Documents.idShift=Shift.idShift").Where("Documents.idKKM=?", k).Where("Shift.ShiftIndex=?", s.ID).Find(&documents)

			var queue int

			row := db.Orm.Select("offlinequeue").Table("KKM").Where("idkkm=?", k).Row()
			err = row.Scan(&queue)

			j := 0
			var lenDocs, lenOps int

			for _, o := range s.Operations {
				if o.GetType() > 0 {
					lenOps++
				}
			}

			lenDocs = len(documents)

			if lenDocs != lenOps {
				t.Error(fmt.Sprintf("Mismatch on number of opertaions on kkm %d Shift %d should be %d, received %d",
					v.ID, s.ID, lenDocs, lenOps))
				continue
			}

			for _, o := range s.Operations {
				otype := o.GetType()
				if otype == 0 {
					continue
				}
				if o.GetDate().Equal(documents[j].DateDocument) == false {
					t.Error(fmt.Sprintf("Mismatch on KKM %d, Shift %d, Operation %d: date sould be %s, received %s", k, i+1, j+i, documents[j].DateDocument.String(), s.Operations[j].GetDate().String()))
				}

				if o.GetType() != documents[j].IdTypedocument {
					t.Error(fmt.Sprintf("Mismatch on KKM %d, Shift %d, Operation %d: Name sould be %d, received %d", k, i+1, j+i, documents[j].IdTypedocument, s.Operations[j].GetType()))
				}
				j++
			}
		}
	}

}
