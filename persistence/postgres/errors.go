package postgres

import "github.com/reitmas32/rkit/core/kerrors"

var (
	ErrorItemIDRequired     = kerrors.NewKError("Item ID is required", 400, nil)
	ErrorItemNotFound       = kerrors.NewKError("Item not found", 404, nil)
	ErrorItemFieldsRequired = kerrors.NewKError("Item fields are required", 400, nil)
	ErrorDatabaseOperation  = kerrors.NewKError("Database operation failed", 500, nil)
	ErrorDuplicateKey       = kerrors.NewKError("Duplicate key error: a record with the same unique key already exists", 409, nil)
	ErrorPageableRequired   = kerrors.NewKError("Pageable is required and must be valid", 400, nil)
	ErrorConvertModelToJSON = kerrors.NewKError("Error converting model to JSON", 500, nil)
	ErrorConvertJSONToMap   = kerrors.NewKError("Error converting JSON to map", 500, nil)
	ErrorConvertMapToModel  = kerrors.NewKError("Error converting map to model", 500, nil)
)
