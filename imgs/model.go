package imgs

import "github.com/jinzhu/gorm"

type ImgModel struct {
	gorm.Model
	FileName      string `gorm:"filename"`
	Path          string `gorm:"column:path"` //path in local
	OwnerUserName string `gorm:"column:owner_username"`
}
