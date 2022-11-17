package crud

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) Create(c *gin.Context) {
	fmt.Println("create")

	c.Status(http.StatusOK)
}
