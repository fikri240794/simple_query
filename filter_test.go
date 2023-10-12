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
		Logic       FilterLogic
		Expectation *Filter
	} = []struct {
		Name        string
		Logic       FilterLogic
		Expectation *Filter
	}{
		{
			Name:  "logic and",
			Logic: FilterLogicAnd,
			Expectation: &Filter{
				Logic: FilterLogicAnd,
			},
		},
		{
			Name:  "logic or",
			Logic: FilterLogicOr,
			Expectation: &Filter{
				Logic: FilterLogicOr,
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
		Operator    FilterOperator
		Value       interface{}
		Expectation *Filter
	} = []struct {
		Name        string
		Field       string
		Operator    FilterOperator
		Value       interface{}
		Expectation *Filter
	}{
		{
			Name:     fmt.Sprintf("operator %s", FilterOperatorEqual),
			Field:    "field1",
			Operator: FilterOperatorEqual,
			Value:    "value1",
			Expectation: &Filter{
				Field:    "field1",
				Operator: FilterOperatorEqual,
				Value:    "value1",
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", FilterOperatorNotEqual),
			Field:    "field1",
			Operator: FilterOperatorNotEqual,
			Value:    true,
			Expectation: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotEqual,
				Value:    true,
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", FilterOperatorGreaterThan),
			Field:    "field1",
			Operator: FilterOperatorGreaterThan,
			Value:    int64(100),
			Expectation: &Filter{
				Field:    "field1",
				Operator: FilterOperatorGreaterThan,
				Value:    int64(100),
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", FilterOperatorGreaterThanOrEqual),
			Field:    "field1",
			Operator: FilterOperatorGreaterThanOrEqual,
			Value:    float64(100),
			Expectation: &Filter{
				Field:    "field1",
				Operator: FilterOperatorGreaterThanOrEqual,
				Value:    float64(100),
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", FilterOperatorLessThan),
			Field:    "field1",
			Operator: FilterOperatorLessThan,
			Value:    uint64(100),
			Expectation: &Filter{
				Field:    "field1",
				Operator: FilterOperatorLessThan,
				Value:    uint64(100),
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", FilterOperatorLessThanOrEqual),
			Field:    "field1",
			Operator: FilterOperatorLessThanOrEqual,
			Value:    "2006-01-02T15:04:05+07:00",
			Expectation: &Filter{
				Field:    "field1",
				Operator: FilterOperatorLessThanOrEqual,
				Value:    "2006-01-02T15:04:05+07:00",
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", FilterOperatorIsNull),
			Field:    "field1",
			Operator: FilterOperatorIsNull,
			Value:    nil,
			Expectation: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIsNull,
				Value:    nil,
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", FilterOperatorIsNotNull),
			Field:    "field1",
			Operator: FilterOperatorIsNotNull,
			Value:    nil,
			Expectation: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIsNotNull,
				Value:    nil,
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", FilterOperatorIn),
			Field:    "field1",
			Operator: FilterOperatorIn,
			Value:    []string{"value1", "value 2", "value3"},
			Expectation: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIn,
				Value:    []string{"value1", "value 2", "value3"},
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", FilterOperatorNotIn),
			Field:    "field1",
			Operator: FilterOperatorNotIn,
			Value:    [3]int64{1, 2, 3},
			Expectation: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotIn,
				Value:    [3]int64{1, 2, 3},
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", FilterOperatorLike),
			Field:    "field1",
			Operator: FilterOperatorLike,
			Value:    "value1",
			Expectation: &Filter{
				Field:    "field1",
				Operator: FilterOperatorLike,
				Value:    "value1",
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", FilterOperatorNotLike),
			Field:    "field1",
			Operator: FilterOperatorNotLike,
			Value:    "value1",
			Expectation: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotLike,
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
		Operator    FilterOperator
		Value       interface{}
		Expectation *Filter
	} = []struct {
		Name        string
		Field       string
		Operator    FilterOperator
		Value       interface{}
		Expectation *Filter
	}{
		{
			Name:     fmt.Sprintf("operator %s", FilterOperatorEqual),
			Field:    "field1",
			Operator: FilterOperatorEqual,
			Value:    "value1",
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorEqual,
						Value:    "value1",
					},
				},
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", FilterOperatorNotEqual),
			Field:    "field1",
			Operator: FilterOperatorNotEqual,
			Value:    true,
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorNotEqual,
						Value:    true,
					},
				},
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", FilterOperatorGreaterThan),
			Field:    "field1",
			Operator: FilterOperatorGreaterThan,
			Value:    int64(100),
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorGreaterThan,
						Value:    int64(100),
					},
				},
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", FilterOperatorGreaterThanOrEqual),
			Field:    "field1",
			Operator: FilterOperatorGreaterThanOrEqual,
			Value:    float64(100),
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorGreaterThanOrEqual,
						Value:    float64(100),
					},
				},
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", FilterOperatorLessThan),
			Field:    "field1",
			Operator: FilterOperatorLessThan,
			Value:    uint64(100),
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorLessThan,
						Value:    uint64(100),
					},
				},
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", FilterOperatorLessThanOrEqual),
			Field:    "field1",
			Operator: FilterOperatorLessThanOrEqual,
			Value:    "2006-01-02T15:04:05+07:00",
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorLessThanOrEqual,
						Value:    "2006-01-02T15:04:05+07:00",
					},
				},
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", FilterOperatorIsNull),
			Field:    "field1",
			Operator: FilterOperatorIsNull,
			Value:    nil,
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorIsNull,
						Value:    nil,
					},
				},
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", FilterOperatorIsNotNull),
			Field:    "field1",
			Operator: FilterOperatorIsNotNull,
			Value:    nil,
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorIsNotNull,
						Value:    nil,
					},
				},
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", FilterOperatorIn),
			Field:    "field1",
			Operator: FilterOperatorIn,
			Value:    []string{"value1", "value 2", "value3"},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorIn,
						Value:    []string{"value1", "value 2", "value3"},
					},
				},
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", FilterOperatorNotIn),
			Field:    "field1",
			Operator: FilterOperatorNotIn,
			Value:    [3]int64{1, 2, 3},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorNotIn,
						Value:    [3]int64{1, 2, 3},
					},
				},
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", FilterOperatorLike),
			Field:    "field1",
			Operator: FilterOperatorLike,
			Value:    "value1",
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorLike,
						Value:    "value1",
					},
				},
			},
		},
		{
			Name:     fmt.Sprintf("operator %s", FilterOperatorNotLike),
			Field:    "field1",
			Operator: FilterOperatorNotLike,
			Value:    "value1",
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorNotLike,
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
			Name: fmt.Sprintf("operator %s", FilterOperatorEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorEqual,
				Value:    "value1",
			},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorEqual,
						Value:    "value1",
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("operator %s", FilterOperatorNotEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotEqual,
				Value:    true,
			},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorNotEqual,
						Value:    true,
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("operator %s", FilterOperatorGreaterThan),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorGreaterThan,
				Value:    int64(100),
			},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorGreaterThan,
						Value:    int64(100),
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("operator %s", FilterOperatorGreaterThanOrEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorGreaterThanOrEqual,
				Value:    float64(100),
			},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorGreaterThanOrEqual,
						Value:    float64(100),
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("operator %s", FilterOperatorLessThan),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorLessThan,
				Value:    uint64(100),
			},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorLessThan,
						Value:    uint64(100),
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("operator %s", FilterOperatorLessThanOrEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorLessThanOrEqual,
				Value:    "2006-01-02T15:04:05+07:00",
			},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorLessThanOrEqual,
						Value:    "2006-01-02T15:04:05+07:00",
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("operator %s", FilterOperatorIsNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIsNull,
				Value:    nil,
			},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorIsNull,
						Value:    nil,
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("operator %s", FilterOperatorIsNotNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIsNotNull,
				Value:    nil,
			},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorIsNotNull,
						Value:    nil,
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("operator %s", FilterOperatorIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIn,
				Value:    []string{"value1", "value 2", "value3"},
			},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorIn,
						Value:    []string{"value1", "value 2", "value3"},
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("operator %s", FilterOperatorNotIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotIn,
				Value:    [3]int64{1, 2, 3},
			},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorNotIn,
						Value:    [3]int64{1, 2, 3},
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("operator %s", FilterOperatorLike),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorLike,
				Value:    "value1",
			},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorLike,
						Value:    "value1",
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("operator %s", FilterOperatorNotLike),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotLike,
				Value:    "value1",
			},
			Expectation: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorNotLike,
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
				Logic: FilterLogicAnd,
				Field: "field1",
			},
			Expectation: errors.New("field is not empty"),
		},
		{
			Name: "logic is not empty and operator is not empty",
			Filter: &Filter{
				Logic:    FilterLogicOr,
				Operator: FilterOperatorEqual,
			},
			Expectation: errors.New("operator is not empty"),
		},
		{
			Name: "logic is not empty and value is not nil or value kind is allowed",
			Filter: &Filter{
				Logic: FilterLogicAnd,
				Value: "value1",
			},
			Expectation: errors.New("value is not empty"),
		},
		{
			Name: "logic is not empty and filters length is zero",
			Filter: &Filter{
				Logic:   FilterLogicAnd,
				Filters: []*Filter{},
			},
			Expectation: errors.New("filters is required"),
		},
		{
			Name: "logic is empty and filters length greater than zero",
			Filter: &Filter{
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorEqual,
						Value:    "value1",
					},
				},
			},
			Expectation: errors.New("logic is required"),
		},
		{
			Name: "logic is empty and filters length is zero and field is empty",
			Filter: &Filter{
				Operator: FilterOperatorEqual,
				Value:    "valu1",
			},
			Expectation: errors.New("field is required"),
		},
		{
			Name: "logic is empty and filters length is zero and operator is empty",
			Filter: &Filter{
				Field: "field1",
				Value: "value1",
			},
			Expectation: errors.New("operator is required"),
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is not %s and operator is not %s and value is nil", FilterOperatorIsNull, FilterOperatorIsNotNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorEqual,
				Value:    nil,
			},
			Expectation: errors.New("value is required"),
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is not %s and operator is not %s and value kind is not allowed", FilterOperatorIsNull, FilterOperatorIsNotNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorEqual,
				Value:    map[string]string{"key1": "value1"},
			},
			Expectation: fmt.Errorf("unsupported %s value type for operator %s", reflect.Map.String(), FilterOperatorEqual),
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is %s and value is not nil", FilterOperatorIsNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIsNull,
				Value:    "value1",
			},
			Expectation: errors.New("value is not empty"),
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is %s and value is not nil", FilterOperatorIsNotNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIsNotNull,
				Value:    "value1",
			},
			Expectation: errors.New("value is not empty"),
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is not %s and operator is not %s and value kind is %s", FilterOperatorIn, FilterOperatorNotIn, reflect.Slice.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorEqual,
				Value:    []int64{1, 2, 3},
			},
			Expectation: fmt.Errorf("unsupported %s value type for operator %s", reflect.Slice.String(), FilterOperatorEqual),
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is not %s and operator is not %s and value kind is %s", FilterOperatorIn, FilterOperatorNotIn, reflect.Array.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorEqual,
				Value:    [3]string{"value1", "value2", "value3"},
			},
			Expectation: fmt.Errorf("unsupported %s value type for operator %s", reflect.Array.String(), FilterOperatorEqual),
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is %s and value kind is not %s and %s", FilterOperatorIn, reflect.Slice.String(), reflect.Array.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIn,
				Value:    int64(123),
			},
			Expectation: fmt.Errorf("unsupported %s value type for operator %s", reflect.Int64.String(), FilterOperatorIn),
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is %s and value length is zero", FilterOperatorIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIn,
				Value:    []int64{},
			},
			Expectation: errors.New("value is required"),
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is %s and kind of element value is not allowed", FilterOperatorIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIn,
				Value:    []map[string]string{{"key1": "value1"}},
			},
			Expectation: fmt.Errorf("unsupported %s type of element value for operator %s", reflect.Map.String(), FilterOperatorIn),
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is %s and kind of element value is %s", FilterOperatorIn, reflect.Slice.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIn,
				Value:    [][]string{{"value1", "value2", "value3"}},
			},
			Expectation: fmt.Errorf("unsupported %s type of element value for operator %s", reflect.Slice.String(), FilterOperatorIn),
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is %s and kind of element value is %s", FilterOperatorIn, reflect.Array.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIn,
				Value:    [][3]string{{"value1", "value2", "value3"}},
			},
			Expectation: fmt.Errorf("unsupported %s type of element value for operator %s", reflect.Array.String(), FilterOperatorIn),
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is %s and value kind is not allowed", FilterOperatorNotIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotIn,
				Value:    int64(123),
			},
			Expectation: fmt.Errorf("unsupported %s value type for operator %s", reflect.Int64.String(), FilterOperatorNotIn),
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is %s and value length is zero", FilterOperatorNotIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotIn,
				Value:    []int64{},
			},
			Expectation: errors.New("value is required"),
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is %s and kind of element value is not allowed", FilterOperatorNotIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotIn,
				Value:    []map[string]string{{"key1": "value1"}},
			},
			Expectation: fmt.Errorf("unsupported %s type of element value for operator %s", reflect.Map.String(), FilterOperatorNotIn),
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is %s and kind of element value is %s", FilterOperatorNotIn, reflect.Slice.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotIn,
				Value:    [][]string{{"value1", "value2", "value3"}},
			},
			Expectation: fmt.Errorf("unsupported %s type of element value for operator %s", reflect.Slice.String(), FilterOperatorNotIn),
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is %s and kind of element value is %s", FilterOperatorNotIn, reflect.Array.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotIn,
				Value:    [][3]string{{"value1", "value2", "value3"}},
			},
			Expectation: fmt.Errorf("unsupported %s type of element value for operator %s", reflect.Array.String(), FilterOperatorNotIn),
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is %s and value kind is not %s", FilterOperatorLike, reflect.String.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorLike,
				Value:    int64(123),
			},
			Expectation: fmt.Errorf("unsupported %s type of value for operator %s", reflect.Int64.String(), FilterOperatorLike),
		},
		{
			Name: fmt.Sprintf("logic is empty and filters length is zero and operator is %s and value kind is not %s", FilterOperatorNotLike, reflect.String.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotLike,
				Value:    int64(123),
			},
			Expectation: fmt.Errorf("unsupported %s type of value for operator %s", reflect.Int64.String(), FilterOperatorNotLike),
		},
		{
			Name: "filter is valid",
			Filter: &Filter{
				Logic: FilterLogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorEqual,
						Value:    int64(123),
					},
					{
						Field:    "field2",
						Operator: FilterOperatorEqual,
						Value:    "value1",
					},
				},
			},
			Expectation: nil,
		},
		{
			Name: "filter is invalid",
			Filter: &Filter{
				Logic: FilterLogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorEqual,
						Value:    int64(123),
					},
					{
						Field:    "field2",
						Operator: FilterOperatorEqual,
						Value:    []string{"a", "b", "c"},
					},
				},
			},
			Expectation: fmt.Errorf("unsupported %s value type for operator %s", reflect.Slice.String(), FilterOperatorEqual),
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

func TestFilter_typedSliceToInterfaceSlice(t *testing.T) {
	var testCases []struct {
		Name        string
		Value       interface{}
		Expectation struct {
			Values []interface{}
			Error  error
		}
	} = []struct {
		Name        string
		Value       interface{}
		Expectation struct {
			Values []interface{}
			Error  error
		}
	}{
		{
			Name:  "value kind is not slice and value kind is not array",
			Value: "value1",
			Expectation: struct {
				Values []interface{}
				Error  error
			}{
				Values: nil,
				Error:  fmt.Errorf("unsupported %s value type", reflect.String.String()),
			},
		},
		{
			Name:  "kind of element value is not allowed",
			Value: []map[string]string{{"key1": "value1"}},
			Expectation: struct {
				Values []interface{}
				Error  error
			}{
				Values: nil,
				Error:  fmt.Errorf("unsupported %s type of element value", reflect.Map.String()),
			},
		},
		{
			Name:  fmt.Sprintf("kind of element value is %s", reflect.Slice.String()),
			Value: [][]string{{"value1", "value2", "value3"}},
			Expectation: struct {
				Values []interface{}
				Error  error
			}{
				Values: nil,
				Error:  fmt.Errorf("unsupported %s type of element value", reflect.Slice.String()),
			},
		},
		{
			Name:  fmt.Sprintf("kind of element value is %s", reflect.Array.String()),
			Value: [][3]string{{"value1", "value2", "value3"}},
			Expectation: struct {
				Values []interface{}
				Error  error
			}{
				Values: nil,
				Error:  fmt.Errorf("unsupported %s type of element value", reflect.Array.String()),
			},
		},
		{
			Name:  "slice of string to slice of interface",
			Value: []string{"value1", "value2", "value3"},
			Expectation: struct {
				Values []interface{}
				Error  error
			}{
				Values: []interface{}{"value1", "value2", "value3"},
				Error:  nil,
			},
		},
	}

	for i := 0; i < len(testCases); i++ {
		t.Run(testCases[i].Name, func(t *testing.T) {
			var (
				actualValues []interface{}
				actualErr    error
			)

			actualValues, actualErr = NewFilter().typedSliceToInterfaceSlice(testCases[i].Value)

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
				if len(testCases[i].Expectation.Values) != len(actualValues) {
					t.Errorf("expectation values length is %d, got %d", len(testCases[i].Expectation.Values), len(actualValues))
				}

				for x := 0; x < len(testCases[i].Expectation.Values); x++ {
					if !deepEqual(testCases[i].Expectation.Values[x], actualValues[x]) {
						t.Errorf("expectation element slice of interface is %v, got %v", testCases[i].Expectation.Values[x], actualValues[x])
					}
				}
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
				Error:          errors.New("dialect is required"),
			},
		},

		// MYSQL
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectMySQL, FilterOperatorEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorEqual,
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
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[FilterOperatorEqual], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectMySQL, FilterOperatorEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorEqual,
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
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[FilterOperatorEqual], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectMySQL, FilterOperatorNotEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotEqual,
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
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[FilterOperatorNotEqual], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectMySQL, FilterOperatorNotEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotEqual,
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
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[FilterOperatorNotEqual], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectMySQL, FilterOperatorGreaterThan),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorGreaterThan,
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
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[FilterOperatorGreaterThan], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectMySQL, FilterOperatorGreaterThan),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorGreaterThan,
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
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[FilterOperatorGreaterThan], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectMySQL, FilterOperatorGreaterThanOrEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorGreaterThanOrEqual,
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
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[FilterOperatorGreaterThanOrEqual], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectMySQL, FilterOperatorGreaterThanOrEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorGreaterThanOrEqual,
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
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[FilterOperatorGreaterThanOrEqual], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectMySQL, FilterOperatorLessThan),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorLessThan,
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
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[FilterOperatorLessThan], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectMySQL, FilterOperatorLessThan),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorLessThan,
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
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[FilterOperatorLessThan], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectMySQL, FilterOperatorLessThanOrEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorLessThanOrEqual,
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
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[FilterOperatorLessThanOrEqual], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectMySQL, FilterOperatorLessThanOrEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorLessThanOrEqual,
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
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[FilterOperatorLessThanOrEqual], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectMySQL, FilterOperatorIsNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIsNull,
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
				ConditionQuery: fmt.Sprintf("%s %s", "field1", filterOperatorMap[FilterOperatorIsNull]),
				Args:           []interface{}{},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectMySQL, FilterOperatorIsNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIsNull,
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
				ConditionQuery: fmt.Sprintf("%s %s", "field1", filterOperatorMap[FilterOperatorIsNull]),
				Args:           []interface{}{"other args value"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectMySQL, FilterOperatorIsNotNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIsNotNull,
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
				ConditionQuery: fmt.Sprintf("%s %s", "field1", filterOperatorMap[FilterOperatorIsNotNull]),
				Args:           []interface{}{},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectMySQL, FilterOperatorIsNotNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIsNotNull,
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
				ConditionQuery: fmt.Sprintf("%s %s", "field1", filterOperatorMap[FilterOperatorIsNotNull]),
				Args:           []interface{}{"other args value"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectMySQL, FilterOperatorIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIn,
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
				ConditionQuery: fmt.Sprintf("%s %s (%s)", "field1", filterOperatorMap[FilterOperatorIn], fmt.Sprintf("%s, %s, %s", placeholderMap[DialectMySQL], placeholderMap[DialectMySQL], placeholderMap[DialectMySQL])),
				Args:           []interface{}{"value1", "value2", "value3"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectMySQL, FilterOperatorIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIn,
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
				ConditionQuery: fmt.Sprintf("%s %s (%s)", "field1", filterOperatorMap[FilterOperatorIn], fmt.Sprintf("%s, %s, %s", placeholderMap[DialectMySQL], placeholderMap[DialectMySQL], placeholderMap[DialectMySQL])),
				Args:           []interface{}{"other args value", "value1", "value2", "value3"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s with value kind is not %s and %s", DialectMySQL, FilterOperatorIn, reflect.Slice.String(), reflect.Array.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIn,
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
				Error:          fmt.Errorf("unsupported %s value type for operator %s", reflect.String.String(), FilterOperatorIn),
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s with kind of element value is not allowed", DialectMySQL, FilterOperatorIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIn,
				Value:    []map[string]string{{"key1": "value1"}},
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
				Error:          fmt.Errorf("unsupported %s type of element value for operator %s", reflect.Map.String(), FilterOperatorIn),
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s with kind of element value is %s", DialectMySQL, FilterOperatorIn, reflect.Slice.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIn,
				Value:    [][]string{{"value1", "value2", "value3"}},
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
				Error:          fmt.Errorf("unsupported %s type of element value for operator %s", reflect.Slice.String(), FilterOperatorIn),
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s with kind of element value is %s", DialectMySQL, FilterOperatorIn, reflect.Array.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIn,
				Value:    [][3]string{{"value1", "value2", "value3"}},
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
				Error:          fmt.Errorf("unsupported %s type of element value for operator %s", reflect.Array.String(), FilterOperatorIn),
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectMySQL, FilterOperatorNotIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotIn,
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
				ConditionQuery: fmt.Sprintf("%s %s (%s)", "field1", filterOperatorMap[FilterOperatorNotIn], fmt.Sprintf("%s, %s, %s", placeholderMap[DialectMySQL], placeholderMap[DialectMySQL], placeholderMap[DialectMySQL])),
				Args:           []interface{}{"value1", "value2", "value3"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectMySQL, FilterOperatorNotIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotIn,
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
				ConditionQuery: fmt.Sprintf("%s %s (%s)", "field1", filterOperatorMap[FilterOperatorNotIn], fmt.Sprintf("%s, %s, %s", placeholderMap[DialectMySQL], placeholderMap[DialectMySQL], placeholderMap[DialectMySQL])),
				Args:           []interface{}{"other args value", "value1", "value2", "value3"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s with value kind is not %s and %s", DialectMySQL, FilterOperatorNotIn, reflect.Slice.String(), reflect.Array.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotIn,
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
				Error:          fmt.Errorf("unsupported %s value type for operator %s", reflect.String.String(), FilterOperatorNotIn),
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s with kind of element value is not allowed", DialectMySQL, FilterOperatorNotIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotIn,
				Value:    []map[string]string{{"key1": "value1"}},
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
				Error:          fmt.Errorf("unsupported %s type of element value for operator %s", reflect.Map.String(), FilterOperatorNotIn),
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s with kind of element value is %s", DialectMySQL, FilterOperatorNotIn, reflect.Slice.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotIn,
				Value:    [][]string{{"value1", "value2", "value3"}},
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
				Error:          fmt.Errorf("unsupported %s type of element value for operator %s", reflect.Slice.String(), FilterOperatorNotIn),
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s with kind of element value is %s", DialectMySQL, FilterOperatorNotIn, reflect.Array.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotIn,
				Value:    [][3]string{{"value1", "value2", "value3"}},
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
				Error:          fmt.Errorf("unsupported %s type of element value for operator %s", reflect.Array.String(), FilterOperatorNotIn),
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectMySQL, FilterOperatorLike),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorLike,
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
				ConditionQuery: fmt.Sprintf("%s %s concat('%%', %s, '%%')", "field1", filterOperatorMap[FilterOperatorLike], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectMySQL, FilterOperatorLike),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorLike,
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
				ConditionQuery: fmt.Sprintf("%s %s concat('%%', %s, '%%')", "field1", filterOperatorMap[FilterOperatorLike], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectMySQL, FilterOperatorNotLike),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotLike,
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
				ConditionQuery: fmt.Sprintf("%s %s concat('%%', %s, '%%')", "field1", filterOperatorMap[FilterOperatorNotLike], placeholderMap[DialectMySQL]),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectMySQL, FilterOperatorNotLike),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotLike,
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
				ConditionQuery: fmt.Sprintf("%s %s concat('%%', %s, '%%')", "field1", filterOperatorMap[FilterOperatorNotLike], placeholderMap[DialectMySQL]),
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
				Logic: FilterLogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorEqual,
						Value:    "value1",
					},
					{
						Logic: FilterLogicOr,
						Filters: []*Filter{
							{
								Field:    "field2",
								Operator: FilterOperatorIsNull,
								Value:    nil,
							},
							{
								Field:    "field3",
								Operator: FilterOperatorIn,
								Value:    []int64{1, 2, 3},
							},
						},
					},
					{
						Field:    "field4",
						Operator: FilterOperatorLike,
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
					filterOperatorMap[FilterOperatorEqual],
					placeholderMap[DialectMySQL],
					FilterLogicAnd,
					filterOperatorMap[FilterOperatorIsNull],
					FilterLogicOr,
					filterOperatorMap[FilterOperatorIn],
					fmt.Sprintf("(%s, %s, %s)", placeholderMap[DialectMySQL], placeholderMap[DialectMySQL], placeholderMap[DialectMySQL]),
					FilterLogicAnd,
					filterOperatorMap[FilterOperatorLike],
					fmt.Sprintf("concat('%%', %s, '%%')", placeholderMap[DialectMySQL]),
				),
				Args:  []interface{}{"value1", 1, 2, 3, "value4"},
				Error: nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with element filters is empty filter", DialectMySQL),
			Filter: &Filter{
				Logic: FilterLogicAnd,
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
				Logic: FilterLogicAnd,
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
				Logic: FilterLogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorIn,
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
				Error:          fmt.Errorf("unsupported %s value type for operator %s", reflect.String.String(), FilterOperatorIn),
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with element values of element filters is invalid", DialectMySQL),
			Filter: &Filter{
				Logic: FilterLogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorNotIn,
						Value:    [][]string{{"a", "b", "c"}},
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
				Error:          fmt.Errorf("unsupported %s type of element value for operator %s", reflect.Slice.String(), FilterOperatorNotIn),
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero", DialectMySQL),
			Filter: &Filter{
				Logic: FilterLogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorEqual,
						Value:    "value1",
					},
					{
						Logic: FilterLogicOr,
						Filters: []*Filter{
							{
								Field:    "field2",
								Operator: FilterOperatorIsNull,
								Value:    nil,
							},
							{
								Field:    "field3",
								Operator: FilterOperatorIn,
								Value:    []int64{1, 2, 3},
							},
						},
					},
					{
						Field:    "field4",
						Operator: FilterOperatorLike,
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
					filterOperatorMap[FilterOperatorEqual],
					placeholderMap[DialectMySQL],
					FilterLogicAnd,
					filterOperatorMap[FilterOperatorIsNull],
					FilterLogicOr,
					filterOperatorMap[FilterOperatorIn],
					fmt.Sprintf("(%s, %s, %s)", placeholderMap[DialectMySQL], placeholderMap[DialectMySQL], placeholderMap[DialectMySQL]),
					FilterLogicAnd,
					filterOperatorMap[FilterOperatorLike],
					fmt.Sprintf("concat('%%', %s, '%%')", placeholderMap[DialectMySQL]),
				),
				Args:  []interface{}{"other args value", "value1", 1, 2, 3, "value4"},
				Error: nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with isRoot is true", DialectMySQL),
			Filter: &Filter{
				Logic: FilterLogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorEqual,
						Value:    "value1",
					},
					{
						Logic: FilterLogicOr,
						Filters: []*Filter{
							{
								Field:    "field2",
								Operator: FilterOperatorIsNull,
								Value:    nil,
							},
							{
								Field:    "field3",
								Operator: FilterOperatorIn,
								Value:    []int64{1, 2, 3},
							},
						},
					},
					{
						Field:    "field4",
						Operator: FilterOperatorLike,
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
					filterOperatorMap[FilterOperatorEqual],
					placeholderMap[DialectMySQL],
					FilterLogicAnd,
					filterOperatorMap[FilterOperatorIsNull],
					FilterLogicOr,
					filterOperatorMap[FilterOperatorIn],
					fmt.Sprintf("(%s, %s, %s)", placeholderMap[DialectMySQL], placeholderMap[DialectMySQL], placeholderMap[DialectMySQL]),
					FilterLogicAnd,
					filterOperatorMap[FilterOperatorLike],
					fmt.Sprintf("concat('%%', %s, '%%')", placeholderMap[DialectMySQL]),
				),
				Args:  []interface{}{"value1", 1, 2, 3, "value4"},
				Error: nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero and isRoot is true", DialectMySQL),
			Filter: &Filter{
				Logic: FilterLogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorEqual,
						Value:    "value1",
					},
					{
						Logic: FilterLogicOr,
						Filters: []*Filter{
							{
								Field:    "field2",
								Operator: FilterOperatorIsNull,
								Value:    nil,
							},
							{
								Field:    "field3",
								Operator: FilterOperatorIn,
								Value:    []int64{1, 2, 3},
							},
						},
					},
					{
						Field:    "field4",
						Operator: FilterOperatorLike,
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
					filterOperatorMap[FilterOperatorEqual],
					placeholderMap[DialectMySQL],
					FilterLogicAnd,
					filterOperatorMap[FilterOperatorIsNull],
					FilterLogicOr,
					filterOperatorMap[FilterOperatorIn],
					fmt.Sprintf("(%s, %s, %s)", placeholderMap[DialectMySQL], placeholderMap[DialectMySQL], placeholderMap[DialectMySQL]),
					FilterLogicAnd,
					filterOperatorMap[FilterOperatorLike],
					fmt.Sprintf("concat('%%', %s, '%%')", placeholderMap[DialectMySQL]),
				),
				Args:  []interface{}{"other args value", "value1", 1, 2, 3, "value4"},
				Error: nil,
			},
		},

		// POSTGRES
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectPostgres, FilterOperatorEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorEqual,
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
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[FilterOperatorEqual], fmt.Sprintf("%s1", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectPostgres, FilterOperatorEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorEqual,
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
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[FilterOperatorEqual], fmt.Sprintf("%s2", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectPostgres, FilterOperatorNotEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotEqual,
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
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[FilterOperatorNotEqual], fmt.Sprintf("%s1", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectPostgres, FilterOperatorNotEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotEqual,
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
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[FilterOperatorNotEqual], fmt.Sprintf("%s2", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectPostgres, FilterOperatorGreaterThan),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorGreaterThan,
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
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[FilterOperatorGreaterThan], fmt.Sprintf("%s1", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectPostgres, FilterOperatorGreaterThan),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorGreaterThan,
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
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[FilterOperatorGreaterThan], fmt.Sprintf("%s2", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectPostgres, FilterOperatorGreaterThanOrEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorGreaterThanOrEqual,
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
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[FilterOperatorGreaterThanOrEqual], fmt.Sprintf("%s1", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectPostgres, FilterOperatorGreaterThanOrEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorGreaterThanOrEqual,
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
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[FilterOperatorGreaterThanOrEqual], fmt.Sprintf("%s2", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectPostgres, FilterOperatorLessThan),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorLessThan,
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
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[FilterOperatorLessThan], fmt.Sprintf("%s1", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectPostgres, FilterOperatorLessThan),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorLessThan,
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
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[FilterOperatorLessThan], fmt.Sprintf("%s2", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectPostgres, FilterOperatorLessThanOrEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorLessThanOrEqual,
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
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[FilterOperatorLessThanOrEqual], fmt.Sprintf("%s1", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectPostgres, FilterOperatorLessThanOrEqual),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorLessThanOrEqual,
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
				ConditionQuery: fmt.Sprintf("%s %s %s", "field1", filterOperatorMap[FilterOperatorLessThanOrEqual], fmt.Sprintf("%s2", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectPostgres, FilterOperatorIsNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIsNull,
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
				ConditionQuery: fmt.Sprintf("%s %s", "field1", filterOperatorMap[FilterOperatorIsNull]),
				Args:           []interface{}{},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectPostgres, FilterOperatorIsNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIsNull,
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
				ConditionQuery: fmt.Sprintf("%s %s", "field1", filterOperatorMap[FilterOperatorIsNull]),
				Args:           []interface{}{"other args value"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectPostgres, FilterOperatorIsNotNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIsNotNull,
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
				ConditionQuery: fmt.Sprintf("%s %s", "field1", filterOperatorMap[FilterOperatorIsNotNull]),
				Args:           []interface{}{},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectPostgres, FilterOperatorIsNotNull),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIsNotNull,
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
				ConditionQuery: fmt.Sprintf("%s %s", "field1", filterOperatorMap[FilterOperatorIsNotNull]),
				Args:           []interface{}{"other args value"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectPostgres, FilterOperatorIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIn,
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
				ConditionQuery: fmt.Sprintf("%s %s (%s)", "field1", filterOperatorMap[FilterOperatorIn], fmt.Sprintf("%s1, %s2, %s3", placeholderMap[DialectPostgres], placeholderMap[DialectPostgres], placeholderMap[DialectPostgres])),
				Args:           []interface{}{"value1", "value2", "value3"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectPostgres, FilterOperatorIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIn,
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
				ConditionQuery: fmt.Sprintf("%s %s (%s)", "field1", filterOperatorMap[FilterOperatorIn], fmt.Sprintf("%s2, %s3, %s4", placeholderMap[DialectPostgres], placeholderMap[DialectPostgres], placeholderMap[DialectPostgres])),
				Args:           []interface{}{"other args value", "value1", "value2", "value3"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s with value kind is not %s and %s", DialectPostgres, FilterOperatorIn, reflect.Slice.String(), reflect.Array.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIn,
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
				Error:          fmt.Errorf("unsupported %s value type for operator %s", reflect.String.String(), FilterOperatorIn),
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s with kind of element value is not allowed", DialectPostgres, FilterOperatorIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIn,
				Value:    []map[string]string{{"key1": "value1"}},
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
				Error:          fmt.Errorf("unsupported %s type of element value for operator %s", reflect.Map.String(), FilterOperatorIn),
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s with kind of element value is %s", DialectPostgres, FilterOperatorIn, reflect.Slice.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIn,
				Value:    [][]string{{"value1", "value2", "value3"}},
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
				Error:          fmt.Errorf("unsupported %s type of element value for operator %s", reflect.Slice.String(), FilterOperatorIn),
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s with kind of element value is %s", DialectPostgres, FilterOperatorIn, reflect.Array.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorIn,
				Value:    [][3]string{{"value1", "value2", "value3"}},
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
				Error:          fmt.Errorf("unsupported %s type of element value for operator %s", reflect.Array.String(), FilterOperatorIn),
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectPostgres, FilterOperatorNotIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotIn,
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
				ConditionQuery: fmt.Sprintf("%s %s (%s)", "field1", filterOperatorMap[FilterOperatorNotIn], fmt.Sprintf("%s1, %s2, %s3", placeholderMap[DialectPostgres], placeholderMap[DialectPostgres], placeholderMap[DialectPostgres])),
				Args:           []interface{}{"value1", "value2", "value3"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectPostgres, FilterOperatorNotIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotIn,
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
				ConditionQuery: fmt.Sprintf("%s %s (%s)", "field1", filterOperatorMap[FilterOperatorNotIn], fmt.Sprintf("%s2, %s3, %s4", placeholderMap[DialectPostgres], placeholderMap[DialectPostgres], placeholderMap[DialectPostgres])),
				Args:           []interface{}{"other args value", "value1", "value2", "value3"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s with value kind is not %s and %s", DialectPostgres, FilterOperatorNotIn, reflect.Slice.String(), reflect.Array.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotIn,
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
				Error:          fmt.Errorf("unsupported %s value type for operator %s", reflect.String.String(), FilterOperatorNotIn),
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s with kind of element value is not allowed", DialectPostgres, FilterOperatorNotIn),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotIn,
				Value:    []map[string]string{{"key1": "value1"}},
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
				Error:          fmt.Errorf("unsupported %s type of element value for operator %s", reflect.Map.String(), FilterOperatorNotIn),
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s with kind of element value is %s", DialectPostgres, FilterOperatorNotIn, reflect.Slice.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotIn,
				Value:    [][]string{{"value1", "value2", "value3"}},
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
				Error:          fmt.Errorf("unsupported %s type of element value for operator %s", reflect.Slice.String(), FilterOperatorNotIn),
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s with kind of element value is %s", DialectPostgres, FilterOperatorNotIn, reflect.Array.String()),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotIn,
				Value:    [][3]string{{"value1", "value2", "value3"}},
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
				Error:          fmt.Errorf("unsupported %s type of element value for operator %s", reflect.Array.String(), FilterOperatorNotIn),
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectPostgres, FilterOperatorLike),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorLike,
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
				ConditionQuery: fmt.Sprintf("%s %s concat('%%', %s, '%%')", "field1", fmt.Sprintf("i%s", filterOperatorMap[FilterOperatorLike]), fmt.Sprintf("%s1", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectPostgres, FilterOperatorLike),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorLike,
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
				ConditionQuery: fmt.Sprintf("%s %s concat('%%', %s, '%%')", "field1", fmt.Sprintf("i%s", filterOperatorMap[FilterOperatorLike]), fmt.Sprintf("%s2", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"other args value", "value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter operator %s", DialectPostgres, FilterOperatorNotLike),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotLike,
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
				ConditionQuery: fmt.Sprintf("%s %s concat('%%', %s, '%%')", "field1", fmt.Sprintf("not i%s", filterOperatorMap[FilterOperatorLike]), fmt.Sprintf("%s1", placeholderMap[DialectPostgres])),
				Args:           []interface{}{"value1"},
				Error:          nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero with filter operator %s", DialectPostgres, FilterOperatorNotLike),
			Filter: &Filter{
				Field:    "field1",
				Operator: FilterOperatorNotLike,
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
				ConditionQuery: fmt.Sprintf("%s %s concat('%%', %s, '%%')", "field1", fmt.Sprintf("not i%s", filterOperatorMap[FilterOperatorLike]), fmt.Sprintf("%s2", placeholderMap[DialectPostgres])),
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
				Logic: FilterLogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorEqual,
						Value:    "value1",
					},
					{
						Logic: FilterLogicOr,
						Filters: []*Filter{
							{
								Field:    "field2",
								Operator: FilterOperatorIsNull,
								Value:    nil,
							},
							{
								Field:    "field3",
								Operator: FilterOperatorIn,
								Value:    []int64{1, 2, 3},
							},
						},
					},
					{
						Field:    "field4",
						Operator: FilterOperatorLike,
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
					filterOperatorMap[FilterOperatorEqual],
					fmt.Sprintf("%s1", placeholderMap[DialectPostgres]),
					FilterLogicAnd,
					filterOperatorMap[FilterOperatorIsNull],
					FilterLogicOr,
					filterOperatorMap[FilterOperatorIn],
					fmt.Sprintf("(%s2, %s3, %s4)", placeholderMap[DialectPostgres], placeholderMap[DialectPostgres], placeholderMap[DialectPostgres]),
					FilterLogicAnd,
					fmt.Sprintf("i%s", filterOperatorMap[FilterOperatorLike]),
					fmt.Sprintf("concat('%%', %s5, '%%')", placeholderMap[DialectPostgres]),
				),
				Args:  []interface{}{"value1", 1, 2, 3, "value4"},
				Error: nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with element filters is empty filter", DialectPostgres),
			Filter: &Filter{
				Logic: FilterLogicAnd,
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
				Logic: FilterLogicAnd,
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
				Logic: FilterLogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorIn,
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
				Error:          fmt.Errorf("unsupported %s value type for operator %s", reflect.String.String(), FilterOperatorIn),
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with element values of element filters is invalid", DialectPostgres),
			Filter: &Filter{
				Logic: FilterLogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorNotIn,
						Value:    [][]string{{"a", "b", "c"}},
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
				Error:          fmt.Errorf("unsupported %s type of element value for operator %s", reflect.Slice.String(), FilterOperatorNotIn),
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero", DialectPostgres),
			Filter: &Filter{
				Logic: FilterLogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorEqual,
						Value:    "value1",
					},
					{
						Logic: FilterLogicOr,
						Filters: []*Filter{
							{
								Field:    "field2",
								Operator: FilterOperatorIsNull,
								Value:    nil,
							},
							{
								Field:    "field3",
								Operator: FilterOperatorIn,
								Value:    []int64{1, 2, 3},
							},
						},
					},
					{
						Field:    "field4",
						Operator: FilterOperatorLike,
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
					filterOperatorMap[FilterOperatorEqual],
					fmt.Sprintf("%s2", placeholderMap[DialectPostgres]),
					FilterLogicAnd,
					filterOperatorMap[FilterOperatorIsNull],
					FilterLogicOr,
					filterOperatorMap[FilterOperatorIn],
					fmt.Sprintf("(%s3, %s4, %s5)", placeholderMap[DialectPostgres], placeholderMap[DialectPostgres], placeholderMap[DialectPostgres]),
					FilterLogicAnd,
					fmt.Sprintf("i%s", filterOperatorMap[FilterOperatorLike]),
					fmt.Sprintf("concat('%%', %s6, '%%')", placeholderMap[DialectPostgres]),
				),
				Args:  []interface{}{"other args value", "value1", 1, 2, 3, "value4"},
				Error: nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with isRoot is true", DialectPostgres),
			Filter: &Filter{
				Logic: FilterLogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorEqual,
						Value:    "value1",
					},
					{
						Logic: FilterLogicOr,
						Filters: []*Filter{
							{
								Field:    "field2",
								Operator: FilterOperatorIsNull,
								Value:    nil,
							},
							{
								Field:    "field3",
								Operator: FilterOperatorIn,
								Value:    []int64{1, 2, 3},
							},
						},
					},
					{
						Field:    "field4",
						Operator: FilterOperatorLike,
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
					filterOperatorMap[FilterOperatorEqual],
					fmt.Sprintf("%s1", placeholderMap[DialectPostgres]),
					FilterLogicAnd,
					filterOperatorMap[FilterOperatorIsNull],
					FilterLogicOr,
					filterOperatorMap[FilterOperatorIn],
					fmt.Sprintf("(%s2, %s3, %s4)", placeholderMap[DialectPostgres], placeholderMap[DialectPostgres], placeholderMap[DialectPostgres]),
					FilterLogicAnd,
					fmt.Sprintf("i%s", filterOperatorMap[FilterOperatorLike]),
					fmt.Sprintf("concat('%%', %s5, '%%')", placeholderMap[DialectPostgres]),
				),
				Args:  []interface{}{"value1", 1, 2, 3, "value4"},
				Error: nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s and args length greater than zero and isRoot is true", DialectPostgres),
			Filter: &Filter{
				Logic: FilterLogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorEqual,
						Value:    "value1",
					},
					{
						Logic: FilterLogicOr,
						Filters: []*Filter{
							{
								Field:    "field2",
								Operator: FilterOperatorIsNull,
								Value:    nil,
							},
							{
								Field:    "field3",
								Operator: FilterOperatorIn,
								Value:    []int64{1, 2, 3},
							},
						},
					},
					{
						Field:    "field4",
						Operator: FilterOperatorLike,
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
					filterOperatorMap[FilterOperatorEqual],
					fmt.Sprintf("%s2", placeholderMap[DialectPostgres]),
					FilterLogicAnd,
					filterOperatorMap[FilterOperatorIsNull],
					FilterLogicOr,
					filterOperatorMap[FilterOperatorIn],
					fmt.Sprintf("(%s3, %s4, %s5)", placeholderMap[DialectPostgres], placeholderMap[DialectPostgres], placeholderMap[DialectPostgres]),
					FilterLogicAnd,
					fmt.Sprintf("i%s", filterOperatorMap[FilterOperatorLike]),
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
				Logic: FilterLogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorEqual,
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
				Error:          fmt.Errorf("unsupported %s value type for operator %s", reflect.Slice.String(), FilterOperatorEqual),
			},
		},
		{
			Name: "filter is valid",
			Filter: &Filter{
				Logic: FilterLogicAnd,
				Filters: []*Filter{
					{
						Field:    "field1",
						Operator: FilterOperatorEqual,
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
