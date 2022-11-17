package crud

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) GetMany(c *gin.Context) {
	fmt.Println("get many")

	c.Status(http.StatusOK)
}
