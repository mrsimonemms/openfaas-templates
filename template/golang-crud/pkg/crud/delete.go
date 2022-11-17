package crud

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h handler) Delete(c *gin.Context) {
	fmt.Println("delete")

	c.Status(http.StatusOK)
}
