package models

// PlaceUsed represents a row from 'dbo.PlaceUsed'.
type PlaceUsed struct {
	Id   int    `gorm:"PRIMARY_KEY;column:idPlaceUsed"` // idPlaceUsed
	Name string `gorm:"column:Name"`                    // Name
}

// TableName sets the insert table name for this struct type
func (self PlaceUsed) TableName() string {
	return "PlaceUsed"
}
