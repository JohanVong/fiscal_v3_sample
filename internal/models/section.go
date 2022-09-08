package models

import (
	"github.com/ericlagergren/decimal"
)

// Section represents a row from 'dbo.Section'.
type Section struct {
	Id          int    `gorm:"PRIMARY_KEY;column:idSection"` // idSection
	Name        string `gorm:"column:Name"`                  // Name
	IdKkm       int    `gorm:"column:idKKM"`                 // idKKM
	IdCompany   int    `gorm:"column:idCompany"`             //
	Nds         int64  `gorm:"column:NDS"`                   // NDS
	Active      bool   `gorm:"column:Active"`
	SectionType int    `gorm:"column:SectionType"`
}

// TableName sets the insert table name for this struct type
func (self Section) TableName() string {
	return "Section"
}

type Repsec struct {
	Name                                       string
	Sales, Purchases, Refunds, PurchaseRefunds *decimal.Big
}
