package querybuilder

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

type QueryCondition struct {
	Query   string
	Factory func(...string) (bson.E, error)
}

// A lot of the conditions are basically the same - this generalises it
func passthru(query string, requiredArgs int) QueryCondition {
	return QueryCondition{
		Query: query,
		Factory: func(a ...string) (mongo bson.E, err error) {
			argLen := len(a)
			if argLen != requiredArgs {
				err = fmt.Errorf("condition requires %d argument(s), found %d", requiredArgs, argLen)
			}

			mongo = bson.E{
				Key:   query,
				Value: a[0],
			}

			return mongo, err
		},
	}
}

// @todo(sje) complete
var conditions []QueryCondition = []QueryCondition{
	passthru("$eq", 1),
	{
		Query: "$between",
	},
	{
		Query: "$cont",
	},
	{
		Query: "$excl",
	},
	{
		Query: "$ends",
	},
	{
		Query: "$isnull",
	},
	{
		Query: "$nin",
	},
	{
		Query: "$notnull",
		Factory: func(_ ...string) (bson.E, error) {
			return bson.E{
				Key:   "$ne",
				Value: nil,
			}, nil
		},
	},
	{
		Query: "$starts",
	},
}
