package common

import (
	"github.com/gin-gonic/gin"
)

func Bind(c *gin.Context, obj interface{}) error {
	//b := binding.Default(c.Request.Method, c.ContentType())
	return c.ShouldBindQuery(obj)
}
