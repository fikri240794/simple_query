package simple_query

import (
	"fmt"
	"strings"
)

type Dialect string

const (
	DialectMySQL    Dialect = "mysql"
	DialectPostgres Dialect = "postgres"
)

var placeholderMap map[Dialect]string = map[Dialect]string{
	DialectMySQL:    "?",
	DialectPostgres: "$",
}

func getPlaceholder(dialect Dialect, startIdx, endIdx int) string {
	var placeholders []string = []string{}

	if startIdx <= 0 || endIdx <= 0 || endIdx < startIdx {
		return ""
	}

	switch dialect {
	case DialectMySQL:
		if startIdx == endIdx {
			return placeholderMap[dialect]
		}
		for i := startIdx; i <= endIdx; i++ {
			placeholders = append(placeholders, placeholderMap[dialect])
		}
		return strings.Join(placeholders, ", ")

	case DialectPostgres:
		if startIdx == endIdx {
			return fmt.Sprintf("%s%d", placeholderMap[dialect], endIdx)
		}
		for i := startIdx; i <= endIdx; i++ {
			placeholders = append(placeholders, fmt.Sprintf("%s%d", placeholderMap[dialect], i))
		}
		return strings.Join(placeholders, ", ")

	default:
		return ""
	}
}
