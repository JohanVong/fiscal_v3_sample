// Package models contains the types for schema 'dbo'.
package models

// Region represents a row from 'dbo.Region'.
type Region struct {
	Id   int    `gorm:"PRIMARY_KEY;column:idRegion"` // idRegion
	Name string `gorm:"column:Name"`                 // Name
}

// TableName sets the insert table name for this struct type
func (self Region) TableName() string {
	return "Region"
}
