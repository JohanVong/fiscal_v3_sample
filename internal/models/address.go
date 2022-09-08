package models

// Address represents a row from 'dbo.Address'.
type Address struct {
	Id     int    `gorm:"PRIMARY_KEY;column:idAddress"` // idAddress
	IdTown int    `gorm:"column:idTown"`                // idTown
	Street string `gorm:"column:Street"`                // Street
	House  string `gorm:"column:House"`                 // House
	Flat   string `gorm:"column:Flat"`                  // Flat
	Town   Town   `gorm:"foreignkey:IdTown; auto_preload:true"`
}

// TableName sets the insert table name for this struct type
func (self Address) TableName() string {
	return "Address"
}
