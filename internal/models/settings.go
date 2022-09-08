package models

type SystemConfig struct {
	Id          int    `gorm:"PRIMARY_KEY;column:id"`
	OptionName  string `gorm:"column:option"`
	OptionValue string `gorm:"column:value"`
}

func (self *SystemConfig) TableName() string {
	return "system_config"
}
