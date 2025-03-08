package models

import "gorm.io/gorm"

type IRWebsiteSetting struct {
	gorm.Model
	Title      string `gorm:"comment:网站标题"`
	Logo       string `gorm:"comment:logo"`
	Company    string `gorm:"comment:公式"`
	ICP        string `gorm:"comment:ICP"`
	Background string `gorm:"comment:网站背景"`
}

func (IRWebsiteSetting) TableName() string {
	return "ir_website_setting"
}

func (IRWebsiteSetting) Read() []string {
	return []string{
		"system",
	}
}

func (IRWebsiteSetting) Write() []string {
	return []string{
		"system",
	}
}

func (IRWebsiteSetting) Delete() []string {
	return []string{
		"system",
	}
}

func (IRWebsiteSetting) Update() []string {
	return []string{
		"system",
	}
}
