package models

// Typekkmparam represents a row from 'dbo.TypeKKMParam'.
type Typekkmparam struct {
	Id   int    `gorm:"PRIMARY_KEY;column:idTypeKKMParam"` // idTypeKKMParam
	Name string `gorm:"column:Name"`                       // Name
}

// TableName sets the insert table name for this struct type
func (self Typekkmparam) TableName() string {
	return "TypeKKMParam"
}
