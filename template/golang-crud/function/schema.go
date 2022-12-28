package function

import "github.com/go-playground/validator/v10"

// Define the MongoDB schema
type Schema struct {
	Name         string `json:"name" bson:"name" validate:"required"`
	EmailAddress string `json:"emailAddress" bson:"emailAddress" validate:"required,email"`
}

// Set the database collection name
func (model *Schema) CollectionName() string {
	return "users"
}

// Database hooks https://github.com/Kamva/mgm#a-models-hooks
func (model *Schema) Creating() error {
	return nil
}

// Custom validation tags https://pkg.go.dev/github.com/go-playground/validator/v10#hdr-Custom_Validation_Functions
// Message uses universal translator https://pkg.go.dev/github.com/go-playground/universal-translator
// {0} is the field name, {1} is the value
var Validation = map[string]struct {
	Message   string
	Validator validator.Func
}{}
