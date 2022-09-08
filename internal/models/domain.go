package models

// Domain represents a row from 'dbo.Domain'.
type Domain struct {
	Id   int    `gorm:"PRIMARY_KEY;column:idDomain"` // idDomain
	Name string `gorm:"column:Name"`                 // Name
}

// TableName sets the insert table name for this struct type
func (self Domain) TableName() string {
	return "Domain"
}
