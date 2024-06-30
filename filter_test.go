package simple_query

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func testFilter_FilterEquality(t *testing.T, expectation, actual *Filter) {
	if expectation.Logic != actual.Logic {
		t.Errorf("expectation logic is %s, got %s", expectation.Logic, actual.Logic)
	}

	if expectation.Field != actual.Field {
		t.Errorf("expectation field is %s, got %s", expectation.Field, actual.Field)
	}

	if expectation.Operator != actual.Operator {
		t.Errorf("expectation operator is %s, got %s", expectation.Operator, actual.Operator)
	}

	if !deepEqual(expectation.Value, actual.Value) {
		t.Errorf("expectation value is %v, got %v", expectation.Value, actual.Value)
	}

	if len(expectation.Filters) != len(actual.Filters) {
		t.Errorf("expectation length of filters is %d, got %d", len(expectation.Filters), len(actual.Filters))
	}

	if len(expectation.Filters) > 0 {
		for i := 0; i < len(expectation.Filters); i++ {
			testFilter_FilterEquality(t, expectation.Filters[i], actual.Filters[i])
		}
	}
}

func TestFilter_NewFilter(t *testing.T) {
	testFilter_FilterEquality(t, &Filter{}, NewFilter())
}

func TestFilter_SetLogic(t *testing.T) {
	var testCases []struct {
		Name        string
		Logic       Logic
		Expectation *Filter
	} = []struct {
		Name        string
		Logic       Logic
		Expectation *Filter
	}{
		{
			Name:  "logic and",
			Logic: LogicAnd,
			Expectation: &Filter{
				Logic: LogicAnd,
			},
		},
		{
			Name:  "logic or",
			Logic: LogicOr,
			Expectation: &Filter{
				Logic: LogicOr,
			},
		},
	}

	for i := 0; i < len(testCases); i++ {
		t.Run(testCases[i].Name, func(t *testing.T) {
			var actual *Filter = NewFilter().
				SetLogic(testCases[i].Logic)
			testFilter_FilterEquality(t, testCases[i].Expectation, actual)
		})
	}
}

func TestFilter_SetCondition(t *testing.T) {
	var testCases []struct {
		Name        string
		Field       string
		Operator    Operator
		Value       interface{}
		Expectation *Filter
	} = []struct {
		Name        string
		Field       string
		Operator    Operator
		Value       interface{}
		Expectation *Filter
	}{
		{
			Name:     fmt.Sprintf("operator %s", OperatorEqual),
			Field:    "field1",
			Operator: OperatorEqual,
			Value:    "value1",
			Expectation: &Filter{
				Field:    "field1",
				Operator: OperatorEqual,
				Value:    "value1",
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", OperatorNotEqual),
			Field:    "field1",
			Operator: OperatorNotEqual,
			Value:    true,
			Expectation: &Filter{
				Field:    "field1",
				Operator: OperatorNotEqual,
				Value:    true,
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", OperatorGreaterThan),
			Field:    "field1",
			Operator: OperatorGreaterThan,
			Value:    int64(100),
			Expectation: &Filter{
				Field:    "field1",
				Operator: OperatorGreaterThan,
				Value:    int64(100),
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", OperatorGreaterThanOrEqual),
			Field:    "field1",
			Operator: OperatorGreaterThanOrEqual,
			Value:    float64(100),
			Expectation: &Filter{
				Field:    "field1",
				Operator: OperatorGreaterThanOrEqual,
				Value:    float64(100),
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", OperatorLessThan),
			Field:    "field1",
			Operator: OperatorLessThan,
			Value:    uint64(100),
			Expectation: &Filter{
				Field:    "field1",
				Operator: OperatorLessThan,
				Value:    uint64(100),
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", OperatorLessThanOrEqual),
			Field:    "field1",
			Operator: OperatorLessThanOrEqual,
			Value:    "2006-01-02T15:04:05+07:00",
			Expectation: &Filter{
				Field:    "field1",
				Operator: OperatorLessThanOrEqual,
				Value:    "2006-01-02T15:04:05+07:00",
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", OperatorIsNull),
			Field:    "field1",
			Operator: OperatorIsNull,
			Value:    nil,
			Expectation: &Filter{
				Field:    "field1",
				Operator: OperatorIsNull,
				Value:    nil,
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", OperatorIsNotNull),
			Field:    "field1",
			Operator: OperatorIsNotNull,
			Value:    nil,
			Expectation: &Filter{
				Field:    "field1",
				Operator: OperatorIsNotNull,
				Value:    nil,
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", OperatorIn),
			Field:    "field1",
			Operator: OperatorIn,
			Value:    []string{"value1", "value 2", "value3"},
			Expectation: &Filter{
				Field:    "field1",
				Operator: OperatorIn,
				Value:    []string{"value1", "value 2", "value3"},
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", OperatorNotIn),
			Field:    "field1",
			Operator: OperatorNotIn,
			Value:    [3]int64{1, 2, 3},
			Expectation: &Filter{
				Field:    "field1",
				Operator: OperatorNotIn,
				Value:    [3]int64{1, 2, 3},
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", OperatorLike),
			Field:    "field1",
			Operator: OperatorLike,
			Value:    "value1",
			Expectation: &Filter{
				Field:    "field1",
				Operator: OperatorLike,
				Value:    "value1",
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", OperatorNotLike),
			Field:    "field1",
			Operator: OperatorNotLike,
			Value:    "value1",
			Expectation: &Filter{
				Field:    "field1",
				Operator: OperatorNotLike,
				Value:    "value1",
			},
		},
	}

	for i := 0; i < len(testCases); i++ {
		t.Run(testCases[i].Name, func(t *testing.T) {
			var actual *Filter = NewFilter().
				SetCondition(
					testCases[i].Field,
					testCases[i].Operator,
					testCases[i].Value,
				)

			testFilter_FilterEquality(t, testCases[i].Expectation, actual)
		})
	}
}

func TestFilter_AddFilter(t *testing.T) {
	var testCases []struct {
		Name        string
		Field       string
		Operator    Operator
		Value       interface{}
		Expectation *Filter
	} = []struct {
		Name        string
		Field       string
		Operator    Operator
		Value       interface{}
		Expectation *Filter
	}{
		{
			Name:     fmt.Sprintf("operator %s", OperatorEqual),
			Field:    "field1",
			Operator: OperatorEqual,
			Value:    "value1",
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorEqual,
						Value:    "value1",
					},
				},
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", OperatorNotEqual),
			Field:    "field1",
			Operator: OperatorNotEqual,
			Value:    true,
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorNotEqual,
						Value:    true,
					},
				},
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", OperatorGreaterThan),
			Field:    "field1",
			Operator: OperatorGreaterThan,
			Value:    int64(100),
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorGreaterThan,
						Value:    int64(100),
					},
				},
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", OperatorGreaterThanOrEqual),
			Field:    "field1",
			Operator: OperatorGreaterThanOrEqual,
			Value:    float64(100),
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorGreaterThanOrEqual,
						Value:    float64(100),
					},
				},
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", OperatorLessThan),
			Field:    "field1",
			Operator: OperatorLessThan,
			Value:    uint64(100),
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorLessThan,
						Value:    uint64(100),
					},
				},
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", OperatorLessThanOrEqual),
			Field:    "field1",
			Operator: OperatorLessThanOrEqual,
			Value:    "2006-01-02T15:04:05+07:00",
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorLessThanOrEqual,
						Value:    "2006-01-02T15:04:05+07:00",
					},
				},
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", OperatorIsNull),
			Field:    "field1",
			Operator: OperatorIsNull,
			Value:    nil,
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorIsNull,
						Value:    nil,
					},
				},
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", OperatorIsNotNull),
			Field:    "field1",
			Operator: OperatorIsNotNull,
			Value:    nil,
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorIsNotNull,
						Value:    nil,
					},
				},
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", OperatorIn),
			Field:    "field1",
			Operator: OperatorIn,
			Value:    []string{"value1", "value 2", "value3"},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorIn,
						Value:    []string{"value1", "value 2", "value3"},
					},
				},
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", OperatorNotIn),
			Field:    "field1",
			Operator: OperatorNotIn,
			Value:    [3]int64{1, 2, 3},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorNotIn,
						Value:    [3]int64{1, 2, 3},
					},
				},
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", OperatorLike),
			Field:    "field1",
			Operator: OperatorLike,
			Value:    "value1",
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorLike,
						Value:    "value1",
					},
				},
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", OperatorNotLike),
			Field:    "field1",
			Operator: OperatorNotLike,
			Value:    "value1",
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorNotLike,
						Value:    "value1",
					},
				},
			},
		},
	}

	for i := 0; i < len(testCases); i++ {
		t.Run(testCases[i].Name, func(t *testing.T) {
			var actual *Filter = NewFilter().
				AddFilter(
					testCases[i].Field,
					testCases[i].Operator,
					testCases[i].Value,
				)

			testFilter_FilterEquality(t, testCases[i].Expectation, actual)
		})
	}
}

func TestFilter_AddFilters(t *testing.T) {
	var testCases []struct {
		Name        string
		Filter      *Filter
		Expectation *Filter
	} = []struct {
		Name        string
		Filter      *Filter
		Expectation *Filter
	}{
		{
			Name: fmt.Sprintf("operator %s", OperatorEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorEqual,
				Value:    "value1",
			},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorEqual,
						Value:    "value1",
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("operator %s", OperatorNotEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorNotEqual,
				Value:    true,
			},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorNotEqual,
						Value:    true,
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("operator %s", OperatorGreaterThan),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorGreaterThan,
				Value:    int64(100),
			},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorGreaterThan,
						Value:    int64(100),
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("operator %s", OperatorGreaterThanOrEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorGreaterThanOrEqual,
				Value:    float64(100),
			},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorGreaterThanOrEqual,
						Value:    float64(100),
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("operator %s", OperatorLessThan),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorLessThan,
				Value:    uint64(100),
			},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorLessThan,
						Value:    uint64(100),
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("operator %s", OperatorLessThanOrEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorLessThanOrEqual,
				Value:    "2006-01-02T15:04:05+07:00",
			},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorLessThanOrEqual,
						Value:    "2006-01-02T15:04:05+07:00",
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("operator %s", OperatorIsNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorIsNull,
				Value:    nil,
			},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorIsNull,
						Value:    nil,
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("operator %s", OperatorIsNotNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorIsNotNull,
				Value:    nil,
			},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorIsNotNull,
						Value:    nil,
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("operator %s", OperatorIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorIn,
				Value:    []string{"value1", "value 2", "value3"},
			},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorIn,
						Value:    []string{"value1", "value 2", "value3"},
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("operator %s", OperatorNotIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorNotIn,
				Value:    [3]int64{1, 2, 3},
			},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorNotIn,
						Value:    [3]int64{1, 2, 3},
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("operator %s", OperatorLike),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorLike,
				Value:    "value1",
			},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorLike,
						Value:    "value1",
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("operator %s", OperatorNotLike),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorNotLike,
				Value:    "value1",
			},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorNotLike,
						Value:    "value1",
					},
				},
			},
		},
	}

	for i := 0; i < len(testCases); i++ {
		t.Run(testCases[i].Name, func(t *testing.T) {
			var actual *Filter = NewFilter().
				AddFilters(testCases[i].Filter)

			testFilter_FilterEquality(t, testCases[i].Expectation, actual)
		})
	}
}

func TestFilter_validate(t *testing.T) {
	var testCases []struct {
		Name        string
		Filter      *Filter
		Expectation error
	} = []struct {
		Name        string
		Filter      *Filter
		Expectation error
	}{
		{
			Name: "logic is not empty and field is not empty",
			Filter: &Filter{
				Logic: LogicAnd,
				Field: "field1",
			},
			Expectation: ErrFieldIsNotEmpty,
		},
		{
			Name: "logic is not empty and operator is not empty",
			Filter: &Filter{
				Logic:    LogicOr,
				Operator: OperatorEqual,
			},
			Expectation: ErrOperatorIsNotEmpty,
		},
		{
			Name: "logic is not empty and value is not nil or value kind is allowed",
			Filter: &Filter{
				Logic: LogicAnd,
				Value: "value1",
			},
			Expectation: ErrValueIsNotEmpty,
		},
		{
			Name: "logic is not empty and filters length is zero",
			Filter: &Filter{
				Logic:   LogicAnd,
				Filters: []*Filter{},
			},
			Expectation: ErrFiltersIsRequired,
		},
		{
			Name: "logic is empty and filters length greater than zero",
			Filter: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorEqual,
						Value:    "value1",
					},
				},
			},
			Expectation: ErrLogicIsRequired,
		},
		{
			Name: "logic is empty and filters length is zero and field is empty",
			Filter: &Filter{
				Operator: OperatorEqual,
				Value:    "valu1",
			},
			Expectation: ErrFieldIsRequired,
		},
		{
			Name: "logic is empty and filters length is zero and operator is empty",
			Filter: &Filter{
				Field: "field1",
				Value: "value1",
			},
			Expectation: ErrOperatorIsRequired,
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is not %s and operator is not %s and value is nil", OperatorIsNull, OperatorIsNotNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorEqual,
				Value:    nil,
			},
			Expectation: ErrValueIsRequired,
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is %s and value is not nil", OperatorIsNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorIsNull,
				Value:    "value1",
			},
			Expectation: errors.New("value is not empty"),
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is %s and value is not nil", OperatorIsNotNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorIsNotNull,
				Value:    "value1",
			},
			Expectation: errors.New("value is not empty"),
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is not %s and operator is not %s and value kind is %s", OperatorIn, OperatorNotIn, reflect.Slice.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorEqual,
				Value:    []int64{1, 2, 3},
			},
			Expectation: fmt.Errorf(ErrUnsupportedValueTypeForOperatorf, reflect.Slice.String(), OperatorEqual),
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is not %s and operator is not %s and value kind is %s", OperatorIn, OperatorNotIn, reflect.Array.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorEqual,
				Value:    [3]string{"value1", "value2", "value3"},
			},
			Expectation: fmt.Errorf(ErrUnsupportedValueTypeForOperatorf, reflect.Array.String(), OperatorEqual),
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is %s and value kind is not %s and %s", OperatorIn, reflect.Slice.String(), reflect.Array.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorIn,
				Value:    int64(123),
			},
			Expectation: fmt.Errorf(ErrUnsupportedValueTypeForOperatorf, reflect.Int64.String(), OperatorIn),
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is %s and value length is zero", OperatorIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorIn,
				Value:    []int64{},
			},
			Expectation: errors.New("value is required"),
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is %s and value kind is not allowed", OperatorNotIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorNotIn,
				Value:    int64(123),
			},
			Expectation: fmt.Errorf(ErrUnsupportedValueTypeForOperatorf, reflect.Int64.String(), OperatorNotIn),
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is %s and value length is zero", OperatorNotIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorNotIn,
				Value:    []int64{},
			},
			Expectation: errors.New("value is required"),
		},
		{
			Name: "filter is valid",
			Filter: &Filter{
				Logic: LogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorEqual,
						Value:    int64(123),
					},
					{
						Field:    "field2",
						Operator: OperatorEqual,
						Value:    "value1",
					},
				},
			},
			Expectation: nil,
		},
		{
			Name: "filter is invalid",
			Filter: &Filter{
				Logic: LogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorEqual,
						Value:    int64(123),
					},
					{
						Field:    "field2",
						Operator: OperatorEqual,
						Value:    []string{"a", "b", "c"},
					},
				},
			},
			Expectation: fmt.Errorf(ErrUnsupportedValueTypeForOperatorf, reflect.Slice.String(), OperatorEqual),
		},
	}

	for i := 0; i < len(testCases); i++ {
		t.Run(testCases[i].Name, func(t *testing.T) {
			var actual error = testCases[i].Filter.validate()

			if testCases[i].Expectation != nil && actual == nil {
				t.Error("expectation error is not nil, got nil")
			}

			if testCases[i].Expectation == nil && actual != nil {
				t.Error("expectation error is nil, got not nil")
			}

			if testCases[i].Expectation != nil && actual != nil && testCases[i].Expectation.Error() != actual.Error() {
				t.Errorf("expectation error is %s, got %s", testCases[i].Expectation.Error(), actual.Error())
			}
		})
	}
}

func TestFilter_toSQLWithArgs(t *testing.T) {
	var testCases []struct {
		Name        string
		Filter      *Filter
		Dialect     Dialect
		Args        []interface{}
		IsRoot      bool
		Expectation struct {
			ConditionQuery string
			Args           []interface{}
			Error          error
		}
	} = []struct {
		Name        string
		Filter      *Filter
		Dialect     Dialect
		Args        []interface{}
		IsRoot      bool
		Expectation struct {
			ConditionQuery string
			Args           []interface{}
			Error          error
		}
	}{
		{
			Name:    "dialect is empty",
			Filter:  &Filter{},
			Dialect: "",
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: "",
				Args:           []interface{}{},
				Error:          ErrDialectIsRequired,
			},
		},

		// MYSQL
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectMySQL, OperatorEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorEqual,
				Value:    "value1",
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[OperatorEqual], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectMySQL, OperatorEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorEqual,
				Value:    "value1",
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[OperatorEqual], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectMySQL, OperatorNotEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorNotEqual,
				Value:    "value1",
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[OperatorNotEqual], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectMySQL, OperatorNotEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorNotEqual,
				Value:    "value1",
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[OperatorNotEqual], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectMySQL, OperatorGreaterThan),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorGreaterThan,
				Value:    "value1",
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[OperatorGreaterThan], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectMySQL, OperatorGreaterThan),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorGreaterThan,
				Value:    "value1",
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[OperatorGreaterThan], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectMySQL, OperatorGreaterThanOrEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorGreaterThanOrEqual,
				Value:    "value1",
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[OperatorGreaterThanOrEqual], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectMySQL, OperatorGreaterThanOrEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorGreaterThanOrEqual,
				Value:    "value1",
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[OperatorGreaterThanOrEqual], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectMySQL, OperatorLessThan),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorLessThan,
				Value:    "value1",
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[OperatorLessThan], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectMySQL, OperatorLessThan),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorLessThan,
				Value:    "value1",
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[OperatorLessThan], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectMySQL, OperatorLessThanOrEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorLessThanOrEqual,
				Value:    "value1",
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[OperatorLessThanOrEqual], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectMySQL, OperatorLessThanOrEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorLessThanOrEqual,
				Value:    "value1",
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[OperatorLessThanOrEqual], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectMySQL, OperatorIsNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorIsNull,
				Value:    nil,
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s", "field1", filterOperatorMap[OperatorIsNull]),
				Args:           []interface{}{},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectMySQL, OperatorIsNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorIsNull,
				Value:    nil,
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s", "field1", filterOperatorMap[OperatorIsNull]),
				Args:           []interface{}{"other args value"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectMySQL, OperatorIsNotNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorIsNotNull,
				Value:    nil,
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s", "field1", filterOperatorMap[OperatorIsNotNull]),
				Args:           []interface{}{},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectMySQL, OperatorIsNotNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorIsNotNull,
				Value:    nil,
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s", "field1", filterOperatorMap[OperatorIsNotNull]),
				Args:           []interface{}{"other args value"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectMySQL, OperatorIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorIn,
				Value:    []string{"value1", "value2", "value3"},
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s (%s)", "field1", filterOperatorMap[OperatorIn], fmt.Sprintf("%s, %s, %s", placeholderMap[DialectMySQL], placeholderMap[DialectMySQL], placeholderMap[DialectMySQL])),
				Args:           []interface{}{"value1", "value2", "value3"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectMySQL, OperatorIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorIn,
				Value:    []string{"value1", "value2", "value3"},
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s (%s)", "field1", filterOperatorMap[OperatorIn], fmt.Sprintf("%s, %s, %s", placeholderMap[DialectMySQL], placeholderMap[DialectMySQL], placeholderMap[DialectMySQL])),
				Args:           []interface{}{"other args value", "value1", "value2", "value3"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s with value kind is not %s and %s", DialectMySQL, OperatorIn, reflect.Slice.String(), reflect.Array.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorIn,
				Value:    "value1",
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: "",
				Args:           []interface{}{},
				Error:          fmt.Errorf(ErrUnsupportedValueTypeForOperatorf, reflect.String.String(), OperatorIn),
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectMySQL, OperatorNotIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorNotIn,
				Value:    []string{"value1", "value2", "value3"},
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s (%s)", "field1", filterOperatorMap[OperatorNotIn], fmt.Sprintf("%s, %s, %s", placeholderMap[DialectMySQL], placeholderMap[DialectMySQL], placeholderMap[DialectMySQL])),
				Args:           []interface{}{"value1", "value2", "value3"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectMySQL, OperatorNotIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorNotIn,
				Value:    []string{"value1", "value2", "value3"},
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s (%s)", "field1", filterOperatorMap[OperatorNotIn], fmt.Sprintf("%s, %s, %s", placeholderMap[DialectMySQL], placeholderMap[DialectMySQL], placeholderMap[DialectMySQL])),
				Args:           []interface{}{"other args value", "value1", "value2", "value3"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s with value kind is not %s and %s", DialectMySQL, OperatorNotIn, reflect.Slice.String(), reflect.Array.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorNotIn,
				Value:    "value1",
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: "",
				Args:           []interface{}{},
				Error:          fmt.Errorf(ErrUnsupportedValueTypeForOperatorf, reflect.String.String(), OperatorNotIn),
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectMySQL, OperatorLike),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorLike,
				Value:    "value1",
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s concat('%%', %s, '%%')", "field1", filterOperatorMap[OperatorLike], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectMySQL, OperatorLike),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorLike,
				Value:    "value1",
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s concat('%%', %s, '%%')", "field1", filterOperatorMap[OperatorLike], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectMySQL, OperatorNotLike),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorNotLike,
				Value:    "value1",
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s concat('%%', %s, '%%')", "field1", filterOperatorMap[OperatorNotLike], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectMySQL, OperatorNotLike),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorNotLike,
				Value:    "value1",
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s concat('%%', %s, '%%')", "field1", filterOperatorMap[OperatorNotLike], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name:    fmt.Sprintf("dialect %s with filters length is zero", DialectMySQL),
			Filter:  &Filter{},
			Dialect: DialectMySQL,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: "",
				Args:           []interface{}{},
				Error:          nil,
			},
		},
		{
			Name:    fmt.Sprintf("dialect %s and args length greater than zero with filters length is zero", DialectMySQL),
			Filter:  &Filter{},
			Dialect: DialectMySQL,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: "",
				Args:           []interface{}{"other args value"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s", DialectMySQL),
			Filter: &Filter{
				Logic: LogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorEqual,
						Value:    "value1",
					},
					{
						Logic: LogicOr,
						Filters: []*Filter{
							{
								Field:    "field2",
								Operator: OperatorIsNull,
								Value:    nil,
							},
							{
								Field:    "field3",
								Operator: OperatorIn,
								Value:    []int64{1, 2, 3},
							},
						},
					},
					{
						Field:    "field4",
						Operator: OperatorLike,
						Value:    "value4",
					},
				},
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf(
					"(field1 %s %s %s (field2 %s %s field3 %s %s) %s field4 %s %s)",
					filterOperatorMap[OperatorEqual],
					placeholderMap[DialectMySQL],
					LogicAnd,
					filterOperatorMap[OperatorIsNull],
					LogicOr,
					filterOperatorMap[OperatorIn],
					fmt.Sprintf("(%s, %s, %s)", placeholderMap[DialectMySQL], placeholderMap[DialectMySQL], placeholderMap[DialectMySQL]),
					LogicAnd,
					filterOperatorMap[OperatorLike],
					fmt.Sprintf("concat('%%', %s, '%%')", placeholderMap[DialectMySQL]),
				),
				Args:  []interface{}{"value1", 1, 2, 3, "value4"},
				Error: nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with element filters is empty filter", DialectMySQL),
			Filter: &Filter{
				Logic: LogicAnd,
				Filters: []*Filter{
					{},
					{},
					{},
				},
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: "",
				Args:           []interface{}{},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with element filters is nil", DialectMySQL),
			Filter: &Filter{
				Logic: LogicAnd,
				Filters: []*Filter{
					nil,
					nil,
					nil,
				},
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: "",
				Args:           []interface{}{},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with element filters is invalid", DialectMySQL),
			Filter: &Filter{
				Logic: LogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorIn,
						Value:    "value1",
					},
				},
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: "",
				Args:           []interface{}{},
				Error:          fmt.Errorf(ErrUnsupportedValueTypeForOperatorf, reflect.String.String(), OperatorIn),
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero", DialectMySQL),
			Filter: &Filter{
				Logic: LogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorEqual,
						Value:    "value1",
					},
					{
						Logic: LogicOr,
						Filters: []*Filter{
							{
								Field:    "field2",
								Operator: OperatorIsNull,
								Value:    nil,
							},
							{
								Field:    "field3",
								Operator: OperatorIn,
								Value:    []int64{1, 2, 3},
							},
						},
					},
					{
						Field:    "field4",
						Operator: OperatorLike,
						Value:    "value4",
					},
				},
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf(
					"(field1 %s %s %s (field2 %s %s field3 %s %s) %s field4 %s %s)",
					filterOperatorMap[OperatorEqual],
					placeholderMap[DialectMySQL],
					LogicAnd,
					filterOperatorMap[OperatorIsNull],
					LogicOr,
					filterOperatorMap[OperatorIn],
					fmt.Sprintf("(%s, %s, %s)", placeholderMap[DialectMySQL], placeholderMap[DialectMySQL], placeholderMap[DialectMySQL]),
					LogicAnd,
					filterOperatorMap[OperatorLike],
					fmt.Sprintf("concat('%%', %s, '%%')", placeholderMap[DialectMySQL]),
				),
				Args:  []interface{}{"other args value", "value1", 1, 2, 3, "value4"},
				Error: nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with isRoot is true", DialectMySQL),
			Filter: &Filter{
				Logic: LogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorEqual,
						Value:    "value1",
					},
					{
						Logic: LogicOr,
						Filters: []*Filter{
							{
								Field:    "field2",
								Operator: OperatorIsNull,
								Value:    nil,
							},
							{
								Field:    "field3",
								Operator: OperatorIn,
								Value:    []int64{1, 2, 3},
							},
						},
					},
					{
						Field:    "field4",
						Operator: OperatorLike,
						Value:    "value4",
					},
				},
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{},
			IsRoot:  true,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf(
					"field1 %s %s %s (field2 %s %s field3 %s %s) %s field4 %s %s",
					filterOperatorMap[OperatorEqual],
					placeholderMap[DialectMySQL],
					LogicAnd,
					filterOperatorMap[OperatorIsNull],
					LogicOr,
					filterOperatorMap[OperatorIn],
					fmt.Sprintf("(%s, %s, %s)", placeholderMap[DialectMySQL], placeholderMap[DialectMySQL], placeholderMap[DialectMySQL]),
					LogicAnd,
					filterOperatorMap[OperatorLike],
					fmt.Sprintf("concat('%%', %s, '%%')", placeholderMap[DialectMySQL]),
				),
				Args:  []interface{}{"value1", 1, 2, 3, "value4"},
				Error: nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero and isRoot is true", DialectMySQL),
			Filter: &Filter{
				Logic: LogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorEqual,
						Value:    "value1",
					},
					{
						Logic: LogicOr,
						Filters: []*Filter{
							{
								Field:    "field2",
								Operator: OperatorIsNull,
								Value:    nil,
							},
							{
								Field:    "field3",
								Operator: OperatorIn,
								Value:    []int64{1, 2, 3},
							},
						},
					},
					{
						Field:    "field4",
						Operator: OperatorLike,
						Value:    "value4",
					},
				},
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{"other args value"},
			IsRoot:  true,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf(
					"field1 %s %s %s (field2 %s %s field3 %s %s) %s field4 %s %s",
					filterOperatorMap[OperatorEqual],
					placeholderMap[DialectMySQL],
					LogicAnd,
					filterOperatorMap[OperatorIsNull],
					LogicOr,
					filterOperatorMap[OperatorIn],
					fmt.Sprintf("(%s, %s, %s)", placeholderMap[DialectMySQL], placeholderMap[DialectMySQL], placeholderMap[DialectMySQL]),
					LogicAnd,
					filterOperatorMap[OperatorLike],
					fmt.Sprintf("concat('%%', %s, '%%')", placeholderMap[DialectMySQL]),
				),
				Args:  []interface{}{"other args value", "value1", 1, 2, 3, "value4"},
				Error: nil,
			},
		},

		// POSTGRES
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectPostgres, OperatorEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorEqual,
				Value:    "value1",
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[OperatorEqual], fmt.Sprintf("%s1", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectPostgres, OperatorEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorEqual,
				Value:    "value1",
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[OperatorEqual], fmt.Sprintf("%s2", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectPostgres, OperatorNotEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorNotEqual,
				Value:    "value1",
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[OperatorNotEqual], fmt.Sprintf("%s1", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectPostgres, OperatorNotEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorNotEqual,
				Value:    "value1",
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[OperatorNotEqual], fmt.Sprintf("%s2", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectPostgres, OperatorGreaterThan),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorGreaterThan,
				Value:    "value1",
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[OperatorGreaterThan], fmt.Sprintf("%s1", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectPostgres, OperatorGreaterThan),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorGreaterThan,
				Value:    "value1",
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[OperatorGreaterThan], fmt.Sprintf("%s2", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectPostgres, OperatorGreaterThanOrEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorGreaterThanOrEqual,
				Value:    "value1",
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[OperatorGreaterThanOrEqual], fmt.Sprintf("%s1", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectPostgres, OperatorGreaterThanOrEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorGreaterThanOrEqual,
				Value:    "value1",
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[OperatorGreaterThanOrEqual], fmt.Sprintf("%s2", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectPostgres, OperatorLessThan),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorLessThan,
				Value:    "value1",
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[OperatorLessThan], fmt.Sprintf("%s1", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectPostgres, OperatorLessThan),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorLessThan,
				Value:    "value1",
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[OperatorLessThan], fmt.Sprintf("%s2", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectPostgres, OperatorLessThanOrEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorLessThanOrEqual,
				Value:    "value1",
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[OperatorLessThanOrEqual], fmt.Sprintf("%s1", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectPostgres, OperatorLessThanOrEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorLessThanOrEqual,
				Value:    "value1",
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[OperatorLessThanOrEqual], fmt.Sprintf("%s2", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectPostgres, OperatorIsNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorIsNull,
				Value:    nil,
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s", "field1", filterOperatorMap[OperatorIsNull]),
				Args:           []interface{}{},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectPostgres, OperatorIsNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorIsNull,
				Value:    nil,
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s", "field1", filterOperatorMap[OperatorIsNull]),
				Args:           []interface{}{"other args value"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectPostgres, OperatorIsNotNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorIsNotNull,
				Value:    nil,
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s", "field1", filterOperatorMap[OperatorIsNotNull]),
				Args:           []interface{}{},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectPostgres, OperatorIsNotNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorIsNotNull,
				Value:    nil,
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s", "field1", filterOperatorMap[OperatorIsNotNull]),
				Args:           []interface{}{"other args value"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectPostgres, OperatorIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorIn,
				Value:    []string{"value1", "value2", "value3"},
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s (%s)", "field1", filterOperatorMap[OperatorIn], fmt.Sprintf("%s1, %s2, %s3", placeholderMap[DialectPostgres], placeholderMap[DialectPostgres], placeholderMap[DialectPostgres])),
				Args:           []interface{}{"value1", "value2", "value3"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectPostgres, OperatorIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorIn,
				Value:    []string{"value1", "value2", "value3"},
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s (%s)", "field1", filterOperatorMap[OperatorIn], fmt.Sprintf("%s2, %s3, %s4", placeholderMap[DialectPostgres], placeholderMap[DialectPostgres], placeholderMap[DialectPostgres])),
				Args:           []interface{}{"other args value", "value1", "value2", "value3"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s with value kind is not %s and %s", DialectPostgres, OperatorIn, reflect.Slice.String(), reflect.Array.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorIn,
				Value:    "value1",
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: "",
				Args:           []interface{}{},
				Error:          fmt.Errorf(ErrUnsupportedValueTypeForOperatorf, reflect.String.String(), OperatorIn),
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectPostgres, OperatorNotIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorNotIn,
				Value:    []string{"value1", "value2", "value3"},
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s (%s)", "field1", filterOperatorMap[OperatorNotIn], fmt.Sprintf("%s1, %s2, %s3", placeholderMap[DialectPostgres], placeholderMap[DialectPostgres], placeholderMap[DialectPostgres])),
				Args:           []interface{}{"value1", "value2", "value3"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectPostgres, OperatorNotIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorNotIn,
				Value:    []string{"value1", "value2", "value3"},
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s (%s)", "field1", filterOperatorMap[OperatorNotIn], fmt.Sprintf("%s2, %s3, %s4", placeholderMap[DialectPostgres], placeholderMap[DialectPostgres], placeholderMap[DialectPostgres])),
				Args:           []interface{}{"other args value", "value1", "value2", "value3"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s with value kind is not %s and %s", DialectPostgres, OperatorNotIn, reflect.Slice.String(), reflect.Array.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorNotIn,
				Value:    "value1",
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: "",
				Args:           []interface{}{},
				Error:          fmt.Errorf(ErrUnsupportedValueTypeForOperatorf, reflect.String.String(), OperatorNotIn),
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectPostgres, OperatorLike),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorLike,
				Value:    "value1",
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s concat('%%', %s, '%%')", "field1", fmt.Sprintf("i%s", filterOperatorMap[OperatorLike]), fmt.Sprintf("%s1", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectPostgres, OperatorLike),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorLike,
				Value:    "value1",
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s concat('%%', %s, '%%')", "field1", fmt.Sprintf("i%s", filterOperatorMap[OperatorLike]), fmt.Sprintf("%s2", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectPostgres, OperatorNotLike),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorNotLike,
				Value:    "value1",
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s concat('%%', %s, '%%')", "field1", fmt.Sprintf("not i%s", filterOperatorMap[OperatorLike]), fmt.Sprintf("%s1", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectPostgres, OperatorNotLike),
			Filter: &Filter{
				Field:    "field1",
				Operator: OperatorNotLike,
				Value:    "value1",
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf("%s %s concat('%%', %s, '%%')", "field1", fmt.Sprintf("not i%s", filterOperatorMap[OperatorLike]), fmt.Sprintf("%s2", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name:    fmt.Sprintf("dialect %s with filters length is zero", DialectPostgres),
			Filter:  &Filter{},
			Dialect: DialectPostgres,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: "",
				Args:           []interface{}{},
				Error:          nil,
			},
		},
		{
			Name:    fmt.Sprintf("dialect %s and args length greater than zero with filters length is zero", DialectPostgres),
			Filter:  &Filter{},
			Dialect: DialectPostgres,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: "",
				Args:           []interface{}{"other args value"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s", DialectPostgres),
			Filter: &Filter{
				Logic: LogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorEqual,
						Value:    "value1",
					},
					{
						Logic: LogicOr,
						Filters: []*Filter{
							{
								Field:    "field2",
								Operator: OperatorIsNull,
								Value:    nil,
							},
							{
								Field:    "field3",
								Operator: OperatorIn,
								Value:    []int64{1, 2, 3},
							},
						},
					},
					{
						Field:    "field4",
						Operator: OperatorLike,
						Value:    "value4",
					},
				},
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf(
					"(field1 %s %s %s (field2 %s %s field3 %s %s) %s field4 %s %s)",
					filterOperatorMap[OperatorEqual],
					fmt.Sprintf("%s1", placeholderMap[DialectPostgres]),
					LogicAnd,
					filterOperatorMap[OperatorIsNull],
					LogicOr,
					filterOperatorMap[OperatorIn],
					fmt.Sprintf("(%s2, %s3, %s4)", placeholderMap[DialectPostgres], placeholderMap[DialectPostgres], placeholderMap[DialectPostgres]),
					LogicAnd,
					fmt.Sprintf("i%s", filterOperatorMap[OperatorLike]),
					fmt.Sprintf("concat('%%', %s5, '%%')", placeholderMap[DialectPostgres]),
				),
				Args:  []interface{}{"value1", 1, 2, 3, "value4"},
				Error: nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with element filters is empty filter", DialectPostgres),
			Filter: &Filter{
				Logic: LogicAnd,
				Filters: []*Filter{
					{},
					{},
					{},
				},
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: "",
				Args:           []interface{}{},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with element filters is nil", DialectPostgres),
			Filter: &Filter{
				Logic: LogicAnd,
				Filters: []*Filter{
					nil,
					nil,
					nil,
				},
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: "",
				Args:           []interface{}{},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with element filters is invalid", DialectPostgres),
			Filter: &Filter{
				Logic: LogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorIn,
						Value:    "value1",
					},
				},
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: "",
				Args:           []interface{}{},
				Error:          fmt.Errorf(ErrUnsupportedValueTypeForOperatorf, reflect.String.String(), OperatorIn),
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero", DialectPostgres),
			Filter: &Filter{
				Logic: LogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorEqual,
						Value:    "value1",
					},
					{
						Logic: LogicOr,
						Filters: []*Filter{
							{
								Field:    "field2",
								Operator: OperatorIsNull,
								Value:    nil,
							},
							{
								Field:    "field3",
								Operator: OperatorIn,
								Value:    []int64{1, 2, 3},
							},
						},
					},
					{
						Field:    "field4",
						Operator: OperatorLike,
						Value:    "value4",
					},
				},
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{"other args value"},
			IsRoot:  false,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf(
					"(field1 %s %s %s (field2 %s %s field3 %s %s) %s field4 %s %s)",
					filterOperatorMap[OperatorEqual],
					fmt.Sprintf("%s2", placeholderMap[DialectPostgres]),
					LogicAnd,
					filterOperatorMap[OperatorIsNull],
					LogicOr,
					filterOperatorMap[OperatorIn],
					fmt.Sprintf("(%s3, %s4, %s5)", placeholderMap[DialectPostgres], placeholderMap[DialectPostgres], placeholderMap[DialectPostgres]),
					LogicAnd,
					fmt.Sprintf("i%s", filterOperatorMap[OperatorLike]),
					fmt.Sprintf("concat('%%', %s6, '%%')", placeholderMap[DialectPostgres]),
				),
				Args:  []interface{}{"other args value", "value1", 1, 2, 3, "value4"},
				Error: nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with isRoot is true", DialectPostgres),
			Filter: &Filter{
				Logic: LogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorEqual,
						Value:    "value1",
					},
					{
						Logic: LogicOr,
						Filters: []*Filter{
							{
								Field:    "field2",
								Operator: OperatorIsNull,
								Value:    nil,
							},
							{
								Field:    "field3",
								Operator: OperatorIn,
								Value:    []int64{1, 2, 3},
							},
						},
					},
					{
						Field:    "field4",
						Operator: OperatorLike,
						Value:    "value4",
					},
				},
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{},
			IsRoot:  true,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf(
					"field1 %s %s %s (field2 %s %s field3 %s %s) %s field4 %s %s",
					filterOperatorMap[OperatorEqual],
					fmt.Sprintf("%s1", placeholderMap[DialectPostgres]),
					LogicAnd,
					filterOperatorMap[OperatorIsNull],
					LogicOr,
					filterOperatorMap[OperatorIn],
					fmt.Sprintf("(%s2, %s3, %s4)", placeholderMap[DialectPostgres], placeholderMap[DialectPostgres], placeholderMap[DialectPostgres]),
					LogicAnd,
					fmt.Sprintf("i%s", filterOperatorMap[OperatorLike]),
					fmt.Sprintf("concat('%%', %s5, '%%')", placeholderMap[DialectPostgres]),
				),
				Args:  []interface{}{"value1", 1, 2, 3, "value4"},
				Error: nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero and isRoot is true", DialectPostgres),
			Filter: &Filter{
				Logic: LogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorEqual,
						Value:    "value1",
					},
					{
						Logic: LogicOr,
						Filters: []*Filter{
							{
								Field:    "field2",
								Operator: OperatorIsNull,
								Value:    nil,
							},
							{
								Field:    "field3",
								Operator: OperatorIn,
								Value:    []int64{1, 2, 3},
							},
						},
					},
					{
						Field:    "field4",
						Operator: OperatorLike,
						Value:    "value4",
					},
				},
			},
			Dialect: DialectPostgres,
			Args:    []interface{}{"other args value"},
			IsRoot:  true,
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: fmt.Sprintf(
					"field1 %s %s %s (field2 %s %s field3 %s %s) %s field4 %s %s",
					filterOperatorMap[OperatorEqual],
					fmt.Sprintf("%s2", placeholderMap[DialectPostgres]),
					LogicAnd,
					filterOperatorMap[OperatorIsNull],
					LogicOr,
					filterOperatorMap[OperatorIn],
					fmt.Sprintf("(%s3, %s4, %s5)", placeholderMap[DialectPostgres], placeholderMap[DialectPostgres], placeholderMap[DialectPostgres]),
					LogicAnd,
					fmt.Sprintf("i%s", filterOperatorMap[OperatorLike]),
					fmt.Sprintf("concat('%%', %s6, '%%')", placeholderMap[DialectPostgres]),
				),
				Args:  []interface{}{"other args value", "value1", 1, 2, 3, "value4"},
				Error: nil,
			},
		},
	}

	for i := 0; i < len(testCases); i++ {
		t.Run(testCases[i].Name, func(t *testing.T) {
			var (
				actualConditionQuery string
				actualArgs           []interface{}
				actualErr            error
			)

			actualConditionQuery, actualArgs, actualErr = testCases[i].Filter.toSQLWithArgs(testCases[i].Dialect, testCases[i].Args, testCases[i].IsRoot)

			if testCases[i].Expectation.Error != nil && actualErr == nil {
				t.Error("expectation error is not nil, got nil")
			}

			if testCases[i].Expectation.Error == nil && actualErr != nil {
				t.Error("expectation error is nil, got not nil")
			}

			if testCases[i].Expectation.Error != nil && actualErr != nil && testCases[i].Expectation.Error.Error() != actualErr.Error() {
				t.Errorf("expectation error is %s, got %s", testCases[i].Expectation.Error.Error(), actualErr.Error())
			}

			if testCases[i].Expectation.Error == nil && actualErr == nil {
				if testCases[i].Expectation.ConditionQuery != actualConditionQuery {
					t.Errorf("expectation conditional query is %s, got %s", testCases[i].Expectation.ConditionQuery, actualConditionQuery)
				}

				if len(testCases[i].Expectation.Args) != len(actualArgs) {
					t.Errorf("expectation args lenght is %d, got %d", len(testCases[i].Expectation.Args), len(actualArgs))
				}

				for x := 0; x < len(testCases[i].Expectation.Args); x++ {
					if !deepEqual(testCases[i].Expectation.Args[x], actualArgs[x]) {
						t.Errorf("expectation element of args is %v, got %v", testCases[i].Expectation.Args[x], actualArgs[x])
					}
				}
			}
		})
	}
}

func TestFilter_ToSQLWithArgs(t *testing.T) {
	var testCases []struct {
		Name        string
		Filter      *Filter
		Dialect     Dialect
		Args        []interface{}
		Expectation struct {
			ConditionQuery string
			Args           []interface{}
			Error          error
		}
	} = []struct {
		Name        string
		Filter      *Filter
		Dialect     Dialect
		Args        []interface{}
		Expectation struct {
			ConditionQuery string
			Args           []interface{}
			Error          error
		}
	}{
		{
			Name: "invalid validation",
			Filter: &Filter{
				Logic: LogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorEqual,
						Value:    []string{"a", "b", "c"},
					},
				},
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{},
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: "",
				Args:           []interface{}{},
				Error:          fmt.Errorf(ErrUnsupportedValueTypeForOperatorf, reflect.Slice.String(), OperatorEqual),
			},
		},
		{
			Name: "filter is valid",
			Filter: &Filter{
				Logic: LogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: OperatorEqual,
						Value:    "value1",
					},
				},
			},
			Dialect: DialectMySQL,
			Args:    []interface{}{},
			Expectation: struct {
				ConditionQuery string
				Args           []interface{}
				Error          error
			}{
				ConditionQuery: "field1 = ?",
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
	}

	for i := 0; i < len(testCases); i++ {
		t.Run(testCases[i].Name, func(t *testing.T) {
			var (
				actualConditionQuery string
				actualArgs           []interface{}
				actualErr            error
			)

			actualConditionQuery, actualArgs, actualErr = testCases[i].Filter.ToSQLWithArgs(testCases[i].Dialect, testCases[i].Args)

			if testCases[i].Expectation.Error != nil && actualErr == nil {
				t.Error("expectation error is not nil, got nil")
			}

			if testCases[i].Expectation.Error == nil && actualErr != nil {
				t.Error("expectation error is nil, got not nil")
			}

			if testCases[i].Expectation.Error != nil && actualErr != nil && testCases[i].Expectation.Error.Error() != actualErr.Error() {
				t.Errorf("expectation error is %s, got %s", testCases[i].Expectation.Error.Error(), actualErr.Error())
			}

			if testCases[i].Expectation.Error == nil && actualErr == nil {
				if testCases[i].Expectation.ConditionQuery != actualConditionQuery {
					t.Errorf("expectation conditional query is %s, got %s", testCases[i].Expectation.ConditionQuery, actualConditionQuery)
				}

				if len(testCases[i].Expectation.Args) != len(actualArgs) {
					t.Errorf("expectation args lenght is %d, got %d", len(testCases[i].Expectation.Args), len(actualArgs))
				}

				for x := 0; x < len(testCases[i].Expectation.Args); x++ {
					if !deepEqual(testCases[i].Expectation.Args[x], actualArgs[x]) {
						t.Errorf("expectation element of args is %v, got %v", testCases[i].Expectation.Args[x], actualArgs[x])
					}
				}
			}
		})
	}
}
