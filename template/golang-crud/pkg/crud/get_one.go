package crud

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	logger "github.com/mrsimonemms/gin-structured-logger"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/common"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/function"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h handler) GetOne(c *gin.Context) {
	id := c.Param("id")

	log := logger.Get(c).With().Str("action", "getOne").Str("recordId", id).Logger()

	log.Debug().Msg("Getting a single record")

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
			log.Error().Err(err).Msg("Error finding record")

			msg = append(msg, gin.H{
				"message": err.Error(),
			})
		}

		common.ErrorHandler(c, status, msg...)
		return
	}

	log.Debug().Msg("Record found")

	c.JSON(http.StatusOK, data)
}
