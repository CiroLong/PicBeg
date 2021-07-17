package imgs

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlePath struct {
	UserName string `uri:"name" binding:"required"`
	ImgName  string `uri:"imgname" binding:"required"`
}

func GetFunc(c *gin.Context) { //获取
	c.String(200, "ok")
}

func PutFunc(c *gin.Context) { //上传
	//绑定参数
	var handlePath HandlePath
	if err := c.ShouldBindUri(&handlePath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "router err",
		})
		return
	}
	//查询是否存在
	_, err := FindOneImg(handlePath.ImgName, handlePath.UserName)
	if err != nil { //不存在则保存在本地
		file, _ := c.FormFile("file")
		log.Println(file.Filename)
		//dst := "./stored_imgs/" + handlePath.UserName + "/" + file.Filename
		c.SaveUploadedFile(file, "./stored_imgs/"+handlePath.UserName+"/"+file.Filename)
		err := NewImg(file.Filename, handlePath.UserName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "there is some mistake!",
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": file.Filename + " uploaded!",
		})
		return
	}
	//存在则返回错误
	c.JSON(http.StatusBadRequest, gin.H{
		"code": 400,
		"msg":  "the file exits, use POST to update",
	})
}

func PostFunc(c *gin.Context) { //更新
	c.String(200, "ok")
}

func DeleteFunc(c *gin.Context) { //删除
	c.String(200, "ok")
}
