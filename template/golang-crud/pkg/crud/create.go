package crud

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	logger "github.com/mrsimonemms/gin-structured-logger"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/common"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/function"

	fn "github.com/mrsimonemms/openfaas-templates/template/golang-crud/function"
)

func (h handler) Create(c *gin.Context) {
	log := logger.Get(c).With().Str("action", "create").Logger()

	log.Debug().Msg("Creating record")

	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Err(err).Msg("Error reading request body")
		common.ErrorHandler(c, http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	var data function.Schema
	if err := json.Unmarshal(jsonData, &data); err != nil {
		log.Error().Err(err).Msg("Error decoding JSON")
		common.ErrorHandler(c, http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	if validationErrs, err := common.Validate(data.Schema, fn.Validation); validationErrs != nil {
		if validationErrs != nil {
			log.Debug().Err(err).Msg("Error validating input data")
			common.ErrorHandler(c, http.StatusBadRequest, gin.H{
				"message": validationErrs,
			})
			return
		}
	}

	err = mgm.Coll(&data).Create(&data)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create record")
		common.ErrorHandler(c, http.StatusServiceUnavailable, gin.H{
			"message": err.Error(),
		})
		return
	}

	log.Debug().Msg("Record created")

	c.JSON(http.StatusOK, data)
}
