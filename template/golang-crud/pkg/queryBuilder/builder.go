package querybuilder

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

type Query struct {
	Field     string
	Condition bson.E
	Value     []string
}

// https://github.com/nestjsx/crud/wiki/Requests#filter-conditions
func QueryFromString(conditionQuery string) (query Query, err error) {
	parts := strings.Split(conditionQuery, "||")

	if len(parts) != 3 {
		err = fmt.Errorf(`filter parameters must be in format "field||$condition||value"`)
		return
	}

	query.Field = parts[0]
	condition := parts[1]
	query.Value = strings.Split(parts[2], ",")

	matchedCondition := false
	for _, q := range conditions {
		if q.Query == condition {
			matchedCondition = true
			d, err := q.Factory(query.Value...)
			if err != nil {
				return query, err
			}
			query.Condition = bson.E{
				Key:   query.Field,
				Value: bson.D{d},
			}
			break
		}
	}

	if !matchedCondition {
		err = fmt.Errorf(`unknown filter condition: %s`, condition)
	}

	return query, err
}

func New(filterInput []string, orInput []string) (bson.D, error) {
	filter := bson.D{}

	for _, v := range filterInput {
		f, err := QueryFromString(v)
		if err != nil {
			return nil, err
		}

		filter = append(filter, f.Condition)
	}

	// @todo(sje): complete
	// for _, v := range orInput {
	// 	f, err := QueryFromString(v)
	// 	if err != nil {
	// 		return nil,  err
	// 	}
	// }

	return filter, nil
}
