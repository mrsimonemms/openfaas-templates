package crud

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) Update(c *gin.Context) {
	fmt.Println("update")

	c.Status(http.StatusOK)
}
