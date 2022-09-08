package models

// Typevalue represents a row from 'dbo.TypeValue'.
type Typevalue struct {
	Id   int    `gorm:"PRIMARY_KEY;column:idTypeValue"` // idTypeValue
	Name string `gorm:"column:Name"`                    // Name
}

// TableName sets the insert table name for this struct type
func (self Typevalue) TableName() string {
	return "TypeValue"
}
