package mongodb

import (
	"errors"

	"github.com/reitmas32/rkit/core/customctx"
	"github.com/reitmas32/rkit/core/result"
	"github.com/reitmas32/rkit/persistence/contracts"
	"github.com/reitmas32/rkit/persistence/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *MongoRepository[E, M]) DeleteById(cc *customctx.CustomContext, id string) result.Result[E] {
	if id == "" {
		cc.Logger().Error(ErrorItemIDRequired.Error())
		cc.AddError(ErrorItemIDRequired)
		return result.Err[E](ErrorItemIDRequired)
	}

	// Fetch the document first to verify it exists and to fire the mutation hook.
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

	delete(raw, "_id")

	model, err := contracts.FromJSON[M](raw)
	if err != nil {
		kErr := ErrorDecodeDocument.WithCause(err)
		cc.Logger().Error(kErr.Error())
		cc.AddError(kErr)
		return result.Err[E](kErr)
	}

	// Notify mutation before deleting.
	mutationResult := models.NotifyMutation[E, M](cc, r.OnMutation, model, "delete")
	if !mutationResult.IsOk() {
		cc.Logger().Error(mutationResult.ToKError().Error())
	}

	if _, err := r.Collection.DeleteOne(cc, bson.M{"id": id}); err != nil {
		dbErr := ErrorDatabaseOperation.WithCause(err)
		cc.Logger().Error(dbErr.Error())
		cc.AddError(dbErr)
		return result.Err[E](dbErr)
	}

	return result.Empty[E]()
}
