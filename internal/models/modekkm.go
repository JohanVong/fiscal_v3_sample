package models

// ModeKkm represents a row from 'dbo.ModeKKM'.
type ModeKkm struct {
	Id   int    `gorm:"PRIMARY_KEY;column:idModeKKM"` // idModeKKM
	Name string `gorm:"column:Name"`                  // Name
}

// TableName sets the insert table name for this struct type
func (self ModeKkm) TableName() string {
	return "ModeKKM"
}
