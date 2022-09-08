package models

// Article represents a row from 'dbo.Article'.
type Article struct {
	Id             int     `gorm:"PRIMARY_KEY;column:idArticle"` // idArticle
	Name           string  `gorm:"column:Name"`                  // Name
	Qr             string  `gorm:"column:QR"`                    // QR
	IdGroupArticle int     `gorm:"column:idGroupArticle"`        // idGroupArticle
	IdSection      int     `gorm:"column:idSection"`             // idSection
	IdCompany      int     `gorm:"column:idCompany"`             //
	IdUnit         int     `gorm:"column:idUnit"`                //
	IdKkm          int     `gorm:"column:idKKM"`                 //
	Price          Decimal `gorm:"column:Price"`                 // Price
	Discount       Decimal `gorm:"column:Discount"`              // Discount
	Markup         Decimal `gorm:"column:Markup"`                // Наценка
	Active         bool    `gorm:"Column:Active"`                // idStatusUser
}

// TableName sets the insert table name for this struct type
func (self Article) TableName() string {
	return "Article"
}
