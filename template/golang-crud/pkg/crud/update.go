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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	fn "github.com/mrsimonemms/openfaas-templates/template/golang-crud/function"
)

func (h handler) Update(c *gin.Context) {
	id := c.Param("id")

	log := logger.Get(c).With().Str("action", "create").Str("recordId", id).Logger()

	log.Debug().Msg("Updating record")

	data := function.Schema{}
	coll := mgm.Coll(&data)

	filter := bson.D{
		{
			Key:   "id",
			Value: id,
		},
	}

	if err := coll.FindOne(mgm.Ctx(), filter).Decode(&data); err != nil {
		status := http.StatusServiceUnavailable
		var msg []gin.H
		if err == mongo.ErrNoDocuments {
			log.Debug().Msg("Record not found")

			status = http.StatusNotFound
		} else {
			log.Error().Err(err).Msg("Error finding record to update")

			msg = append(msg, gin.H{
				"message": err.Error(),
			})
		}

		common.ErrorHandler(c, status, msg...)
		return
	}

	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Err(err).Msg("Error reading request body")
		common.ErrorHandler(c, http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Just update the schema - avoids setting anything we manage, eg "id"
	update := function.Schema{
		Schema: data.Schema,
	}
	if err := json.Unmarshal(jsonData, &update); err != nil {
		log.Error().Err(err).Msg("Error decoding JSON")
		common.ErrorHandler(c, http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	// Put the updated schema back in the data
	data.Schema = update.Schema
	if validationErrs, err := common.Validate(data.Schema, fn.Validation); validationErrs != nil {
		if validationErrs != nil {
			log.Debug().Err(err).Msg("Error validating input data")
			common.ErrorHandler(c, http.StatusBadRequest, gin.H{
				"message": validationErrs,
			})
			return
		}
	}

	if err := coll.Update(&data); err != nil {
		log.Error().Err(err).Msg("Failed to update record")
		common.ErrorHandler(c, http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	log.Debug().Msg("Record updated")

	c.JSON(http.StatusOK, data)
}
