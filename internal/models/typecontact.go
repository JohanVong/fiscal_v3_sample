package models

// Typecontact represents a row from 'dbo.TypeContact'.
type Typecontact struct {
	Id   int    `gorm:"PRIMARY_KEY;column:idTypeContact"` // idTypeContact
	Name string `gorm:"column:Name"`                      // Name
}

// TableName sets the insert table name for this struct type
func (self Typecontact) TableName() string {
	return "TypeContact"
}
