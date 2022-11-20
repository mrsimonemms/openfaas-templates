package querybuilder

import "go.mongodb.org/mongo-driver/bson"

type Queries struct {
	Query    string
	Operator string // Operators use built-in MongoDB commands
	Factory  func() (bson.E, error)
}

var queryBuilder []Queries = []Queries{
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
		Query:    "$nin",
		Operator: "$notin",
	},
	{
		Query: "$notnull",
		Factory: func() (bson.E, error) {
			return bson.E{
				Key: "$ne",
				Value: nil,
			}, nil
		},
	},
	{
		Query: "$starts",
	},
}
