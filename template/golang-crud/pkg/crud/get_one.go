package crud

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) GetOne(c *gin.Context) {
	id := c.Param("id")
	fmt.Println("get one", id)

	c.Status(http.StatusOK)
}
