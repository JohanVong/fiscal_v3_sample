// Package models contains the types for schema 'dbo'.
package models

//import "github.com/jinzhu/gorm"

//import "encoding/json"

// Code generated by xo. DO NOT EDIT.

// Kkm represents a row from 'dbo.KKM'.
type Kkm struct {
	Id           int       `gorm:"PRIMARY_KEY;column:idKKM"` // idKKM
	IdCompany    int       `gorm:"column:idCompany"`         // idCompany
	IdModeKkm    int       `gorm:"column:idModeKKM"`         // idModeKKM
	IdOfd        int       `gorm:"column:idOFD"`             // idOFD
	IdNk         int       `gorm:"column:idNK"`              // idNK
	IdPlaceUsed  int       `gorm:"column:idPlaceUsed"`       // idPlaceUsed
	IdAddress    int       `gorm:"column:idAddress"`         // idOwnership
	Name         string    `gorm:"column:Name"`              // Name
	IdStatusKkm  int       `gorm:"column:idStatusKKM"`       // idStatusKKM
	Rnm          string    `gorm:"column:RNM"`               // RNM
	Lock         bool      `gorm:"column:Lock"`              // Lock
	Company      Company   `gorm:"foreignkey:IdCompany;auto_preload:true; association_autoupdate:false"`
	ModeKKM      ModeKkm   `gorm:"foreignkey:IdModeKkm;auto_preload:true; association_autoupdate:false" json:"-"`
	PlaceUsed    PlaceUsed `gorm:"foreignkey:IdPlaceUsed;auto_preload:true; association_autoupdate:false"`
	StatusKKM    StatusKkm `gorm:"foreignkey:IdStatusKkm;auto_preload:true; association_autoupdate:false"`
	Ofd          Ofd       `gorm:"foreignkey:IdOfd;auto_preload:true; association_autoupdate:false"`
	Nk           Nk        `gorm:"foreignkey:IdNk;auto_preload:true; association_autoupdate:false"`
	Address      Address   `gorm:"foreignkey:IdAddress;auto_preload:true; association_autoupdate:false"`
	IdCPCR       uint32    `gorm:"column:idCPCR"`
	TokenCPCR    uint32    `gorm:"column:tokenCPCR" json:"TokenCPCR"`
	ReqNumCPCR   uint16    `gorm:"column:reqnumCPCR" json:"-"`
	OfflineQueue uint32    `gorm:"column:OfflineQueue" json:"OfflineQueue"`
	ShiftIndex   int       `gorm:"column:ShiftIndex"`
	IdShift      int       `gorm:"column:idShift"`
	Allowed      bool      `gorm:"-" json:"Allowed"`
	IdSection    uint      `gorm:"column:IdSection" json:"IdSection"`
	Znm          int       `gorm:"column:znm" json:"Znm"`
	IsActive     bool      `gorm:"column:IsActive" json:"IsActive"`
	Autorenew    bool      `gorm:"column:Autorenew" json:"Autorenew"`
	OfdCode      int       `gorm:"column:ofd_code" json:"OfdCode"`
}

// TableName sets the insert table name for this struct type
func (self Kkm) TableName() string {
	return "KKM"
}

/*
func (self *Kkm) AfterCreate(scope *gorm.Scope) (err error) {
	if self.Znm == 0 {
		var znm []int
		scope.DB().Debug().Model(self).Pluck("znm", &znm)
		if len(znm) > 0 {
			self.Znm = znm[0]
		}
	}

	return
}
*/

/*
func (kkm *Kkm) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID        int       `json:"Id"`
		IdCompany int       `json:"IdCompany"`
		Mode      ModeKkm   `json:"Mode"`
		OFD       Ofd       `json:"Ofd"`
		NK        Nk        `json:"Nk"`
		PlaceUsed PlaceUsed `json:"PlaceUsed"`
		Name      string    `json:"Name"`
		Status    StatusKkm `json:"Status"`
		Rnm       string    `json:"Rnm"`
		Allowed   bool      `json:"Allowed"`
	}{
		ID:        kkm.Id,
		IdCompany: kkm.IdCompany,
		Mode:      kkm.ModeKKM,
		OFD:       kkm.Ofd,
		NK:        kkm.Nk,
		PlaceUsed: kkm.PlaceUsed,
		Name:      kkm.Name,
		Status:    kkm.StatusKKM,
		Rnm:       kkm.Rnm,
		Allowed:   kkm.Allowed,
	})
}
*/