package ofd

import (
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/JohanVong/fiscal_v3_sample/internal/db"
)

func InitDB2(db *gorm.DB) {
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
}

func TestProcessOfflineQueue2(t *testing.T) {

	var err error
	db.Orm, err = gorm.Open("sqlite3", "./testDB.db")
	db.Orm.DB().SetMaxOpenConns(1)
	if err != nil {
		panic(err)
	}
	defer func() {
		log.Println(os.Remove("./testDB.db"))
	}()
	InitDB2(db.Orm)

}
