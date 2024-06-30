package simple_query

import "errors"

type Dialect string

const (
	DialectMySQL    Dialect = "mysql"
	DialectPostgres Dialect = "postgres"
)

var placeholderMap map[Dialect]string = map[Dialect]string{
	DialectMySQL:    "?",
	DialectPostgres: "$",
}

type Logic string
type Operator string

const (
	LogicAnd Logic = "and"
	LogicOr  Logic = "or"

	OperatorEqual              Operator = "equal"
	OperatorNotEqual           Operator = "not_equal"
	OperatorGreaterThan        Operator = "greater_than"
	OperatorGreaterThanOrEqual Operator = "greater_than_or_equal"
	OperatorLessThan           Operator = "less_than"
	OperatorLessThanOrEqual    Operator = "less_than_or_equal"
	OperatorIsNull             Operator = "is_null"
	OperatorIsNotNull          Operator = "is_not_null"
	OperatorIn                 Operator = "in"
	OperatorNotIn              Operator = "not_in"
	OperatorLike               Operator = "like"
	OperatorNotLike            Operator = "not_like"
)

var filterOperatorMap map[Operator]string = map[Operator]string{
	OperatorEqual:              "=",
	OperatorNotEqual:           "!=",
	OperatorGreaterThan:        ">",
	OperatorGreaterThanOrEqual: ">=",
	OperatorLessThan:           "<",
	OperatorLessThanOrEqual:    "<=",
	OperatorIsNull:             "is null",
	OperatorIsNotNull:          "is not null",
	OperatorIn:                 "in",
	OperatorNotIn:              "not in",
	OperatorLike:               "like",
	OperatorNotLike:            "not like",
}

var (
	ErrUnsupportedValueTypef            string = "unsupported %s value type"
	ErrForOperatorf                     string = "%s for operator %s"
	ErrUnsupportedValueTypeForOperatorf string = "unsupported %s value type for operator %s"
)

var (
	ErrFieldsIsRequired                    error = errors.New("fields is required")
	ErrFieldIsRequired                     error = errors.New("field is required")
	ErrTableIsRequired                     error = errors.New("table is required")
	ErrFieldIsNotEmpty                     error = errors.New("field is not empty")
	ErrOperatorIsNotEmpty                  error = errors.New("operator is not empty")
	ErrValueIsNotEmpty                     error = errors.New("value is not empty")
	ErrFiltersIsRequired                   error = errors.New("filters is required")
	ErrLogicIsRequired                     error = errors.New("logic is required")
	ErrOperatorIsRequired                  error = errors.New("operator is required")
	ErrValueIsRequired                     error = errors.New("value is required")
	ErrDialectIsRequired                   error = errors.New("dialect is required")
	ErrValuesIsRequired                    error = errors.New("values is required")
	ErrValueLengthIsNotEqualToFieldsLength error = errors.New("value length is not equal to fields length")
	ErrFilterIsRequired                    error = errors.New("filter is required")
)
