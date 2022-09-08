package models

// Typebalance represents a row from 'dbo.TypeBalance'.
type Typebalance struct {
	Id   int    `gorm:"PRIMARY_KEY;column:idTypeBalance" json:"Id"` // idTypeBalance
	Name string `gorm:"column:Name" json:"Name"`                    // Name
}

// TableName sets the insert table name for this struct type
func (self Typebalance) TableName() string {
	return "TypeBalance"
}
