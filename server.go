package main

import (
	"WebSummerCamp/common"
	"WebSummerCamp/users"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	db := common.Init()
	defer db.Close()
	users.AutoMigrate()

	imgchange := r.Group("/api/img")
	imgchange.Use(Login)
	imgchange.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "ok",
		})
	})

	testpath := r.Group("/api/ping")
	testpath.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8080")
}

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