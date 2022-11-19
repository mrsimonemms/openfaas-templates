package crud

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/common"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/function"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	fn "github.com/mrsimonemms/openfaas-templates/template/golang-crud/function"
)

func (h handler) Update(c *gin.Context) {
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

	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
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
		common.ErrorHandler(c, http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	// Put the updated schema back in the data
	data.Schema = update.Schema
	if err := common.Validate(data.Schema, fn.Validation); err != nil {
		if err != nil {
			common.ErrorHandler(c, http.StatusBadRequest, gin.H{
				"message": err,
			})
			return
		}
	}

	if err := coll.Update(&data); err != nil {
		common.ErrorHandler(c, http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}
