package models

// Town represents a row from 'dbo.Town'.
type Town struct {
	Id       int    `gorm:"PRIMARY_KEY;column:idTown"` // idTown
	Name     string `gorm:"column:Name"`               // Name
	TimeZone int    `gorm:"column:TimeZone"`           //time zone shift
	IdRegion int    `gorm:"column:idRegion"`           // idRegion
	Region   Region `gorm:"foreignkey:IdRegion;auto_preload:true"`
}

// TableName sets the insert table name for this struct type
func (self Town) TableName() string {
	return "Town"
}
