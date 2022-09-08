package calculations

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/martinlindhe/gogost/gost28147"
	"github.com/martinlindhe/gogost/gost341194"

	"github.com/JohanVong/fiscal_v3_sample/configs"
	"github.com/JohanVong/fiscal_v3_sample/internal/db"
	"github.com/JohanVong/fiscal_v3_sample/internal/models"
)

func CalculateNDS(price *decimal.Big, ndsPercent int64) *decimal.Big {
	price_f64, _ := price.Float64()
	np_f64 := float64(ndsPercent)
	nds_f64 := (price_f64 * np_f64) / (np_f64 + 100)
	nds_dm := decimal.New(int64(nds_f64*1000000), 6)
	return nds_dm
}

func CalculateDecNDS(price *decimal.Big, ndsPercent models.Decimal) *decimal.Big {
	//defer recoverpanic.RecoverPanic()
	var sum = new(decimal.Big)

	return sum.Mul(new(decimal.Big).Quo(price, new(decimal.Big).Add(decimal.New(1, 0), &ndsPercent.Big)), &ndsPercent.Big)
}

func CalculateMD(price, marckup, discount int) int {
	p := float64(price)
	m := float64(marckup)
	d := float64(discount)

	res := p + p*(m/100) - p*(d/100)
	return int(math.Round(res))
}

func CalculateDiscountMarkupSum(price int, percent int) int {
	return int(math.Round(float64(percent) * float64(price) / 100))
}

func CalculateChecksum(d models.Document) (string, error) {
	//defer recoverpanic.RecoverPanic()
	var hasher struct {
		Shift, User, Kkm int
		Value            decimal.Big
		Number           string
		Date             time.Time
		Signature        []byte
	}
	var (
		keyfile       string
		data, rawdata []byte
		err           error
		cert          *x509.Certificate
	)
	keyfile = configs.Gostkey
	h := gost341194.New(&gost28147.GostR3411_94_CryptoProParamSet)
	rawdata, err = ioutil.ReadFile(keyfile)
	if err != nil {
		return "", err
	}
	cert, err = x509.ParseCertificate(rawdata)
	if err != nil {
		return "", err
	}
	//key = cert.Signature
	//-----------------------------------
	hasher.Shift = d.IdShift
	hasher.User = d.IdUser
	hasher.Kkm = d.IdKkm
	hasher.Value = d.Value.Big
	hasher.Number = d.NumberDoc
	hasher.Date = d.DateDocument
	hasher.Signature = cert.Signature
	//-------------------------------------------
	data, err = json.Marshal(hasher)
	if err != nil {
		return "", err
	}
	h.Write(data)
	hash := h.Sum(nil)
	return fmt.Sprintf("%x", hash), nil
}

func IsChecksumValid(d models.Document) bool {
	//defer recoverpanic.RecoverPanic()
	checksum, err := CalculateChecksum(d)
	if err != nil {
		return false
	} else {
		if checksum == d.Checksum {
			return true
		} else {
			return false
		}
	}
}

func CalculateChain(s, last string) string {
	//defer recoverpanic.RecoverPanic()
	var hasher struct {
		CurrentDocHash string
		PrevDocHash    string
		Signature      string
	}
	//var doc models.Document
	//models.Orm.Last(&doc)
	hasher.CurrentDocHash = s
	hasher.PrevDocHash = last
	hasher.Signature = configs.Signature
	h := gost341194.New(&gost28147.GostR3411_94_CryptoProParamSet)
	data, err := json.Marshal(hasher)
	if err != nil {
		return ""
	}
	h.Write(data)
	hash := h.Sum(nil)
	return fmt.Sprintf("%x", hash)
}

func DocValidator(cur, prev models.Document) bool {
	//defer recoverpanic.RecoverPanic()
	var hasher struct {
		CurrentDocHash string
		PrevDocHash    string
		Signature      string
	}
	if !IsChecksumValid(cur) {
		return false
	}
	if !IsChecksumValid(prev) {
		return false
	}
	hasher.CurrentDocHash = cur.Checksum
	hasher.PrevDocHash = prev.Checksum
	hasher.Signature = configs.Signature
	data, err := json.Marshal(hasher)
	if err != nil {
		return false
	}
	h := gost341194.New(&gost28147.GostR3411_94_CryptoProParamSet)
	h.Write(data)
	hash := fmt.Sprintf("%x", h.Sum(nil))
	if hash == cur.DocChain {
		return true
	} else {
		return false
	}
}

func IsChainValid(from, to time.Time) bool {
	var (
		docs  []models.Document
		valid bool
	)
	db.Orm.Where("DateDocument BETWEEN ? AND ?", from, to).Find(&docs)
	for i, doc := range docs[1:] {
		valid = DocValidator(doc, docs[i-1])
		if !valid {
			return false
		}
	}
	return true
}

func GetAfp(kkmid int) uint32 {
	//defer recoverpanic.RecoverPanic()
	sum := sha256.Sum256([]byte(strconv.Itoa(kkmid) + strconv.Itoa(int(time.Now().UnixNano()))))
	ap := binary.BigEndian.Uint32(sum[:])
	if ap > uint32(math.MaxInt32) {
		ap = GetAfp(kkmid)
	}
	return ap
}

func PutOch(t int64) {
	models.Och <- t
}

func AddOperTime() {
	//defer recoverpanic.RecoverPanic()
	var t int64

	//ch := make(chan int64, 150)
	for {
		t = <-models.Och
		models.AvOp.Count++
		models.AvOp.Time = models.AvOp.Time + t
	}

}

func GetAvrgOpertime() int64 {
	//defer recoverpanic.RecoverPanic()
	if models.AvOp.Count == 0 {
		return 0
	} else {
		tt := models.AvOp.Time / int64(models.AvOp.Count)
		models.AvOp.Count = 1
		models.AvOp.Time = tt

		return tt
	}

}
