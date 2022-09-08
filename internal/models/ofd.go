package models

// Ofd represents a row from 'dbo.OFD'.
type Ofd struct {
	Id   int    `gorm:"PRIMARY_KEY;column:idOFD"` // idOFD
	Name string `gorm:"column:Name"`              // Name
	Url  string `gorm:"column:OfdUrl" json:"Url"`
}

// TableName sets the insert table name for this struct type
func (self Ofd) TableName() string {
	return "OFD"
}
