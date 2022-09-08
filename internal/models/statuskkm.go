package models

// StatusKkm represents a row from 'dbo.StatusKKM'.
type StatusKkm struct {
	Id   int    `gorm:"PRIMARY_KEY;column:idStatusKKM"` // idStatusKKM
	Name string `gorm:"column:Name"`                    // Name
}

// TableName sets the insert table name for this struct type
func (self StatusKkm) TableName() string {
	return "StatusKKM"
}
