package models

// StatusShift represents a row from 'dbo.StatusShift'.
type StatusShift struct {
	Id   int    `gorm:"PRIMARY_KEY;column:idStatusShift"` // idStatusShift
	Name string `gorm:"column:Name"`                      // Name
}

// TableName sets the insert table name for this struct type
func (self StatusShift) TableName() string {
	return "StatusShift"
}
