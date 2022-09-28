package bot

import (
	"fmt"
	"time"

	"github.com/programzheng/black-key/internal/model"
	"gorm.io/gorm"
)

type LineNotification struct {
	gorm.Model
	Service      string    `gorm:"not null"` //"Messaging API"
	PushCycle    string    `gorm:"not null"` //發送週期
	PushDateTime time.Time //發送時間，YYYY-MM-DD HH:MM:SS的格式代表有指定時間
	Limit        int       //限制次數(-1為不限制)
	UserID       string
	GroupID      string
	RoomID       string
	Type         string `gorm:"not null"`
	Template     string `gorm:"type:text"`
}

func (ln *LineNotification) Add() (*LineNotification, error) {
	if err := model.DB.Create(&ln).Error; err != nil {
		return nil, err
	}

	return ln, nil
}

func LineNotificationFirstByID(id uint) (*LineNotification, error) {
	ln := &LineNotification{}
	err := model.DB.First(ln, id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return ln, nil
}

func (ln *LineNotification) Get(maps map[string]interface{}, not map[string]interface{}) ([]*LineNotification, error) {
	var lns []*LineNotification
	err := model.DB.Where(maps).Not(not).Find(&lns).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return lns, nil
}

func (ln *LineNotification) GetByPushDateTimeRange(comparison string, dateTime string) ([]*LineNotification, error) {
	var lns []*LineNotification

	conditional := fmt.Sprintf("push_date_time %s ?", comparison)
	err := model.DB.Where(conditional, dateTime).Find(&lns).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return lns, nil
}

func (ln *LineNotification) GetByUpdatedRange(comparison string, dateTime string) ([]*LineNotification, error) {
	var lns []*LineNotification

	conditional := fmt.Sprintf("updated_at %s ?", comparison)
	err := model.DB.Where(conditional, dateTime).Find(&lns).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return lns, nil
}

func (ln *LineNotification) Save() error {
	err := model.DB.Save(ln).Error
	if err != nil {
		return err
	}
	return nil
}

func (ln *LineNotification) Delete() error {
	err := model.DB.Delete(&ln).Error
	if err != nil {
		return err
	}
	return nil
}

func (ln *LineNotification) PermanentlyDelete() error {
	err := model.DB.Unscoped().Delete(&ln).Error
	if err != nil {
		return err
	}
	return nil
}
