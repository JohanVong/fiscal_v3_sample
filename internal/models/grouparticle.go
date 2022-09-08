package models

// Grouparticle represents a row from 'dbo.GroupArticle'.
type Grouparticle struct {
	Id      int    `gorm:"PRIMARY_KEY;column:idGroupArticle"` // idGroupArticle
	Name    string `gorm:"column:Name"`                       // Name
	IdGroup int    `gorm:"column:IDGroup"`                    // IDGroup
}

// TableName sets the insert table name for this struct type
func (self Grouparticle) TableName() string {
	return "GroupArticle"
}
