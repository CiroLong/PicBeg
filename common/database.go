package common

import (
	"fmt"
	"io/ioutil"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

func Init() *gorm.DB {
	data, err := ioutil.ReadFile("./password.txt")
	if err != nil {
		panic(err)
	}
	var pwd = string(data)

	db, err := gorm.Open("mysql", pwd)
	if err != nil {
		fmt.Println("db err:(Init)", err)
	}
	db.DB().SetMaxIdleConns(10)
	DB = db
	return DB
}

func GetDB() *gorm.DB {
	return DB
}
