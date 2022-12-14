package models

// Statususer represents a row from 'dbo.StatusUser'.
type Statususer struct {
	Id   int    `gorm:"PRIMARY_KEY;column:idStatusUser"` // idStatusUser
	Name string `gorm:"column:Name"`                     // Name
}

// TableName sets the insert table name for this struct type
func (self Statususer) TableName() string {
	return "StatusUser"
}
