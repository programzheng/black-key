package file

import (
	"time"

	"github.com/programzheng/black-key/internal/helper"
	"github.com/programzheng/black-key/internal/model"

	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	HashID      string  `gorm:"unique"`
	Reference   *string `json:"-"`
	System      string  `json:"-"`
	Type        string
	Path        string `json:"-"`
	Name        string `json:"-"`
	ThirdPatyID string `json:"-"`
}

type Files []*File

func (f *File) AfterCreate(tx *gorm.DB) (err error) {
	// 設定給前端呼叫圖片的ID
	hashID := helper.ConvertToString(f.ID) + "_" + helper.ConvertToString(time.Now().Unix())
	hashID = helper.CreateMD5(hashID)
	tx.Model(f).Update("HashID", hashID)
	return
}

func (f File) Add() (File, error) {
	if err := model.DB.Save(&f).Error; err != nil {
		return File{}, err
	}
	return f, nil
}

func Get(where map[string]interface{}) (Files, error) {
	var files Files
	err := model.DB.Where(where).Find(&files).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return files, nil
}

func BatchUpdates(where map[string]interface{}, updates map[string]interface{}) (Files, error) {
	var files Files
	err := model.DB.Model(&files).Where(where).Updates(updates).Find(&files).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return files, nil
}

func (f File) Update() (File, error) {
	if err := model.DB.Save(&f).Error; err != nil {
		return File{}, err
	}
	return f, nil
}

func Delete(where map[string]interface{}) error {
	err := model.DB.Where(where).Delete(File{}).Error
	if err != nil {
		return err
	}
	return nil
}
