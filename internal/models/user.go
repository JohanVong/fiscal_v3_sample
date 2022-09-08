package models

import (
	"database/sql"

	"github.com/JohanVong/fiscal_v3_sample/internal/db"
)

type User struct {
	Id           int        `gorm:"PRIMARY_KEY; column:idUser"`   // idUser
	IdTypeUser   int        `json:"-" gorm:"Column:idTypeUser"`   // idTypeUser
	IdStatusUser int        `json:"-" gorm:"Column:idStatusUser"` // idStatusUser
	IdGroupUser  int        `json:"-" gorm:"Column:idGroupUser"`  // idGroupUser
	PhoneLogin   string     `gorm:"Column:PhoneLogin"`            // PhoneLogin
	Password     string     `json:"-" gorm:"Column:Password"`     // Password
	Name         string     `gorm:"Column:Name"`                  // Name
	Lock         bool       `gorm:"column:Lock"`                  // Lock
	IdShift      int        `gorm:"column:idShift"`
	StatusUser   Statususer `json:"-" gorm:"foreignkey:IdStatusUser;has_one;PRELOAD:true"`
	GroupUser    Groupuser  `json:"-" gorm:"foreignkey:IdGroupUser;has_one;PRELOAD:true"`
	TypeUser     Typeuser   `json:"-" gorm:"foreignkey:IdTypeUser;has_one;PRELOAD:true"`
}

func (u *User) TableName() string {
	return "User"
}

func GetUsers() (users []User, err error) {
	rows, err := db.Orm.DB().Query("SELECT u.idUser, t.idTypeUser, t.Name, s.idStatusUser, s.Name, g.idGroupUser, g.Name, u.PhoneLogin, u.Name, u.Password FROM [User] u JOIN TypeUser t ON t.idTypeUser=u.idTypeUser JOIN StatusUser s ON s.idStatusUser=u.idStatusUser JOIN GroupUser g ON g.idGroupUser=u.idGroupUser")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var (
			idUser         sql.NullInt64
			idTypeuser     sql.NullInt64
			idStatusUser   sql.NullInt64
			idGroupUser    sql.NullInt64
			typeUserName   sql.NullString
			statusUserName sql.NullString
			groupUserName  sql.NullString
			phoneLogin     sql.NullString
			userName       sql.NullString
			password       sql.NullString
		)

		err := rows.Scan(
			&idUser,
			&idTypeuser,
			&typeUserName,
			&idStatusUser,
			&statusUserName,
			&idGroupUser,
			&groupUserName,
			&phoneLogin,
			&userName,
			&password)

		user := User{
			Id:           int(idUser.Int64),
			IdTypeUser:   int(idTypeuser.Int64),
			IdStatusUser: int(idStatusUser.Int64),
			IdGroupUser:  int(idGroupUser.Int64),
			PhoneLogin:   phoneLogin.String,
			Password:     password.String,
			Name:         userName.String,
			StatusUser: Statususer{
				Id:   int(idStatusUser.Int64),
				Name: statusUserName.String,
			},
			TypeUser: Typeuser{
				Id:   int(idTypeuser.Int64),
				Name: typeUserName.String,
			},
		}
		if err != nil {
			return nil, err
		}
		users = append(users, user)

	}
	return
}
