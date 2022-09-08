package models

// Typeaddress represents a row from 'dbo.TypeAddress'.
type Typeaddress struct {
	Id   int    `gorm:"PRIMARY_KEY;column:idTypeAddress"` // idTypeAddress
	Name string `gorm:"column:Name"`                      // Name
}

// TableName sets the insert table name for this struct type
func (self Typeaddress) TableName() string {
	return "TypeAddress"
}
