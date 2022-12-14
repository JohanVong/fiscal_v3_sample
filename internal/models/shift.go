package models

import "time"

// Shift represents a row from 'dbo.Shift'.
type Shift struct {
	Id                  int         `gorm:"PRIMARY_KEY;column:idShift"` // idShift
	IdUser              int         `gorm:"column:idUser"`
	IdKkm               int         `gorm:"column:idKKM"`         // idKKM
	IdStatusShift       int         `gorm:"column:idStatusShift"` // idStatusShift
	IdCompany           int         `gorm:"column:idCompany"`     //
	DateOpen            time.Time   `gorm:"column:DateOpen"`      // DateOpen
	DateClose           time.Time   `gorm:"column:DateClose"`     // DateClose
	BalanceOpen         Decimal     `gorm:"column:BalanceOpen"`
	BalanceClose        Decimal     `gorm:"column:BalanceClose"`
	BalanceCloseNonCash Decimal     `gorm:"column:BalanceCloseNonCash"`
	KKM                 Kkm         `gorm:"foreignkey:IdKkm;auto_preload:true; association_autoupdate:false"`
	StatusShift         StatusShift `gorm:"foreignkey:IdStatusShift;auto_preload:true; association_autoupdate:false"`
	User                User        `gorm:"foreignkey:IdUser;auto_preload:true; association_autoupdate:false"`
	Company             Company     `gorm:"foreignkey:IdCompany;auto_preload:true; association_autoupdate:false"`
	ShiftIndex          int         `gorm:"column:ShiftIndex"`
}

// TableName sets the insert table name for this struct type
func (self Shift) TableName() string {
	return "Shift"
}
