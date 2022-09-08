package models

import (
	"time"
)

type ZReport struct {
	Id                        int `gorm:"PRIMARY_KEY;AUTO_INCREMENT"` // idArticle
	ShiftIndex                int
	IdShift                   int `gorm:"column:idShift"`
	IdKkm                     int
	IdUser                    int `gorm:"column:idUser"`
	DateOpen                  time.Time
	DateClose                 time.Time
	BalanceOpen               Decimal
	BalanceClose              Decimal
	BalanceCloseNonCash       Decimal
	Count                     int
	SalesQty                  int
	SalesAmount               Decimal
	CumulativeSales           Decimal
	PurchasesQty              int
	PurchasesAmount           Decimal
	CumulativePurchases       Decimal
	ExpensesQty               int
	ExpensesAmount            Decimal
	CumulativeExpenses        Decimal
	RefundsQty                int
	RefundsAmount             Decimal
	CumulativeRefunds         Decimal
	PurchaseRefundsQty        int
	PurchaseRefundsAmount     Decimal
	CumulativePurchaseRefunds Decimal
	IncomesQty                int
	IncomesAmount             Decimal
	CumulativeIncomes         Decimal
	Sections                  map[int]Repsec `gorm:"-"`
	KKM                       Kkm            `gorm:"foreignkey:IdKkm;auto_preload:true;" json:"Kkm"`
	User                      User           `gorm:"foreignkey:IdUser;auto_preload:true;"`
}

// TableName sets the insert table name for this struct type
func (self ZReport) TableName() string {
	return "ZReport"
}
