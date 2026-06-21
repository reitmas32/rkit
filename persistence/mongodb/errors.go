package mongodb

import "github.com/reitmas32/rkit/core/kerrors"

var (
	ErrorItemIDRequired     = kerrors.NewKError("Item ID is required", 400, nil)
	ErrorItemNotFound       = kerrors.NewKError("Item not found", 404, nil)
	ErrorItemFieldsRequired = kerrors.NewKError("Item fields are required", 400, nil)
	ErrorDatabaseOperation  = kerrors.NewKError("Database operation failed", 500, nil)
	ErrorDuplicateKey       = kerrors.NewKError("Duplicate key error: a record with the same unique key already exists", 409, nil)
	ErrorPageableRequired   = kerrors.NewKError("Pageable is required and must be valid", 400, nil)
	ErrorConvertModelToDoc  = kerrors.NewKError("Error converting model to MongoDB document", 500, nil)
	ErrorDecodeDocument     = kerrors.NewKError("Error decoding MongoDB document", 500, nil)
	ErrorDecodeCursor       = kerrors.NewKError("Error iterating MongoDB cursor", 500, nil)

	// ErrorInvalidFieldName is returned when a filter/sort field is not a safe
	// identifier or is not permitted by the repository's FieldPolicy. The
	// rejected field and the reason are attached as metadata. This keeps a
	// client-controlled field name from being used as a MongoDB operator key.
	ErrorInvalidFieldName = kerrors.NewValidation("invalid or not-allowed filter/sort field", nil)
)
