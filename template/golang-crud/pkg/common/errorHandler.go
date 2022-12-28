package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context, code int, additional ...gin.H) {
	obj := gin.H{
		"error": http.StatusText(code),
	}

	for _, a := range additional {
		for k, v := range a {
			obj[k] = v
		}
	}

	c.JSON(code, obj)
}
