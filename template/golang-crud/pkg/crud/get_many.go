package crud

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/common"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/function"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h handler) GetMany(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		// Invalid input - use default value
		page = 1
	}
	resultsPerPage, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		// Invalid input - use default value
		resultsPerPage = int(h.Config.Limit)
	}
	if resultsPerPage > int(h.Config.MaxLimit) {
		// Limit request is too high - set to max
		resultsPerPage = int(h.Config.MaxLimit)
	}

	ctx := mgm.Ctx()
	data := function.Schema{}
	coll := mgm.Coll(&data)
	filter := bson.D{}

	// Get total number of pages
	count, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		common.ErrorHandler(c, http.StatusServiceUnavailable, gin.H{
			"message": err.Error(),
		})
		return
	}

	offset := (page - 1) * resultsPerPage
	totalPages := math.Ceil(float64(count) / float64(resultsPerPage))

	opts := options.Find().
		SetSort(bson.D{
			{
				Key:   "created_at",
				Value: -1,
			},
		}).
		SetLimit(int64(resultsPerPage)).
		SetSkip(int64(offset))

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		common.ErrorHandler(c, http.StatusServiceUnavailable, gin.H{
			"message": err.Error(),
		})
		return
	}

	results := make([]function.Schema, 0)
	if err := cursor.All(ctx, &results); err != nil {
		common.ErrorHandler(c, http.StatusServiceUnavailable, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.Paginate{
		Count:        len(results),
		Page:         page,
		TotalPages:   int(totalPages),
		TotalResults: int(count),
		Data:         results,
	})
}
