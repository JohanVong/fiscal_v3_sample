package models

// Kkmparam represents a row from 'dbo.KKMParam'.
type Kkmparam struct {
	Id             int    `gorm:"PRIMARY_KEY;column:idKKMParam"` // idKKMParam
	Value          string `gorm:"column:Value"`                  // Value
	IdTypekkmparam int    `gorm:"column:idTypeKKMParam"`         // idTypeKKMParam
	IdKkm          int    `gorm:"column:idKKM"`                  // idKKM
	IdTypevalue    int    `gorm:"column:idTypeValue"`            // idTypeValue
}

// TableName sets the insert table name for this struct type
func (self Kkmparam) TableName() string {
	return "KKMParam"
}
