package controller

import (
	"WebSummerCamp/users"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	//get data

	L := users.NewLoginValidator()
	err := L.Bind(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid patterns",
		})
		c.Abort()
		return
	}
	//check data in database
	L.UserModel, err = users.FindOneUser(L.User.Username)
	if err != nil {
		err2 := users.NewUser(L.User.Username, L.User.Password) //new user
		if err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err2,
			})
			c.Abort()
			return
		}
	} else {
		//check the password

		err3 := L.UserModel.CheckPassword(L.User.Password)
		if err3 != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "wrong password!",
			})
			c.Abort()
			return
		}
	}
	//fmt.Println(L)
	//else OK
}
