package crud

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	logger "github.com/mrsimonemms/gin-structured-logger"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/common"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/function"
	"go.mongodb.org/mongo-driver/bson"
)

func (h handler) Delete(c *gin.Context) {
	coll := mgm.Coll(&function.Schema{})

	id := c.Param("id")

	log := logger.Get(c).With().Str("action", "create").Str("recordId", id).Logger()

	log.Debug().Msg("Deleting record")

	result, err := coll.DeleteOne(mgm.Ctx(), bson.D{
		{
			Key:   "id",
			Value: id,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete record")
		common.ErrorHandler(c, http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if result.DeletedCount == 0 {
		log.Debug().Msg("Record unknown - cannot delete")
		common.ErrorHandler(c, http.StatusNotFound)
		return
	}

	log.Debug().Msg("Record deleted")

	// Send no content with JSON headers
	c.JSON(http.StatusNoContent, struct{}{})
}
