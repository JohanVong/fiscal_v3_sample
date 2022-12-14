package models

// Position represents a row from 'dbo.Position'.
type Position struct {
	Id          int     `gorm:"PRIMARY_KEY;column:idPosition"` // idPosition
	IdDocument  int     `gorm:"column:idDocuments"`            // idDocuments
	IdArticle   int     `gorm:"column:idArticle"`              // idArticle
	IdSection   int     `gorm:"column:idSection"`              // idSection
	IdCompany   int     `gorm:"column:idCompany"`              //
	IdUnit      int     `gorm:"column:idUnit"`                 // Единица измерения
	Number      int     `gorm:"column:Number"`                 // Number
	Name        string  `gorm:"column:Name"`
	Price       Decimal `gorm:"column:Price"` // Price
	Nds         Decimal `gorm:"column:Nds"`   // Price
	NdsDiscount Decimal `gorm:"column:NdsDiscount"`
	NdsMarkup   Decimal `gorm:"column:NdsMarkup"`
	Discount    Decimal `gorm:"column:Discount"` // Discount
	Markup      Decimal `gorm:"column:Markup"`   // Наценка
	Qty         Decimal `gorm:"column:Qty"`      // Количество
	Total       Decimal `gorm:"column:Total"`    // Total
	Storno      bool    `gorm:"column:Storno"`
	ProductCode string  `gorm:"column:ProductCode"`
	Article     Article `json:",omitempty" gorm:"foreignkey:IdArticle;auto_preload:true"`
	Section     Section `json:",omitempty" gorm:"foreignkey:IdSection;auto_preload:true"`
	Unit        Unit    `json:",omitempty" gorm:"foreignkey:IdUnit;auto_preload:true"`
}

// TableName sets the insert table name for this struct type
func (self Position) TableName() string {
	return "Position"
}
