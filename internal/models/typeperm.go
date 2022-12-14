package models

// Typeperm represents a row from 'dbo.TypePerm'.
type Typeperm struct {
	Id   int    `gorm:"PRIMARY_KEY;column:idTypePerm"` // idTypePerm
	Name string `gorm:"column:Name"`                   // Name
}

// TableName sets the insert table name for this struct type
func (self Typeperm) TableName() string {
	return "TypePerm"
}
