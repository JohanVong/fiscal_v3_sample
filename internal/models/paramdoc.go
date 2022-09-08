package models

// Paramdoc represents a row from 'dbo.ParamDoc'.
type Paramdoc struct {
	Id             int    `gorm:"PRIMARY_KEY;column:idParamDoc"` // idParamDoc
	Value          string `gorm:"column:Value"`                  // Value
	IdTypevalue    int    `gorm:"column:idTypeValue"`            // idTypeValue
	IdDocuments    int    `gorm:"column:idDocuments"`            // idDocuments
	IdTypeparamdoc int    `gorm:"column:idTypeParamDoc"`         // idTypeParamDoc
}

// TableName sets the insert table name for this struct type
func (self Paramdoc) TableName() string {
	return "ParamDoc"
}
