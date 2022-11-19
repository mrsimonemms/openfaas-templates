package crud

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/common"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/function"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h handler) GetOne(c *gin.Context) {
	data := function.Schema{}
	coll := mgm.Coll(&data)
	filter := bson.D{
		{
			Key:   "id",
			Value: c.Param("id"),
		},
	}

	if err := coll.FindOne(mgm.Ctx(), filter).Decode(&data); err != nil {
		status := http.StatusServiceUnavailable
		var msg []gin.H
		if err == mongo.ErrNoDocuments {
			status = http.StatusNotFound
		} else {
			msg = append(msg, gin.H{
				"message": err.Error(),
			})
		}

		common.ErrorHandler(c, status, msg...)
		return
	}

	c.JSON(http.StatusOK, data)
}
