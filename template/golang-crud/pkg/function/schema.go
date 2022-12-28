package function

import (
	"context"

	"github.com/google/uuid"
	"github.com/kamva/mgm/v3"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/function"
)

type Schema struct {
	mgm.DefaultModel `bson:",inline"`
	ID               string `json:"id" bson:"id"`
	function.Schema  `bson:",inline"`
}

func (model *Schema) Creating(ctx context.Context) error {
	// Call the DefaultModel Creating hook
	if err := model.DefaultModel.Creating(); err != nil {
		return err
	}

	//  Call any Creating hook on the schema
	if err := model.Schema.Creating(); err != nil {
		return err
	}

	model.ID = uuid.NewString()

	return nil
}
