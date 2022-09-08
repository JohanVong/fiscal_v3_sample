package models

import (
	"database/sql"

	"github.com/JohanVong/fiscal_v3_sample/internal/db"
)

// Permuser represents a row from 'dbo.PermUser'.
type Permuser struct {
	Id           int        `gorm:"PRIMARY_KEY;column:idPermUser"` // idPermUser
	IdGroupUser  int        `gorm:"column:idGroupUser"`            // idGroupUser
	IdTypeObject int        `gorm:"column:idTypeObject"`           // idTypeObject
	Nameobject   string     `gorm:"column:NameObject"`             // NameObject
	IdTypePerm   int        `gorm:"column:idTypePerm"`             // idTypePerm
	GroupUser    Groupuser  `gorm:"foreignkey:IdGroupUser;PRELOAD:true"`
	TypeObject   Typeobject `gorm:"foreignkey:IdTypeObject;PRELOAD:true"`
	TypePerm     Typeperm   `gorm:"foreignkey:IdTypePerm;PRELOAD:true"`
}

// TableName sets the insert table name for this struct type
func (self Permuser) TableName() string {
	return "PermUser"
}

func AddPermission(perm *Permuser) (sql.Result, error) {
	return db.Orm.DB().Exec("INSERT INTO PermUser (idGroupUser, idTypeObject, idTypePerm, NameObject) VALUES ((SELECT idGroupUser FROM GroupUser WHERE Name=?), (SELECT idTypeObject FROM TypeObject WHERE Name=?), (SELECT idTypePerm FROM TypePerm WHERE Name=?), ?)", perm.GroupUser.Name, perm.TypeObject.Name, perm.TypePerm.Name, perm.Nameobject)
}
