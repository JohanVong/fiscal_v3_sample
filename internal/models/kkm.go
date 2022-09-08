package models

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
