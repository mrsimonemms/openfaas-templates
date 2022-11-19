package crud

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/common"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/function"
	"go.mongodb.org/mongo-driver/bson"
)

func (h handler) Delete(c *gin.Context) {
	coll := mgm.Coll(&function.Schema{})

	result, err := coll.DeleteOne(mgm.Ctx(), bson.D{
		{
			Key:   "id",
			Value: c.Param("id"),
		},
	})
	if err != nil {
		common.ErrorHandler(c, http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if result.DeletedCount == 0 {
		common.ErrorHandler(c, http.StatusNotFound)
		return
	}

	// Send no content with JSON headers
	c.JSON(http.StatusNoContent, struct{}{})
}
