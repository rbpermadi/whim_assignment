// Package response is used to write response based on Bukalapak API V4 standard
package response

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ResponseBody struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Errors  []ErrorInfo `json:"errors,omitempty"`
	Meta    MetaInfo    `json:"meta"`
}

// MetaInfo holds meta data
type MetaInfo struct {
	HTTPStatus int         `json:"http_status"`
	Offset     int         `json:"offset,omitempty"`
	Limit      int         `json:"limit,omitempty"`
	Total      int64       `json:"total,omitempty"`
	Sort       string      `json:"sort,omitempty"`
	Facets     interface{} `json:"facets,omitempty"`
}

// ErrorBody holds data for error response
type ErrorBody struct {
	Errors []ErrorInfo `json:"errors"`
	Meta   interface{} `json:"meta"`
}

// ErrorInfo holds error detail
type ErrorInfo struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Field   string `json:"field,omitempty"`
}

// CustomError holds data for customized error
type CustomError struct {
	Message  string
	Field    string
	Code     int
	HTTPCode int
}

// Error is a function to convert error to string.
// It exists to satisfy error interface
func (c CustomError) Error() string {
	return c.Message
}

// BuildSuccess is a function to create ResponseBody
func BuildSuccess(data interface{}, meta MetaInfo) ResponseBody {
	return ResponseBody{
		Data: data,
		Meta: meta,
	}
}

// BuildError is a function to create ErrorBody
func BuildError(errors []error) ErrorBody {
	var (
		ce CustomError
		ok bool
	)

	if len(errors) == 0 {
		ce = UnexpectedServerError
	} else {
		err := errors[0]
		ce, ok = err.(CustomError)
		if !ok {
			ce = UnexpectedServerError
		}
	}

	return ErrorBody{
		Errors: []ErrorInfo{
			{
				Message: ce.Message,
				Code:    ce.Code,
				Field:   ce.Field,
			},
		},
		Meta: MetaInfo{
			HTTPStatus: ce.HTTPCode,
		},
	}
}

// BuildErrors is a function to create ErrorBody
func BuildErrors(errors []error) ErrorBody {
	var (
		ce         CustomError
		ok         bool
		errorInfos []ErrorInfo
	)

	for _, err := range errors {
		ce, ok = err.(CustomError)
		if !ok {
			ce = UnexpectedServerError
		}

		errorInfo := ErrorInfo{
			Code:    ce.Code,
			Message: ce.Message,
		}

		errorInfos = append(errorInfos, errorInfo)
	}

	return ErrorBody{
		Errors: errorInfos,
		Meta: MetaInfo{
			HTTPStatus: ce.HTTPCode,
		},
	}
}

// BuildErrorAndStatus is a function to Differentiate Error and create Error Body and Response Status Code
func BuildErrorAndStatus(err error, fieldName string) (ErrorBody, int) {

	if strings.Contains(err.Error(), "strconv.ParseInt: parsing") ||
		strings.Contains(err.Error(), "strconv.Atoi: parsing") ||
		strings.Contains(err.Error(), "invalid character") ||
		strings.Contains(err.Error(), "json: cannot unmarshal") {
		return BuildError([]error{InvalidParameterError}), InvalidParameterError.HTTPCode
	} else if strings.Contains(err.Error(), "sql: Scan error") {
		return ErrorBody{
				Errors: []ErrorInfo{
					{
						Message: http.StatusText(http.StatusInternalServerError),
						Code:    http.StatusInternalServerError,
					},
				},
				Meta: MetaInfo{
					HTTPStatus: http.StatusInternalServerError,
				}},
			http.StatusInternalServerError
	} else if strings.Contains(err.Error(), "Duplicate entry") {
		ce := RecordConflictError

		return BuildError([]error{ce}), RecordConflictError.HTTPCode
	} else if strings.Contains(err.Error(), "cannot be null") {
		ce := ParamCannotBeNullError
		ce.Message = err.Error()

		return BuildError([]error{ce}), ParamCannotBeNullError.HTTPCode
	} else if strings.Contains(err.Error(), "Not Found") {
		ce := NotFoundError
		ce.Message = err.Error()

		return BuildError([]error{ce}), NotFoundError.HTTPCode
	} else if strings.Contains(err.Error(), "Bad Request") {
		ce := BadRequestError
		ce.Message = err.Error()

		return BuildError([]error{ce}), BadRequestError.HTTPCode
	}

	return BuildError([]error{UnexpectedServerError}), UnexpectedServerError.HTTPCode
}

// Write is a function to write data in json format
func Write(w http.ResponseWriter, result interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(result)
}
