package mongodb

import (
	"errors"

	"github.com/reitmas32/rkit/core/customctx"
	"github.com/reitmas32/rkit/core/kerrors"
	"github.com/reitmas32/rkit/core/result"
	"github.com/reitmas32/rkit/persistence/contracts"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *MongoRepository[E, M]) GetById(cc *customctx.CustomContext, id string) result.Result[E] {
	if id == "" {
		cc.Logger().Error(ErrorItemIDRequired.Error())
		cc.AddError(ErrorItemIDRequired)
		return result.Err[E](ErrorItemIDRequired)
	}

	// Decode into a raw map so we can use the JSON-based ModelToEntity conversion
	// without requiring BSON struct tags on the model.
	var raw map[string]interface{}
	err := r.Collection.FindOne(cc, bson.M{"id": id}).Decode(&raw)
	if err != nil {
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

	// Remove the MongoDB-managed _id field before converting to avoid confusion.
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

	return result.Ok(entityResult.Value())
}
