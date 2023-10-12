package simple_query

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type FilterLogic string
type FilterOperator string

const (
	FilterLogicAnd FilterLogic = "and"
	FilterLogicOr  FilterLogic = "or"

	FilterOperatorEqual              FilterOperator = "equal"
	FilterOperatorNotEqual           FilterOperator = "not_equal"
	FilterOperatorGreaterThan        FilterOperator = "greater_than"
	FilterOperatorGreaterThanOrEqual FilterOperator = "greater_than_or_equal"
	FilterOperatorLessThan           FilterOperator = "less_than"
	FilterOperatorLessThanOrEqual    FilterOperator = "less_than_or_equal"
	FilterOperatorIsNull             FilterOperator = "is_null"
	FilterOperatorIsNotNull          FilterOperator = "is_not_null"
	FilterOperatorIn                 FilterOperator = "in"
	FilterOperatorNotIn              FilterOperator = "not_in"
	FilterOperatorLike               FilterOperator = "like"
	FilterOperatorNotLike            FilterOperator = "not_like"
)

var filterOperatorMap map[FilterOperator]string = map[FilterOperator]string{
	FilterOperatorEqual:              "=",
	FilterOperatorNotEqual:           "!=",
	FilterOperatorGreaterThan:        ">",
	FilterOperatorGreaterThanOrEqual: ">=",
	FilterOperatorLessThan:           "<",
	FilterOperatorLessThanOrEqual:    "<=",
	FilterOperatorIsNull:             "is null",
	FilterOperatorIsNotNull:          "is not null",
	FilterOperatorIn:                 "in",
	FilterOperatorNotIn:              "not in",
	FilterOperatorLike:               "like",
	FilterOperatorNotLike:            "not like",
}

type Filter struct {
	Logic    FilterLogic
	Field    string
	Operator FilterOperator
	Value    interface{}
	Filters  []*Filter
}

func NewFilter() *Filter {
	return &Filter{}
}

func (f *Filter) SetLogic(logic FilterLogic) *Filter {
	f.Logic = logic
	return f
}

func (f *Filter) SetCondition(field string, operator FilterOperator, value interface{}) *Filter {
	f.Field = field
	f.Operator = operator
	f.Value = value
	return f
}

func (f *Filter) AddFilter(field string, operator FilterOperator, value interface{}) *Filter {
	f.Filters = append(f.Filters, &Filter{Field: field, Operator: operator, Value: value})
	return f
}

func (f *Filter) AddFilters(filters ...*Filter) *Filter {
	f.Filters = append(f.Filters, filters...)
	return f
}

func (f *Filter) validate() error {
	var reflectValue reflect.Value = reflect.ValueOf(f.Value)

	if f.Logic != "" && f.Field != "" {
		return errors.New("field is not empty")
	}

	if f.Logic != "" && f.Operator != "" {
		return errors.New("operator is not empty")
	}

	if f.Logic != "" && (f.Value != nil || allowedKindValue[reflectValue.Kind()]) {
		return errors.New("value is not empty")
	}

	if f.Logic != "" && len(f.Filters) == 0 {
		return errors.New("filters is required")
	}

	if f.Logic == "" && len(f.Filters) > 0 {
		return errors.New("logic is required")
	}

	if f.Logic == "" && len(f.Filters) == 0 {
		if f.Field == "" {
			return errors.New("field is required")
		}

		if f.Operator == "" {
			return errors.New("operator is required")
		}

		if f.Operator != FilterOperatorIsNull && f.Operator != FilterOperatorIsNotNull {
			if f.Value == nil {
				return errors.New("value is required")
			}

			if !allowedKindValue[reflectValue.Kind()] {
				return fmt.Errorf("unsupported %s value type for operator %s", reflectValue.Kind().String(), f.Operator)
			}
		}

		if (f.Operator == FilterOperatorIsNotNull || f.Operator == FilterOperatorIsNull) && f.Value != nil {
			return errors.New("value is not empty")
		}

		if f.Operator != FilterOperatorIn && f.Operator != FilterOperatorNotIn && (reflectValue.Kind() == reflect.Slice || reflectValue.Kind() == reflect.Array) {
			return fmt.Errorf("unsupported %s value type for operator %s", reflectValue.Kind().String(), f.Operator)
		}

		if f.Operator == FilterOperatorIn || f.Operator == FilterOperatorNotIn {
			if reflectValue.Kind() != reflect.Slice && reflectValue.Kind() != reflect.Array {
				return fmt.Errorf("unsupported %s value type for operator %s", reflectValue.Kind().String(), f.Operator)
			}

			if reflectValue.Len() == 0 {
				return errors.New("value is required")
			}

			for i := 0; i < reflectValue.Len(); i++ {
				if !allowedKindValue[reflectValue.Index(i).Kind()] || reflectValue.Index(i).Kind() == reflect.Slice || reflectValue.Index(i).Kind() == reflect.Array {
					return fmt.Errorf("unsupported %s type of element value for operator %s", reflectValue.Index(i).Kind(), f.Operator)
				}
			}
		}

		if (f.Operator == FilterOperatorLike || f.Operator == FilterOperatorNotLike) && reflectValue.Kind() != reflect.String {
			return fmt.Errorf("unsupported %s type of value for operator %s", reflectValue.Kind().String(), f.Operator)
		}
	}

	for i := 0; i < len(f.Filters); i++ {
		var err error = f.Filters[i].validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *Filter) typedSliceToInterfaceSlice(value interface{}) ([]interface{}, error) {
	var (
		reflectValue   reflect.Value
		interfaceSlice []interface{}
	)

	reflectValue = reflect.ValueOf(value)
	if reflectValue.Kind() != reflect.Slice && reflectValue.Kind() != reflect.Array {
		return nil, fmt.Errorf("unsupported %s value type", reflectValue.Kind().String())
	}

	interfaceSlice = []interface{}{}
	for i := 0; i < reflectValue.Len(); i++ {
		if !allowedKindValue[reflectValue.Index(i).Kind()] || reflectValue.Index(i).Kind() == reflect.Slice || reflectValue.Index(i).Kind() == reflect.Array {
			return nil, fmt.Errorf("unsupported %s type of element value", reflectValue.Index(i).Kind().String())
		}
		interfaceSlice = append(interfaceSlice, reflectValue.Index(i).Interface())
	}

	return interfaceSlice, nil
}

func (f *Filter) toSQLWithArgs(dialect Dialect, args []interface{}, isRoot bool) (string, []interface{}, error) {
	var (
		whereClause          string
		conditionQueries     []string
		conditionQueryFormat string
		filterOperator       string
		placeholderStartIdx  int
		placeholderEndIdx    int
		placeholder          string
		conditionQuery       string
		err                  error
	)

	if dialect == "" {
		err = errors.New("dialect is required")
		return "", args, err
	}

	switch f.Operator {
	case FilterOperatorEqual, FilterOperatorNotEqual, FilterOperatorGreaterThan, FilterOperatorGreaterThanOrEqual, FilterOperatorLessThan, FilterOperatorLessThanOrEqual:
		conditionQueryFormat = "%s %s %s"
		filterOperator = filterOperatorMap[f.Operator]
		args = append(args, f.Value)
		placeholderStartIdx = len(args)
		placeholderEndIdx = len(args)
		placeholder = getPlaceholder(dialect, placeholderStartIdx, placeholderEndIdx)
		conditionQuery = fmt.Sprintf(conditionQueryFormat, f.Field, filterOperator, placeholder)

		return conditionQuery, args, nil

	case FilterOperatorIsNull, FilterOperatorIsNotNull:
		conditionQueryFormat = "%s %s"
		filterOperator = filterOperatorMap[f.Operator]
		conditionQuery = fmt.Sprintf(conditionQueryFormat, f.Field, filterOperator)

		return conditionQuery, args, nil

	case FilterOperatorIn, FilterOperatorNotIn:
		var interfaceSlice []interface{}

		conditionQueryFormat = "%s %s (%s)"
		filterOperator = filterOperatorMap[f.Operator]

		interfaceSlice, err = f.typedSliceToInterfaceSlice(f.Value)
		if err != nil {
			return "", nil, fmt.Errorf("%s for operator %s", err.Error(), f.Operator)
		}

		args = append(args, interfaceSlice...)
		placeholderStartIdx = len(args) - (len(interfaceSlice) - 1)
		placeholderEndIdx = len(args)
		placeholder = getPlaceholder(dialect, placeholderStartIdx, placeholderEndIdx)
		conditionQuery = fmt.Sprintf(conditionQueryFormat, f.Field, filterOperator, placeholder)

		return conditionQuery, args, nil

	case FilterOperatorLike, FilterOperatorNotLike:
		conditionQueryFormat = "%s %s concat('%%', %s, '%%')"

		switch dialect {
		case DialectMySQL:
			filterOperator = filterOperatorMap[f.Operator]
		case DialectPostgres:
			filterOperator = fmt.Sprintf("i%s", filterOperatorMap[FilterOperatorLike])
			if f.Operator == FilterOperatorNotLike {
				filterOperator = fmt.Sprintf("not i%s", filterOperatorMap[FilterOperatorLike])
			}
		}

		args = append(args, f.Value)
		placeholderStartIdx = len(args)
		placeholderEndIdx = len(args)
		placeholder = getPlaceholder(dialect, placeholderStartIdx, placeholderEndIdx)
		conditionQuery = fmt.Sprintf(conditionQueryFormat, f.Field, filterOperator, placeholder)

		return conditionQuery, args, nil
	}

	if len(f.Filters) == 0 {
		return "", args, nil
	}

	for i := 0; i < len(f.Filters); i++ {
		var (
			subConditionQuery string
			subArgs           []interface{}
		)

		if f.Filters[i] == nil {
			return "", args, nil
		}

		subConditionQuery, subArgs, err = f.Filters[i].toSQLWithArgs(dialect, args, false)
		if err != nil {
			return "", args, err
		}

		if subConditionQuery != "" {
			conditionQueries = append(conditionQueries, subConditionQuery)
		}

		args = subArgs
	}

	if len(conditionQueries) == 0 {
		return "", args, nil
	}

	whereClause = fmt.Sprintf("(%s)", strings.Join(conditionQueries, fmt.Sprintf(" %s ", f.Logic)))
	if isRoot {
		whereClause = strings.Join(conditionQueries, fmt.Sprintf(" %s ", f.Logic))
	}

	return whereClause, args, nil
}

func (f *Filter) ToSQLWithArgs(dialect Dialect, args []interface{}) (string, []interface{}, error) {
	var err error = f.validate()
	if err != nil {
		return "", args, err
	}

	return f.toSQLWithArgs(dialect, args, true)
}
