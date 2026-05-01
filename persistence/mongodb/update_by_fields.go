package mongodb

import (
	"errors"

	"github.com/reitmas32/rkit/core/customctx"
	"github.com/reitmas32/rkit/core/kerrors"
	"github.com/reitmas32/rkit/core/result"
	"github.com/reitmas32/rkit/persistence/contracts"
	"github.com/reitmas32/rkit/persistence/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *MongoRepository[E, M]) UpdateByFields(cc *customctx.CustomContext, id string, fields map[string]any) result.Result[E] {
	if id == "" {
		cc.Logger().Error(ErrorItemIDRequired.Error())
		cc.AddError(ErrorItemIDRequired)
		return result.Err[E](ErrorItemIDRequired)
	}

	if len(fields) == 0 {
		cc.Logger().Error(ErrorItemFieldsRequired.Error())
		cc.AddError(ErrorItemFieldsRequired)
		return result.Err[E](ErrorItemFieldsRequired)
	}

	filter := bson.M{"id": id}

	// Verify the document exists before updating.
	var check map[string]interface{}
	if err := r.Collection.FindOne(cc, filter).Decode(&check); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			cc.Logger().Error(ErrorItemNotFound.Error())
			cc.AddError(ErrorItemNotFound)
			return result.Err[E](ErrorItemNotFound)
		}
		dbErr := ErrorDatabaseOperation.WithCause(err)
		cc.Logger().Error(dbErr.Error())
		cc.AddError(dbErr)
		return result.Err[E](dbErr)
	}

	update := bson.M{"$set": fields}
	if _, err := r.Collection.UpdateOne(cc, filter, update); err != nil {
		dbErr := ErrorDatabaseOperation.WithCause(err)
		cc.Logger().Error(dbErr.Error())
		cc.AddError(dbErr)
		return result.Err[E](dbErr)
	}

	// Fetch the updated document and return it as entity.
	var raw map[string]interface{}
	if err := r.Collection.FindOne(cc, filter).Decode(&raw); err != nil {
		dbErr := ErrorDatabaseOperation.WithCause(err)
		cc.Logger().Error(dbErr.Error())
		cc.AddError(dbErr)
		return result.Err[E](dbErr)
	}

	delete(raw, "_id")

	model, err := contracts.FromJSON[M](raw)
	if err != nil {
		kErr := ErrorDecodeDocument.WithCause(err)
		cc.Logger().Error(kErr.Error())
		cc.AddError(kErr)
		return result.Err[E](kErr)
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

	mutationResult := models.NotifyMutation[E, M](cc, r.OnMutation, model, "update")
	if !mutationResult.IsOk() {
		cc.Logger().Error(mutationResult.ToKError().Error())
	}

	return result.Ok(entityResult.Value())
}
