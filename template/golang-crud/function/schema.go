package function

// Define the MongoDB schema
type Schema struct {
	Name         string `json:"name" bson:"name"`
	EmailAddress string `json:"emailAddress" bson:"emailAddress"`
}

// Set the database collection name
func (model *Schema) CollectionName() string {
	return "users"
}

// Database hooks https://github.com/Kamva/mgm#a-models-hooks
func (model *Schema) Creating() error {
	return nil
}
