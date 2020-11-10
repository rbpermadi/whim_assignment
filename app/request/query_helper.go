package request

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// QueryHelper represent helper to get query url data
type QueryHelper struct {
	r  *http.Request
	uv url.Values
}

// GetString to get query url value with string data type, return empty string if query url not found
func (q *QueryHelper) GetString(p string, vars ...string) string {
	defValue := ""
	if len(vars) > 0 {
		defValue = vars[0]
	}

	sv := q.uv.Get(p)
	if sv != "" {
		return sv
	}
	return defValue
}

// GetInt to get query url value with integer data type, return 0 if query url not found
func (q *QueryHelper) GetInt(p string, vars ...int) int {
	defValue := 0
	if len(vars) > 0 {
		defValue = vars[0]
	}
	sv := q.uv.Get(p)
	if sv != "" {
		if v, err := strconv.Atoi(sv); err == nil {
			return v
		}
	}
	return defValue
}

// GetInt to get query url value with integer data type, return 0 if query url not found
func (q *QueryHelper) GetInt64(p string, vars ...int) int64 {
	defValue := 0
	if len(vars) > 0 {
		defValue = vars[0]
	}
	sv := q.uv.Get(p)
	if sv != "" {
		if v, err := strconv.Atoi(sv); err == nil {
			return int64(v)
		}
	}
	return int64(defValue)
}

// GetFloat to get query url value with float data type, return 0.0 if query url not found
func (q *QueryHelper) GetFloat(p string, vars ...float64) float64 {
	defValue := 0.0
	if len(vars) > 0 {
		defValue = vars[0]
	}
	sv := q.uv.Get(p)
	if sv != "" {
		if v, err := strconv.ParseFloat(sv, 64); err == nil {
			return v
		}
	}
	return defValue
}

// GetBool to get query url value with boolean data type, return false if query url not found
func (q *QueryHelper) GetBool(p string, vars ...bool) bool {
	defValue := false
	if len(vars) > 0 {
		defValue = vars[0]
	}
	sv := q.uv.Get(p)
	if sv != "" {
		v, err := strconv.ParseBool(sv)
		if err != nil {
			return defValue
		}
		return v
	}
	return defValue
}

// GetDate to get query url value with date data type using RFC3339 to parse string into time.Time, return nil if query url not found or parse format error
func (q *QueryHelper) GetDate(p string, vars ...time.Time) *time.Time {
	var defValue *time.Time
	if len(vars) > 0 {
		t := vars[0]
		defValue = &t
	}
	sv := q.uv.Get(p)
	if sv != "" {
		v, err := time.Parse(time.RFC3339, sv)
		if err != nil {
			fmt.Println("=====>", err)
			return defValue
		}

		return &v
	}

	return defValue
}

// NewQueryHelper is a function to create query helper struct
func NewQueryHelper(r *http.Request) *QueryHelper {
	return &QueryHelper{r, r.URL.Query()}
}
