package billing

import (
	"github.com/programzheng/black-key/pkg/model"

	"gorm.io/gorm"
)

type Billing struct {
	gorm.Model
	Title  string `gorm:"comment:'標題'"`
	Amount int    `gorm:"comment:'總付款金額'"`
	Payer  string `gorm:"comment:'付款人'"`
	Note   string `gorm:"comment:'備註'"`
}

func (b Billing) Add() (Billing, error) {
	if err := model.DB.Create(&b).Error; err != nil {
		return Billing{}, err
	}

	return b, nil
}
