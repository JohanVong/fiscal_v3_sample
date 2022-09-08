package models

// Таблица отношений пользователя к компании много ко многому.
type UserRel struct {
	Id          int       `gorm:"PRIMARY_KEY;column:Id"` // id
	IdUser      int       `gorm:"column:idUser"`         //
	IdCompany   int       `gorm:"column:idCompany"`      //
	Active      bool      `gorm:"Column:Active"`         // idStatusUser
	IdGroupUser int       `gorm:"Column:idGroupUser"`    // idGroupUser
	GroupUser   Groupuser `gorm:"foreignkey:IdGroupUser;PRELOAD:true"`
	User        User      `gorm:"foreignkey:IdUser;auto_preload:true"`
	Company     Company   `gorm:"foreignkey:IdCompany;auto_preload:true"`
}

// TableName sets the insert table name for this struct type
func (self UserRel) TableName() string {
	return "UserRel"
}
