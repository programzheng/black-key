package file

import (
	"time"

	"github.com/programzheng/black-key/pkg/filesystem"
	"github.com/programzheng/black-key/pkg/helper"
	"github.com/programzheng/black-key/pkg/model"

	"gorm.io/gorm"
)

type File struct {
	ID        uint       `gorm:"primary_key" json:"-"`
	HashID    string     `gorm:"unique"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`
	Reference string     `json:"-"`
	System    string     `json:"-"`
	Type      string
	Path      string `json:"-"`
	Name      string `json:"-"`
}

type Files []*File

func (f *File) AfterFind() (err error) {
	f.Path = filesystem.Driver.GetHostURL() + "/" + f.Path + f.Name
	return
}

func (f *File) AfterCreate(tx *gorm.DB) (err error) {
	// 設定給前端呼叫圖片的ID
	hashID := helper.ConvertToString(f.ID) + "_" + helper.ConvertToString(time.Now().Unix())
	hashID = helper.CreateMD5(hashID)
	tx.Model(f).Update("HashID", hashID)
	return
}

func (f File) Add() (File, error) {
	model.Migrate(&f)
	if err := model.DB.Save(&f).Error; err != nil {
		return File{}, err
	}
	return f, nil
}

func Get(ids []interface{}, maps interface{}) (Files, error) {
	var files Files
	if ids == nil {
		err := model.DB.Where(maps).Find(&files).Error

		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}

		return files, nil
	}
	err := model.DB.Where(ids).Where(maps).Find(&files).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return files, nil
}

func BatchUpdates(maps interface{}, updates interface{}) (Files, error) {
	var files Files
	err := model.DB.Model(&files).Where(maps).Updates(updates).Find(&files).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return files, nil
	// err := model.DB.Model(&files).Where(ids).Where(maps).Updates(updates).Find(&files).Error

	// if err != nil && err != gorm.ErrRecordNotFound {
	// 	return nil, err
	// }

	return files, nil
}

func (f File) Update() (File, error) {
	if err := model.DB.Save(&f).Error; err != nil {
		return File{}, err
	}
	return f, nil
}
