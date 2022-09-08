package models

// Company represents a row from 'dbo.Company'.
type Company struct {
	Id          int       `gorm:"PRIMARY_KEY;column:idCompany"` // idCompany
	IdUser      int       `gorm:"column:idUser"`                // idUser
	IdOwnership int       `gorm:"column:idOwnership"`           // idOwnership
	IdAddress   int       `gorm:"column:idAddress"`             // idOwnership
	Bin         string    `gorm:"column:BIN"`                   // BIN
	ShortName   string    `gorm:"column:ShortName"`             // ShortName
	FullName    string    `gorm:"column:FullName"`              // FullName
	Fio         string    `gorm:"column:FIO"`                   // FIO
	Nds         string    `gorm:"column:NDS"`                   // NDS
	User        User      `json:"-" gorm:"foreignkey:IdUser;auto_preload:true"`
	Ownership   Ownership `gorm:"foreignkey:IdOwnership;auto_preload:true"`
	Address     Address   `gorm:"foreignkey:IdAddress;auto_preload:true"`
	KkmsAmount  int       `gorm:"-"`
}

// TableName sets the insert table name for this struct type
func (self Company) TableName() string {
	return "Company"
}
