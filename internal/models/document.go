package models

import (
	"time"
)

// Document represents a row from 'dbo.Documents'.
type Document struct {
	Id               int          `gorm:"PRIMARY_KEY;column:idDocuments"` // idDocuments
	Uid              string       `gorm:"column:Uid" json:"-"`
	IdShift          int          `gorm:"column:idShift"`                                // idShift
	IdUser           int          `gorm:"column:idUser"`                                 // idUser
	IdTypedocument   int          `json:"-" gorm:"column:idTypeDocument"`                // idTypeDocument
	IdKkm            int          `gorm:"column:idKKM"`                                  // idKKM
	IdCompany        int          `gorm:"column:idCompany"`                              //
	DateDocument     time.Time    `gorm:"column:DateDocument;default:CURRENT_TIMESTAMP"` // DateDocument
	NumberDoc        string       `gorm:"column:NumberDoc"`                              // NumberDoc
	Checksum         string       `gorm:"column:CheckSum"`                               // CheckSum
	DocChain         string       `json:"-" gorm:"column:DocChain"`                      // CheckSum
	IdDomain         int          `gorm:"column:idDomain"`                               // idDomain
	Value            Decimal      `gorm:"column:Value"`                                  // К оплате
	Cash             Decimal      `gorm:"column:Cash"`                                   //Получено нал
	NonCash          Decimal      `gorm:"column:NonCash"`                                //Получено безнал
	Mobile           Decimal      `gorm:"column:Mobile"`                                 //Получено мобильным переводом
	Coins            int          `gorm:"column:Coins"`                                  //Получено монетами
	Change           Decimal      `gorm:"column:Change"`                                 //Сдача
	FiscalNumber     uint64       `gorm:"column:FiscalNumber"`
	AutonomousNumber uint32       `gorm:"column:AutonomousNumber"`
	CustomerIIN      string       `gorm:"column:CustomerIIN"` // ИИН покупателя
	CheckLink        string       `json:"-" gorm:"column:CheckLink"`
	Offline          bool         `gorm:"column:Offline" json:"Offline"`
	Token            uint32       `gorm:"column:token" json:"-"`
	ReqNum           uint16       `gorm:"column:reqNum" json:"-"`
	Domain           Domain       `gorm:"foreignkey:IdDomain;auto_preload:true; association_autoupdate:false;association_autocreate:false"`
	TypeDocument     Typedocument `gorm:"foreignkey:IdTypedocument;auto_preload:true; association_autoupdate:false;association_autocreate:false"`
	KKM              Kkm          `json:"KKM,omitempty" gorm:"foreignkey:IdKkm;auto_preload:true; association_autoupdate:false;association_autocreate:false"`
	User             *User        `json:"User,omitempty" gorm:"foreignkey:IdUser;auto_preload:true; association_autoupdate:false;association_autocreate:false"`
	Shift            *Shift       `json:"Shift,omitempty" gorm:"foreignkey:IdShift;PRELOAD:false; association_autoupdate:false;association_autocreate:false"`
	ReceiptDate      time.Time    `gorm:"column:ReceiptDate;default:CURRENT_TIMESTAMP"`

	Positions []*Position `json:"Positions,omitempty" gorm:"-"`
}

// TableName sets the insert table name for this struct type
func (self Document) TableName() string {
	return "Documents"
}
