package models

// Unit represents a row from 'dbo.Unit'.
type Unit struct {
	Id        int    `gorm:"PRIMARY_KEY;column:idUnit"` // idUnit
	Code      int    `gorm:"column:Code"`               // Code
	NameRU    string `gorm:"column:NameRU"`             // NameRU
	NameKAZ   string `gorm:"column:NameKAZ"`            // NameKAZ
	ShortName string `gorm:"column:ShortName"`          // ShortName
}

// TableName sets the insert table name for this struct type
func (self Unit) TableName() string {
	return "Unit"
}
