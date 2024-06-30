package simple_query

import (
	"fmt"
	"reflect"
	"strings"
)

type Filter struct {
	Logic    Logic
	Field    string
	Operator Operator
	Value    interface{}
	Filters  []*Filter
}

func NewFilter() *Filter {
	return &Filter{}
}

func (f *Filter) SetLogic(logic Logic) *Filter {
	f.Logic = logic
	return f
}

func (f *Filter) SetCondition(field string, operator Operator, value interface{}) *Filter {
	f.Field = field
	f.Operator = operator
	f.Value = value
	return f
}

func (f *Filter) AddFilter(field string, operator Operator, value interface{}) *Filter {
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
		return ErrFieldIsNotEmpty
	}

	if f.Logic != "" && f.Operator != "" {
		return ErrOperatorIsNotEmpty
	}

	if f.Logic != "" && (f.Value != nil || reflectValue.Kind() != reflect.Invalid) {
		return ErrValueIsNotEmpty
	}

	if f.Logic != "" && len(f.Filters) == 0 {
		return ErrFiltersIsRequired
	}

	if f.Logic == "" && len(f.Filters) > 0 {
		return ErrLogicIsRequired
	}

	if f.Logic == "" && len(f.Filters) == 0 {
		if f.Field == "" {
			return ErrFieldIsRequired
		}

		if f.Operator == "" {
			return ErrOperatorIsRequired
		}

		if f.Operator != OperatorIsNull && f.Operator != OperatorIsNotNull && f.Value == nil && reflectValue.Kind() == reflect.Invalid {
			return ErrValueIsRequired
		}

		if (f.Operator == OperatorIsNotNull || f.Operator == OperatorIsNull) && (f.Value != nil || reflectValue.Kind() != reflect.Invalid) {
			return ErrValueIsNotEmpty
		}

		if f.Operator != OperatorIn && f.Operator != OperatorNotIn && (reflectValue.Kind() == reflect.Slice || reflectValue.Kind() == reflect.Array) {
			return fmt.Errorf(ErrUnsupportedValueTypeForOperatorf, reflectValue.Kind().String(), f.Operator)
		}

		if f.Operator == OperatorIn || f.Operator == OperatorNotIn {
			if reflectValue.Kind() != reflect.Slice && reflectValue.Kind() != reflect.Array {
				return fmt.Errorf(ErrUnsupportedValueTypeForOperatorf, reflectValue.Kind().String(), f.Operator)
			}

			if reflectValue.Len() == 0 {
				return ErrValueIsRequired
			}
		}
	}

	for i := range f.Filters {
		var err error = f.Filters[i].validate()
		if err != nil {
			return err
		}
	}

	return nil
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
		err = ErrDialectIsRequired
		return "", args, err
	}

	switch f.Operator {
	case OperatorEqual, OperatorNotEqual, OperatorGreaterThan, OperatorGreaterThanOrEqual, OperatorLessThan, OperatorLessThanOrEqual:
		conditionQueryFormat = "%s %s %s"
		filterOperator = filterOperatorMap[f.Operator]
		args = append(args, f.Value)
		placeholderStartIdx = len(args)
		placeholderEndIdx = len(args)
		placeholder = getPlaceholder(dialect, placeholderStartIdx, placeholderEndIdx)
		conditionQuery = fmt.Sprintf(conditionQueryFormat, f.Field, filterOperator, placeholder)

		return conditionQuery, args, nil

	case OperatorIsNull, OperatorIsNotNull:
		conditionQueryFormat = "%s %s"
		filterOperator = filterOperatorMap[f.Operator]
		conditionQuery = fmt.Sprintf(conditionQueryFormat, f.Field, filterOperator)

		return conditionQuery, args, nil

	case OperatorIn, OperatorNotIn:
		var interfaceSlice []interface{}

		conditionQueryFormat = "%s %s (%s)"
		filterOperator = filterOperatorMap[f.Operator]

		interfaceSlice, err = typedSliceToInterfaceSlice(f.Value)
		if err != nil {
			err = fmt.Errorf(ErrForOperatorf, err.Error(), f.Operator)
			return "", nil, err
		}

		args = append(args, interfaceSlice...)
		placeholderStartIdx = len(args) - (len(interfaceSlice) - 1)
		placeholderEndIdx = len(args)
		placeholder = getPlaceholder(dialect, placeholderStartIdx, placeholderEndIdx)
		conditionQuery = fmt.Sprintf(conditionQueryFormat, f.Field, filterOperator, placeholder)

		return conditionQuery, args, nil

	case OperatorLike, OperatorNotLike:
		conditionQueryFormat = "%s %s concat('%%', %s, '%%')"

		switch dialect {
		case DialectMySQL:
			filterOperator = filterOperatorMap[f.Operator]
		case DialectPostgres:
			filterOperator = fmt.Sprintf("i%s", filterOperatorMap[OperatorLike])
			if f.Operator == OperatorNotLike {
				filterOperator = fmt.Sprintf("not i%s", filterOperatorMap[OperatorLike])
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

	for i := range f.Filters {
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
