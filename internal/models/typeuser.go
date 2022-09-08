package models

// Typeuser represents a row from 'dbo.TypeUser'.
type Typeuser struct {
	Id   int    `gorm:"PRIMARY_KEY;column:idTypeUser"` // idTypeUser
	Name string `gorm:"column:Name"`                   // Name
}

// TableName sets the insert table name for this struct type
func (self Typeuser) TableName() string {
	return "TypeUser"
}
