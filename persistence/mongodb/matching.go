package mongodb

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/reitmas32/rkit/core/customctx"
	"github.com/reitmas32/rkit/core/kerrors"
	"github.com/reitmas32/rkit/core/result"
	"github.com/reitmas32/rkit/persistence/contracts"
	"github.com/reitmas32/rkit/persistence/criteria"
	"github.com/reitmas32/rkit/persistence/pagination"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// MatchingResult holds paginated data together with the total count of matching documents.
type MatchingResult[E contracts.IEntity] struct {
	Data  []E
	Total int64
}

// Matching returns paginated entities that match the given criteria.
func (r *MongoRepository[E, M]) Matching(
	cc *customctx.CustomContext,
	crit criteria.Criteria,
	pageable *pagination.Pageable,
) result.Result[[]E] {
	if pageable == nil {
		pageable = pagination.NewPageableWithoutSort(0, 10)
	}

	if !pageable.IsValid() {
		cc.Logger().Error(ErrorPageableRequired.Error())
		cc.AddError(ErrorPageableRequired)
		return result.Err[[]E](ErrorPageableRequired)
	}

	if vErr := r.validateCriteria(crit, pageable); vErr != nil {
		cc.Logger().Error(vErr.Detail())
		cc.AddError(vErr)
		return result.Err[[]E](vErr)
	}

	filter := buildFilter(crit.Filters)
	opts := buildFindOptions(pageable)

	cursor, err := r.Collection.Find(cc, filter, opts)
	if err != nil {
		dbErr := ErrorDatabaseOperation.WithCause(err)
		cc.Logger().Error(dbErr.Error())
		cc.AddError(dbErr)
		return result.Err[[]E](dbErr)
	}
	defer cursor.Close(cc)

	entities, kErr := decodeCursor[E, M](cc, cursor)
	if kErr != nil {
		return result.Err[[]E](kErr)
	}

	return result.Ok(entities)
}

// MatchingWithTotal returns paginated entities and the total count of matching documents.
// Use this when you need pagination metadata (total pages, total items).
func (r *MongoRepository[E, M]) MatchingWithTotal(
	cc *customctx.CustomContext,
	crit criteria.Criteria,
	pageable *pagination.Pageable,
) result.Result[MatchingResult[E]] {
	if pageable == nil {
		pageable = pagination.NewPageableWithoutSort(0, 10)
	}

	if !pageable.IsValid() {
		cc.Logger().Error(ErrorPageableRequired.Error())
		cc.AddError(ErrorPageableRequired)
		return result.Err[MatchingResult[E]](ErrorPageableRequired)
	}

	if vErr := r.validateCriteria(crit, pageable); vErr != nil {
		cc.Logger().Error(vErr.Detail())
		cc.AddError(vErr)
		return result.Err[MatchingResult[E]](vErr)
	}

	filter := buildFilter(crit.Filters)

	total, err := r.Collection.CountDocuments(cc, filter)
	if err != nil {
		dbErr := ErrorDatabaseOperation.WithCause(err)
		cc.Logger().Error(dbErr.Error())
		cc.AddError(dbErr)
		return result.Err[MatchingResult[E]](dbErr)
	}

	opts := buildFindOptions(pageable)

	cursor, err := r.Collection.Find(cc, filter, opts)
	if err != nil {
		dbErr := ErrorDatabaseOperation.WithCause(err)
		cc.Logger().Error(dbErr.Error())
		cc.AddError(dbErr)
		return result.Err[MatchingResult[E]](dbErr)
	}
	defer cursor.Close(cc)

	entities, kErr := decodeCursor[E, M](cc, cursor)
	if kErr != nil {
		return result.Err[MatchingResult[E]](kErr)
	}

	return result.Ok(MatchingResult[E]{Data: entities, Total: total})
}

// buildFilter converts criteria.Filters into a MongoDB bson.D filter document.
func buildFilter(filters criteria.Filters) bson.D {
	filter := bson.D{}

	for _, f := range filters.Get() {
		field := string(f.Field)
		value := f.Value

		var condition bson.M

		switch f.Operator {
		case criteria.OperatorEqual:
			condition = bson.M{"$eq": value}
		case criteria.OperatorNotEqual:
			condition = bson.M{"$ne": value}
		case criteria.OperatorGreaterThan:
			condition = bson.M{"$gt": value}
		case criteria.OperatorGreaterEqual:
			condition = bson.M{"$gte": value}
		case criteria.OperatorLessThan:
			condition = bson.M{"$lt": value}
		case criteria.OperatorLessEqual:
			condition = bson.M{"$lte": value}
		case criteria.OperatorLike:
			pattern := likeToRegex(fmt.Sprintf("%v", value))
			condition = bson.M{"$regex": pattern, "$options": "i"}
		case criteria.OperatorNotLike:
			pattern := likeToRegex(fmt.Sprintf("%v", value))
			condition = bson.M{"$not": bson.M{"$regex": pattern, "$options": "i"}}
		case criteria.OperatorIn:
			condition = bson.M{"$in": value}
		case criteria.OperatorNotIn:
			condition = bson.M{"$nin": value}
		default:
			continue
		}

		filter = append(filter, bson.E{Key: field, Value: condition})
	}

	return filter
}

// buildFindOptions converts pagination and sort settings into MongoDB FindOptions.
func buildFindOptions(pageable *pagination.Pageable) *options.FindOptionsBuilder {
	opts := options.Find().
		SetSkip(int64(pageable.Offset())).
		SetLimit(int64(pageable.Limit()))

	if pageable.Sort != nil && pageable.Sort.IsValid() {
		direction := 1
		if strings.ToUpper(string(pageable.Sort.Direction)) == "DESC" {
			direction = -1
		}
		opts.SetSort(bson.D{{Key: pageable.Sort.Field, Value: direction}})
	}

	return opts
}

// decodeCursor iterates a MongoDB cursor and converts each document to a domain entity.
// cc implements context.Context so it is accepted by the mongo.Cursor methods directly.
func decodeCursor[E contracts.IEntity, M contracts.IModel](
	cc *customctx.CustomContext,
	cursor interface {
		Next(context.Context) bool
		Decode(interface{}) error
		Close(context.Context) error
	},
) ([]E, *kerrors.KError) {
	entities := make([]E, 0)

	for cursor.Next(cc) {
		var raw map[string]interface{}
		if err := cursor.Decode(&raw); err != nil {
			return nil, ErrorDecodeCursor.WithCause(err)
		}

		delete(raw, "_id")

		model, err := contracts.FromJSON[M](raw)
		if err != nil {
			return nil, ErrorDecodeDocument.WithCause(err)
		}

		entityResult := contracts.ModelToEntity[E, M](model)
		if !entityResult.IsOk() {
			kErr := entityResult.ToKError()
			if kErr != nil {
				return nil, kErr
			}
			return nil, kerrors.NewKError("Error converting model to entity", 500, nil)
		}

		entities = append(entities, entityResult.Value())
	}

	return entities, nil
}

// likeToRegex converts a SQL LIKE pattern (%foo%) to a MongoDB regex pattern (.*foo.*).
func likeToRegex(pattern string) string {
	// Escape all regex metacharacters first so a user-supplied value cannot
	// inject a regular expression (e.g. a catastrophic-backtracking ReDoS that
	// would stall MongoDB). QuoteMeta leaves the SQL LIKE wildcards %/_ intact,
	// so they can then be translated to their regex equivalents.
	escaped := regexp.QuoteMeta(pattern)
	escaped = strings.ReplaceAll(escaped, "%", ".*")
	escaped = strings.ReplaceAll(escaped, "_", ".")
	return escaped
}
