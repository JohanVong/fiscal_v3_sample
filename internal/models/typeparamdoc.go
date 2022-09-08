package models

// Typeparamdoc represents a row from 'dbo.TypeParamDoc'.
type Typeparamdoc struct {
	Id             int    `gorm:"PRIMARY_KEY;column:idTypeParamDoc"` // idTypeParamDoc
	Name           string `gorm:"column:Name"`                       // Name
	IdTypedocument int    `gorm:"column:idTypeDocument"`             // idTypeDocument
}

// TableName sets the insert table name for this struct type
func (self Typeparamdoc) TableName() string {
	return "TypeParamDoc"
}
