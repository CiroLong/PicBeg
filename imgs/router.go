package imgs

import (
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

type HandlePath struct {
	UserName string `uri:"name" binding:"required"`
	ImgName  string `uri:"imgname" binding:"required"`
}

func GetFunc(c *gin.Context) { //获取  ok
	//绑定参数
	var handlePath HandlePath
	if err := c.ShouldBindUri(&handlePath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "router err",
		})
		return
	}
	//获取路径
	model, err := FindOneImg(handlePath.ImgName, handlePath.UserName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err,
		})
	}
	//返回图片
	filepath := model.GetPath()
	_, errByfile := os.Open(filepath)

	//获取文件的名称
	fileName := path.Base(filepath)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	if errByfile != nil {
		log.Println("获取文件失败")
		c.Redirect(http.StatusFound, "/404")
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")

	c.File(filepath)
}

func PostFunc(c *gin.Context) { //上传  ok
	//绑定参数
	var handlePath HandlePath
	if err := c.ShouldBindUri(&handlePath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "router err",
		})
		return
	}
	var realuser = c.Query("username")
	if realuser != handlePath.UserName {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "No permission.",
		})
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

func PutFunc(c *gin.Context) { //更新 ok
	//绑定参数
	var handlePath HandlePath
	if err := c.ShouldBindUri(&handlePath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "router err",
		})
		return
	}
	var realuser = c.Query("username")
	if realuser != handlePath.UserName {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "No permission.",
		})
	}
	//查询是否存在
	_, err := FindOneImg(handlePath.ImgName, handlePath.UserName)
	if err != nil { //不存在则返回错误
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "the file doesn't exit, use PUT to upload",
		})
		return
	}
	//存在则更新
	file, _ := c.FormFile("file")
	log.Println(file.Filename)
	dst := "./stored_imgs/" + handlePath.UserName + "/" + file.Filename
	os.Remove(dst)
	c.SaveUploadedFile(file, dst)
	errByImg := NewImg(file.Filename, handlePath.UserName)
	if errByImg != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "there is some mistake!",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": file.Filename + " updated!",
	})
}

func DeleteFunc(c *gin.Context) { //删除 ok
	var handlePath HandlePath
	if err := c.ShouldBindUri(&handlePath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "router err",
		})
		return
	}
	var realuser = c.Query("username")
	if realuser != handlePath.UserName {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "No permission.",
		})
	}
	//查询是否存在
	model, err := FindOneImg(handlePath.ImgName, handlePath.UserName)
	if err != nil { //不存在则返回错误
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "the file doesn't exit",
		})
		return
	}
	//存在则删除
	dst := "./stored_imgs/" + handlePath.UserName + "/" + handlePath.ImgName
	os.Remove(dst)
	errByDelete := DeleteImg(model)
	if errByDelete != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "there is some mistake",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  handlePath.ImgName + " deleted",
	})
}

func GetAllPath(c *gin.Context) { //查看用户所有图片路径 ok
	var handlePath struct {
		UserName string `uri:"name" binding:"required"`
	} //写复杂了....
	if err := c.ShouldBindUri(&handlePath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "router err",
		})
		return
	}
	models, errByFind := FindAllImgs(handlePath.UserName)
	if errByFind != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "there is no img by user " + handlePath.UserName,
		})
	}
	var paths []string
	for _, model := range models {
		var path = model.GetPath()
		path = strings.Replace(path, "./stored_imgs/", "", -1)
		path = "/api/img/" + path
		paths = append(paths, path)
	}
	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"msg":   "ok",
		"paths": paths,
	})
}
