package models

// Typeobject represents a row from 'dbo.TypeObject'.
type Typeobject struct {
	Id   int    `gorm:"PRIMARY_KEY;column:idTypeObject"` // idTypeObject
	Name string `gorm:"column:Name"`                     // Name
}

// TableName sets the insert table name for this struct type
func (self Typeobject) TableName() string {
	return "TypeObject"
}
