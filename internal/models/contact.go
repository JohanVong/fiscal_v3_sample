package models

// Contact represents a row from 'dbo.Contact'.
type Contact struct {
	Id            int         `gorm:"PRIMARY_KEY;column:idContact"` // idContact
	IdTypecontact int         `gorm:"column:idTypeContact"`         // idTypeContact
	Name          string      `gorm:"column:Name"`                  // Name
	TypeContact   Typecontact `gorm:"foreignkey:IdTypecontact;auto_preload:true"`
}

// TableName sets the insert table name for this struct type
func (self Contact) TableName() string {
	return "Contact"
}
