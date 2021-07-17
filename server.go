package main

import (
	"WebSummerCamp/common"
	"WebSummerCamp/imgs"
	"WebSummerCamp/users"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	db := common.Init()
	defer db.Close()
	db.AutoMigrate(&users.UserModel{})
	db.AutoMigrate(&imgs.ImgModel{})

	// 为 multipart forms 设置较低的内存限制 (默认是 32 MiB)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB

	imgchange := r.Group("/api/img")
	imgchange.Use(Login)
	{
		imgchange.GET("/:name/:imgname", imgs.GetFunc)
		imgchange.PUT("/:name/:imgname", imgs.PutFunc)
		imgchange.POST("/:name/:imgname", imgs.PostFunc)
		imgchange.DELETE("/:name/:imgname", imgs.DeleteFunc)

		imgchange.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  "login ok",
			})
		})
	}
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
