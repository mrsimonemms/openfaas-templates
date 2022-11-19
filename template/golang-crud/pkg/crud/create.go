package crud

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/common"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/function"

	fn "github.com/mrsimonemms/openfaas-templates/template/golang-crud/function"
)

func (h handler) Create(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		common.ErrorHandler(c, http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	var data function.Schema
	if err := json.Unmarshal(jsonData, &data); err != nil {
		common.ErrorHandler(c, http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	if err := common.Validate(data.Schema, fn.Validation); err != nil {
		if err != nil {
			common.ErrorHandler(c, http.StatusBadRequest, gin.H{
				"message": err,
			})
			return
		}
	}

	err = mgm.Coll(&data).Create(&data)
	if err != nil {
		common.ErrorHandler(c, http.StatusServiceUnavailable, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}
