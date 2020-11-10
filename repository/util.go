package repository

import (
	"fmt"
	"strings"
	"time"
)

const (
	MysqlTimeFormat = "2006-01-02 15:04:05"
)

func buildQuery(query string, args ...interface{}) string {
	return fmt.Sprintf(query, args...)
}

func sqlTime(t time.Time) string {
	return t.Format(MysqlTimeFormat)
}

func int64sToString(a []int64, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}
