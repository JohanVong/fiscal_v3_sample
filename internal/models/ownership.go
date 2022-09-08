package models

// Ownership represents a row from 'dbo.Ownership'.
type Ownership struct {
	Id   int    `gorm:"PRIMARY_KEY;column:idOwnership"` // idOwnership
	Name string `gorm:"column:Name"`                    // Name
}

// TableName sets the insert table name for this struct type
func (self Ownership) TableName() string {
	return "Ownership"
}
