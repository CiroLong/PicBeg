package users

import (
	"WebSummerCamp/common"
	"errors"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	ID           uint   `gorm:"primaryKey"`
	UserName     string `gorm:"type:varchar(20);column:username;unique_index"`
	PasswordHash string `gorm:"column:password;not null"`
}

func NewUser(username, password string) error {
	db := common.GetDB()

	if len(password) == 0 {
		return errors.New("password shoule not be empty")
	}
	//password hash
	bytepassword := []byte(password)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytepassword, bcrypt.DefaultCost)

	user := UserModel{
		UserName:     username,
		PasswordHash: string(passwordHash),
	}
	db.Create(&user)
	os.MkdirAll("./stored_imgs/"+username, os.ModePerm)
	return nil
}

func FindOneUser(username string) (UserModel, error) {
	db := common.GetDB()
	var user UserModel
	result := db.Where("username = ?", username).First(&user)
	if result.RowsAffected == 0 {
		return user, errors.New("no such user")
	} else {
		return user, nil
	}
}

func (u *UserModel) CheckPassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.PasswordHash)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

/*func Saveone(data interface{}) error {
	db := common.GetDB()
	err := db.Save(data).Error
	return err
}*/
