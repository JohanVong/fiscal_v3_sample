package models

// Statuscashier represents a row from 'dbo.StatusCashier'.
type Statuscashier struct {
	Id   int    `gorm:"PRIMARY_KEY;column:idStatusCashier"` // idStatusCashier
	Name string `gorm:"column:Name"`                        // Name
}

// TableName sets the insert table name for this struct type
func (self Statuscashier) TableName() string {
	return "StatusCashier"
}
