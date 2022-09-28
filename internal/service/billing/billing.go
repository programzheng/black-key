package billing

import (
	"time"

	"github.com/programzheng/black-key/internal/model/billing"
	"github.com/programzheng/black-key/internal/service"

	"github.com/jinzhu/copier"
)

type Billing struct {
	ID        uint `json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Title     string `json:"title"`
	Amount    int    `json:"amount"`
	Payer     string `json:"payer"`
	Note      string `json:"note"`

	service.Page
}

func (b *Billing) Add() (Billing, error) {
	model := billing.Billing{}
	copier.Copy(&model, &b)
	result, err := model.Add()
	if err != nil {
		return Billing{}, err
	}
	billing := Billing{}
	copier.Copy(&billing, &result)

	return billing, nil
}
