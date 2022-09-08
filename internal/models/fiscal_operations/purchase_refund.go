package fiscal_operations

import (
	"errors"
	"strconv"
	"time"

	"github.com/ericlagergren/decimal"

	db "github.com/JohanVong/fiscal_v3_sample/internal/db"
	validator "gopkg.in/asaskevich/govalidator.v9"

	"github.com/JohanVong/fiscal_v3_sample/internal/models"
	"github.com/JohanVong/fiscal_v3_sample/internal/services/calculations"
)

type PurchaseRefundReqest struct {
	IdKkm         uint64      `json:"-"`
	IdUser        int         `json:"-"`
	IdDomain      int         `json:"IdDomain" valid:"required~Поле IdDomain не должно быть пустым"`
	Cash          decimal.Big `json:"Cash" valid:"twoDecimalPlaces~Поле Cash не должно иметь более двух десятичных знаков,isNonNegative~поле Cash должно быть положительным"`
	NonCash       decimal.Big `json:"NonCash" valid:"twoDecimalPlaces~Поле NonCash не должно иметь более двух десятичных знаков,isNonNegative~поле NonCash должно быть положительным"`
	Positions     []Position  `json:"Positions" valid:"required~Поле Positions не должно быть пустым"`
	Total         decimal.Big `json:"Total" valid:"required~Поле Total не должно быть пустым,twoDecimalPlaces~Поле Total не должно иметь более двух десятичных знаков,isNonNegative~поле Total должно быть положительным"`
	Uid           string      `json:"Uid" valid:"required~Заголовок не содержит Uid, uuidv4~Uid не соответствует формату"`
	AFP           uint32      `json:"AFP"`
	ChkLocation   string      `json:"ChkLocation"`
	GenerateCheck bool        `json:"GenerateCheck"`
	ReceiptDate   time.Time   `json:"ReceiptDate"`
	IdCompany     int         `json:"-"`
	IdShift       int         `json:"-"`
	CustomerIIN   string      `json:"CustomerIIN"`
}

func (self PurchaseRefundReqest) Valid() (err error) {

	var (
		valid          bool
		position_price decimal.Big
		position_total decimal.Big
		control        decimal.Big
	)
	valid, err = validator.ValidateStruct(&self)
	if !valid {
		return
	}
	if self.NonCash.Cmp(&self.Total) > 0 {
		err = errors.New("Сумма безналичной оплаты не может быть больше стоимости товара")
		return
	}
	if self.NonCash.Cmp(&self.Total) == 0 && self.Cash.Cmp(decimal.New(0, 0)) > 0 {
		err = errors.New("Сумма оплаты наличным расчетом не может быть больше 0, если сумма оплаты безналичным расчетом равна общей сумме оплаты")
		return
	}

	for _, pos := range self.Positions {
		if pos.Storno {
			continue
		}
		position_price.Add(&pos.Price, new(decimal.Big).Sub(&pos.Markup, &pos.Discount))
		position_total.Copy(new(decimal.Big).Mul(&position_price, &pos.Qty))
		control.Add(&control, &position_total)
	}

	if control.Cmp(&self.Total) != 0 {
		err = errors.New("Сумма ИТОГО не верна")
		return
	}

	return
}

func (self PurchaseRefundReqest) Create() (doc models.Document, positions []models.Position, err error) {
	var (
		pos     models.Position
		address models.Address
		section models.Section
	)

	tx := db.Orm.Begin()
	defer func() {
		tx.RollbackUnlessCommitted()
	}()
	tx.Set("gorm:auto_preload", true).Where("idAddress = (SELECT idAddress from KKM WHERE idKKM = ?)").First(&address)

	doc.Uid = self.Uid
	doc.IdShift = int(self.IdShift)
	doc.CustomerIIN = self.CustomerIIN
	doc.IdUser = int(self.IdUser)
	doc.IdTypedocument = OPERATION_PURCHASE_REFUND
	doc.IdKkm = int(self.IdKkm)
	doc.IdCompany = int(self.IdCompany)
	doc.DateDocument = time.Now().In(time.FixedZone(address.Town.Name, address.Town.TimeZone*60*60))
	doc.NumberDoc = strconv.Itoa(int(self.IdCompany)) + strconv.Itoa(int(self.IdKkm)) + strconv.Itoa(int(time.Now().UnixNano()/100000000))
	doc.IdDomain = self.IdDomain
	doc.Value.Big = self.Total
	doc.Cash.Big = self.Cash
	doc.NonCash.Big = self.NonCash
	doc.Change.Copy(new(decimal.Big).Add(&self.Cash, new(decimal.Big).Sub(&self.NonCash, &self.Total)))
	doc.AutonomousNumber = self.AFP
	if self.AFP != 0 {
		doc.Offline = true
		doc.ReceiptDate = self.ReceiptDate
		err = tx.Exec("UPDATE KKM SET OfflineQueue = OfflineQueue + ? WHERE idKKM = ?", 1, doc.IdKkm).Error
		if err != nil {
			return
		}

	} else {
		doc.Offline = false
		doc.ReceiptDate = doc.DateDocument
	}
	doc.Checksum, _ = calculations.GetDocumentChecksum(doc)
	err = tx.Create(&doc).Error
	if err != nil {
		return
	}

	for i, itm := range self.Positions {
		if section.Id == 0 || section.Id != itm.IdSection {
			tx.First(&section, itm.IdSection)
		}
		priceDM := new(decimal.Big)
		pos.IdDocument = doc.Id
		pos.Number = i
		pos.IdCompany = int(self.IdCompany)
		pos.Price.Big = itm.Price
		pos.Discount.Big = itm.Discount
		pos.Markup.Big = itm.Markup
		pos.Qty.Big = itm.Qty
		pos.IdSection = itm.IdSection
		pos.Name = itm.Name

		priceDM.Add(&pos.Price.Big, new(decimal.Big).Sub(&pos.Markup.Big, &pos.Discount.Big))
		pos.Nds.Copy(calculations.CalculateNDS(priceDM, section.Nds))
		pos.NdsDiscount.Copy(calculations.CalculateNDS(&pos.Discount.Big, section.Nds))
		pos.NdsMarkup.Copy(calculations.CalculateNDS(&pos.Markup.Big, section.Nds))
		pos.Total.Copy(new(decimal.Big).Mul(priceDM, &itm.Qty))
		pos.Storno = itm.Storno
		pos.ProductCode = itm.ProductCode
		pos.IdUnit = itm.IdUnit
		err = tx.Create(&pos).Error
		if err != nil {
			return
		}
		pos = models.Position{}
	}
	// err = tx.Create(&positions).Error
	// if err != nil {
	// 	return
	// }
	if self.Cash.Cmp(new(decimal.Big)) > 0 {
		sum_to_balance := new(models.Decimal).Sub(&self.Cash, &doc.Change.Big)
		sum_f, _ := sum_to_balance.Float64()
		err = tx.Exec("UPDATE Balance SET Amount = Amount + ? WHERE idKKM = ? AND idTypeBalance = ?", sum_f, self.IdKkm, 1).Error
		if err != nil {
			return
		}
	}
	if self.NonCash.Cmp(new(decimal.Big)) > 0 {
		non_cash_f, _ := self.NonCash.Float64()
		err = tx.Exec("UPDATE Balance SET Amount = Amount + ? WHERE idKKM = ? AND idTypeBalance = ?", non_cash_f, self.IdKkm, 2).Error
		if err != nil {
			return
		}
	}
	err = tx.Set("gorm:auto_preload", true).Where("idDocuments = ?", doc.Id).Find(&positions).Error
	if err != nil {
		return
	}
	tx.Commit()

	return
}
