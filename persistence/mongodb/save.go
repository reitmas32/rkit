package mongodb

import (
	"encoding/json"

	"github.com/reitmas32/rkit/core/customctx"
	"github.com/reitmas32/rkit/core/kerrors"
	"github.com/reitmas32/rkit/core/result"
	"github.com/reitmas32/rkit/persistence/contracts"
	"github.com/reitmas32/rkit/persistence/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *MongoRepository[E, M]) Save(cc *customctx.CustomContext, item E) result.Result[E] {
	if item.GetID() == "" {
		cc.Logger().Error(ErrorItemIDRequired.Error())
		cc.AddError(ErrorItemIDRequired)
		return result.Err[E](ErrorItemIDRequired)
	}

	modelResult := contracts.EntityToModel[E, M](item)
	if !modelResult.IsOk() {
		kErr := modelResult.ToKError()
		if kErr != nil {
			cc.Logger().Error(kErr.Error())
			cc.AddError(kErr)
			return result.Err[E](kErr)
		}
		genericErr := kerrors.NewKError("Error converting entity to model", 500, nil).WithCause(modelResult.Error())
		cc.Logger().Error(genericErr.Error())
		cc.AddError(genericErr)
		return result.Err[E](genericErr)
	}

	model := modelResult.Value()

	// Convert model to map via JSON so document keys match the model's JSON tags.
	// The entity ID is stored as the "id" field inside the document.
	data, err := json.Marshal(model)
	if err != nil {
		kErr := ErrorConvertModelToDoc.WithCause(err)
		cc.Logger().Error(kErr.Error())
		cc.AddError(kErr)
		return result.Err[E](kErr)
	}

	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		kErr := ErrorConvertModelToDoc.WithCause(err)
		cc.Logger().Error(kErr.Error())
		cc.AddError(kErr)
		return result.Err[E](kErr)
	}

	if _, err := r.Collection.InsertOne(cc, doc); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			dupErr := ErrorDuplicateKey.WithCause(err)
			cc.Logger().Error(dupErr.Error())
			cc.AddError(dupErr)
			return result.Err[E](dupErr)
		}
		dbErr := ErrorDatabaseOperation.WithCause(err)
		cc.Logger().Error(dbErr.Error())
		cc.AddError(dbErr)
		return result.Err[E](dbErr)
	}

	entityResult := contracts.ModelToEntity[E, M](model)
	if !entityResult.IsOk() {
		kErr := entityResult.ToKError()
		if kErr != nil {
			cc.Logger().Error(kErr.Error())
			cc.AddError(kErr)
			return result.Err[E](kErr)
		}
		genericErr := kerrors.NewKError("Error converting model to entity", 500, nil).WithCause(entityResult.Error())
		cc.Logger().Error(genericErr.Error())
		cc.AddError(genericErr)
		return result.Err[E](genericErr)
	}

	mutationResult := models.NotifyMutation[E, M](cc, r.OnMutation, model, "save")
	if !mutationResult.IsOk() {
		cc.Logger().Error(mutationResult.ToKError().Error())
	}

	return result.Ok(entityResult.Value())
}
