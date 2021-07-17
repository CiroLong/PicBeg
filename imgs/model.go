package imgs

import (
	"WebSummerCamp/common"
	"errors"

	"github.com/jinzhu/gorm"
)

type ImgModel struct {
	gorm.Model
	FileName      string `gorm:"filename"`
	Path          string `gorm:"column:path"` //path in local such as "/stored_imgs/jack/test.img"
	OwnerUserName string `gorm:"column:owner_username"`
}

func NewImg(filename, username string) error {
	db := common.GetDB()
	if len(filename) == 0 {
		return errors.New("filename can't be empty")
	}
	model := ImgModel{
		FileName:      filename,
		Path:          "./stored_imgs/" + username + "/" + filename,
		OwnerUserName: username,
	}
	db.Create(&model)
	return nil
}

//return a model found and nil, or return error for not found
func FindOneImg(filename, username string) (ImgModel, error) {
	db := common.GetDB()
	var model ImgModel
	result := db.Where("owner_username = ? AND file_name = ?", username, filename).First(&model)
	if result.RowsAffected == 0 {
		return model, errors.New("no such user")
	} else {
		return model, nil
	}
}

func UpdateImg(model ImgModel) error {
	db := common.GetDB()
	db.Model(&model).Update()
	return nil
}

func DeleteImg(model ImgModel) error {
	db := common.GetDB()
	db.Delete(&model)
	return nil
}
