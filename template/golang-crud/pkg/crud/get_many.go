package crud

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	logger "github.com/mrsimonemms/gin-structured-logger"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/common"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/function"
	querybuilder "github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/queryBuilder"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h handler) GetMany(c *gin.Context) {
	log := logger.Get(c).With().Str("action", "getMany").Logger()

	log.Debug().Msg("Getting multiple records")

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
	sort := c.QueryMap("sort")
	if len(sort) == 0 {
		sort["created_at"] = "DESC"
	}

	var sortBy bson.D
	for key, direction := range sort {
		var dirValue int
		switch strings.ToUpper(direction) {
		case "ASC":
			dirValue = 1
		case "DESC":
			dirValue = -1
		default:
			log.Debug().Err(err).Str("direction", direction).Msg("Invalid sort direction received")

			common.ErrorHandler(c, http.StatusBadRequest, gin.H{
				"message": "Invalid direction: " + direction + ". Use ASC or DESC",
			})
			return
		}

		sortBy = append(sortBy, bson.E{
			Key:   key,
			Value: dirValue,
		})
	}

	ctx := mgm.Ctx()
	data := function.Schema{}
	coll := mgm.Coll(&data)

	filter, err := querybuilder.New(c.QueryArray("filter"), c.QueryArray("or"))
	if err != nil {
		log.Debug().Err(err).Msg("Query builder error")

		common.ErrorHandler(c, http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Get total number of pages
	count, err := coll.CountDocuments(ctx, filter)
	if err != nil {
		log.Error().Err(err).Msg("Failed to count documents in database")

		common.ErrorHandler(c, http.StatusServiceUnavailable, gin.H{
			"message": err.Error(),
		})
		return
	}

	offset := (page - 1) * resultsPerPage
	totalPages := math.Ceil(float64(count) / float64(resultsPerPage))

	opts := options.Find().
		SetSort(sortBy).
		SetLimit(int64(resultsPerPage)).
		SetSkip(int64(offset))

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find documents in database")

		common.ErrorHandler(c, http.StatusServiceUnavailable, gin.H{
			"message": err.Error(),
		})
		return
	}

	results := make([]function.Schema, 0)
	if err := cursor.All(ctx, &results); err != nil {
		log.Error().Err(err).Msg("Failed to retrieve documents from database")

		common.ErrorHandler(c, http.StatusServiceUnavailable, gin.H{
			"message": err.Error(),
		})
		return
	}

	countResult := len(results)

	log.Debug().
		Int("count", countResult).
		Int("page", page).
		Float64("totalPages", totalPages).
		Int64("totalResults", count).
		Msg("Retrieved records")

	c.JSON(http.StatusOK, common.Paginate{
		Count:        countResult,
		Page:         page,
		TotalPages:   int(totalPages),
		TotalResults: int(count),
		Data:         results,
	})
}
