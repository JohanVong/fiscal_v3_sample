package models

// Typedocument represents a row from 'dbo.TypeDocument'.
type Typedocument struct {
	Id   int    `gorm:"PRIMARY_KEY;column:idTypeDocument"` // idTypeDocument
	Name string `gorm:"column:Name"`                       // Name
}

// TableName sets the insert table name for this struct type
func (self Typedocument) TableName() string {
	return "TypeDocument"
}
