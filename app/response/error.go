package response

import (
	"net/http"
)

var (
	UnexpectedServerError = CustomError{
		Message:  "Unexpected server error",
		Code:     10000,
		HTTPCode: http.StatusInternalServerError,
	}

	// RecordConflictError represents Duplicate entry for unique field error
	RecordConflictError = CustomError{
		Message:  "Record conflict",
		Code:     10202,
		HTTPCode: http.StatusConflict,
	}

	// DataNotUpdatableError represents Data not updatable error
	DataNotUpdatableError = CustomError{
		Message:  "Data not updatable",
		Code:     10210,
		HTTPCode: http.StatusNotAcceptable,
	}

	// ModelNotExistsError represents Model not found error
	ModelNotExistsError = CustomError{
		Message:  "Model not exists or has been deleted",
		Code:     10213,
		HTTPCode: http.StatusNotFound,
	}

	// InvalidParameterError represents Invalid parameter error
	InvalidParameterError = CustomError{
		Message:  "Invalid parameter",
		Code:     10111,
		HTTPCode: http.StatusUnprocessableEntity,
	}

	// BadRequestError represents bad request error
	BadRequestError = CustomError{
		Message:  "Bad Request",
		Code:     10005,
		HTTPCode: http.StatusBadRequest,
	}

	// InvalidFileTypeError represents Uploaded file type not supported
	InvalidFileTypeError = CustomError{
		Message:  "File type not supported",
		Code:     71001,
		HTTPCode: http.StatusUnsupportedMediaType,
	}

	// ParamCannotBeNullError represents Param cannot be null error
	ParamCannotBeNullError = CustomError{
		Message:  "cannot be null",
		Code:     10212,
		HTTPCode: http.StatusUnprocessableEntity,
	}

	//NotFoundError represents not found
	NotFoundError = CustomError{
		Message:  "Not Found",
		Code:     10213,
		HTTPCode: http.StatusNotFound,
	}
)
