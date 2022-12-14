package models

// Groupuser represents a row from 'dbo.GroupUser'.
type Groupuser struct {
	Id   int    `gorm:"PRIMARY_KEY;column:idGroupUser"` // idGroupUser
	Name string `gorm:"column:Name"`                    // Name
}

// TableName sets the insert table name for this struct type
func (self Groupuser) TableName() string {
	return "GroupUser"
}
