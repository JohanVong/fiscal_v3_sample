package models

// Cashier represents a row from 'dbo.Cashier'.
type Cashier struct {
	IdCashier       int           `gorm:"PRIMARY_KEY;column:idCashier"` // idCashier
	IdCompany       int           `gorm:"column:idCompany"`             // idCompany
	IdUser          int           `gorm:"column:idUser"`                // idUser
	Fio             string        `gorm:"column:FIO"`                   // FIO
	IdStatuscashier int           `gorm:"column:idStatusCashier"`       // idStatusCashier
	Lock            bool          `gorm:"column:Lock"`                  // Lock
	IdShift         int           `gorm:"column:idShift"`
	Company         Company       `gorm:"foreignkey:IdCompany;auto_preload:true"`
	User            User          `gorm:"foreignkey:IdUser;auto_preload:true"`
	StatusCashier   Statuscashier `gorm:"foreignkey:Idstatuscashier;auto_preload:true"`
}

// TableName sets the insert table name for this struct type
func (self Cashier) TableName() string {
	return "Cashier"
}
