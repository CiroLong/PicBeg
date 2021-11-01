package main

import (
	"WebSummerCamp/common"
	"WebSummerCamp/controller"
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
	imgchange.Use(controller.Login)
	{
		imgchange.GET("/:name/:imgname", imgs.GetFunc)
		imgchange.PUT("/:name/:imgname", imgs.PutFunc)   //更新
		imgchange.POST("/:name/:imgname", imgs.PostFunc) //新建
		imgchange.DELETE("/:name/:imgname", imgs.DeleteFunc)
		imgchange.GET("/:name", imgs.GetAllPath)
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
