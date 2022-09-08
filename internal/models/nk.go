package models

// Nk represents a row from 'dbo.NK'.
type Nk struct {
	Id       int    `gorm:"PRIMARY_KEY;column:idNK"` // idNK
	Name     string `gorm:"column:Name"`             // Name
	IdRegion int    `gorm:"column:idRegion"`         // idRegion
	Code     int    `gorm:"column:Code"`             // Code
}

// TableName sets the insert table name for this struct type
func (self Nk) TableName() string {
	return "NK"
}
