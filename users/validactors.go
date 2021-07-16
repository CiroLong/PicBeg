package users

import (
	"WebSummerCamp/common"

	"github.com/gin-gonic/gin"
)

type LoginValidator struct {
	User struct {
		Username string `json:"username" form:"username" binding:"required,alphanum,min=4,max=255"`
		Password string `json:"password" form:"password" binding:"required,min=8,max=255"`
	} `json:"user"`
	UserModel UserModel `json:"-"`
}

func (L *LoginValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, L)
	if err != nil {
		return err
	}

	L.UserModel.UserName = L.User.Username
	return nil
}

func NewLoginValidator() LoginValidator {
	l := LoginValidator{}
	return l
}
