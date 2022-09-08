package models

import (
	"encoding/json"
)

// Balance represents a row from 'dbo.Balance'.
type Balance struct {
	Id            int         `gorm:"PRIMARY_KEY;column:idBalance"` // idBalance
	IdKkm         int         `gorm:"column:idKKM"`                 // idKKM
	Amount        Decimal     `gorm:"column:Amount"`                // Amount
	IdTypeBalance int         `gorm:"column:idTypeBalance"`         // idTypeBalance
	Type          Typebalance `gorm:"foreignkey:IdTypeBalance;auto_preload:true; association_autoupdate:false"`
}

// TableName sets the insert table name for this struct type
func (self Balance) TableName() string {
	return "Balance"
}

func (balance *Balance) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID          int         `json:"Id"`
		IdKkm       int         `json:"IdKKM"`
		Amount      Decimal     `json:"Amount"`
		TypeBalance Typebalance `json:"TypeBalance"`
	}{
		ID:          balance.Id,
		IdKkm:       balance.IdKkm,
		Amount:      balance.Amount,
		TypeBalance: balance.Type,
	})
}

type ResponseBalances struct {
	Balances         []Balance
	FiscalNumber     uint64 `json:"FiscalNumber"`
	AutonomousNumber uint32 `json:"AutonomousNumber"`
	IdDocument       int    `json:"IdDocument"`
	Location         string `json:"Location"`
	Receipt          string
}
