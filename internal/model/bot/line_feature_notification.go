package bot

import (
	"time"

	"github.com/programzheng/black-key/internal/model"
	"gorm.io/gorm"
)

type LineFeatureNotification struct {
	gorm.Model
	Feature      string    `gorm:"not null"`
	PushCycle    string    `gorm:"not null"`
	PushDateTime time.Time //發送時間，YYYY-MM-DD HH:MM:SS的格式代表有指定時間
	Limit        int       //限制次數(-1為不限制)
	UserID       string
	GroupID      string
	RoomID       string
	Request      string
}

func (lfn *LineFeatureNotification) Get(maps map[string]interface{}, not map[string]interface{}) ([]*LineFeatureNotification, error) {
	var lfns []*LineFeatureNotification
	err := model.DB.Where(maps).Not(not).Find(&lfns).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return lfns, nil
}

func (lfn *LineFeatureNotification) Save() error {
	err := model.DB.Save(lfn).Error
	if err != nil {
		return err
	}
	return nil
}

func (lfn *LineFeatureNotification) Delete() error {
	err := model.DB.Delete(&lfn).Error
	if err != nil {
		return err
	}
	return nil
}
