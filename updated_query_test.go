package simple_query

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestUpdateQuery_Update(t *testing.T) {
	var (
		expectation *UpdateQuery
		actual      *UpdateQuery
	)

	expectation = &UpdateQuery{
		Table:       "table1",
		FieldsValue: map[string]interface{}{},
	}
	actual = Update("table1")

	if expectation.Table != actual.Table {
		t.Errorf("expected table is %s, got %s", expectation.Table, actual.Table)
	}

	if len(expectation.FieldsValue) != len(actual.FieldsValue) {
		t.Errorf("expected length of fields value is %d, got %d", len(expectation.FieldsValue), len(actual.FieldsValue))
	}
}

func TestUpdateQuery_Set(t *testing.T) {
	var (
		expectation *UpdateQuery
		actual      *UpdateQuery
	)

	expectation = &UpdateQuery{
		Table: "table1",
		FieldsValue: map[string]interface{}{
			"field1": "value1",
			"field2": 2,
		},
	}

	actual = Update("table1").
		Set("field1", "value1").
		Set("field2", 2)

	if expectation.Table != actual.Table {
		t.Errorf("expected table is %s, got %s", expectation.Table, actual.Table)
	}

	if len(expectation.FieldsValue) != len(actual.FieldsValue) {
		t.Errorf("expected length of fields value is %d, got %d", len(expectation.FieldsValue), len(actual.FieldsValue))
	}

	if len(expectation.FieldsValue) > 0 {
		for field, value := range expectation.FieldsValue {
			if !deepEqual(value, actual.FieldsValue[field]) {
				t.Errorf("expected value is %v, got %v", value, actual.FieldsValue[field])
			}
		}
	}
}

func TestUpdateQuery_Where(t *testing.T) {
	var (
		expectation *UpdateQuery
		actual      *UpdateQuery
	)

	expectation = &UpdateQuery{
		Table: "table1",
		FieldsValue: map[string]interface{}{
			"field1": "value1",
			"field2": 2,
		},
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
	}

	actual = Update("table1").
		Set("field1", "value1").
		Set("field2", 2).
		Where(
			NewFilter().
				SetLogic(FilterLogicAnd).
				AddFilter("field1", FilterOperatorEqual, "value1"),
		)

	if expectation.Table != actual.Table {
		t.Errorf("expected table is %s, got %s", expectation.Table, actual.Table)
	}

	if len(expectation.FieldsValue) != len(actual.FieldsValue) {
		t.Errorf("expected length of fields value is %d, got %d", len(expectation.FieldsValue), len(actual.FieldsValue))
	}

	if len(expectation.FieldsValue) > 0 {
		for field, value := range expectation.FieldsValue {
			if !deepEqual(value, actual.FieldsValue[field]) {
				t.Errorf("expected value is %v, got %v", value, actual.FieldsValue[field])
			}
		}
	}

	if !deepEqual(expectation.Filter, actual.Filter) {
		t.Errorf("expectation filter is %v, got %v", expectation.Filter, actual.Filter)
	}
}

func TestUpdateQuery_validate(t *testing.T) {
	var testCases []struct {
		Name        string
		UpdateQuery *UpdateQuery
		Expectation error
	} = []struct {
		Name        string
		UpdateQuery *UpdateQuery
		Expectation error
	}{
		{
			Name:        "table is empty",
			UpdateQuery: &UpdateQuery{},
			Expectation: errors.New("table is required"),
		},
		{
			Name: "fields and value is empty",
			UpdateQuery: &UpdateQuery{
				Table:       "table1",
				FieldsValue: map[string]interface{}{},
			},
			Expectation: errors.New("fields is required"),
		},
		{
			Name: "field is empty",
			UpdateQuery: &UpdateQuery{
				Table: "table1",
				FieldsValue: map[string]interface{}{
					"": "field1",
				},
			},
			Expectation: errors.New("field is required"),
		},
		{
			Name: "value is not allowed",
			UpdateQuery: &UpdateQuery{
				Table: "table1",
				FieldsValue: map[string]interface{}{
					"field1": map[string]string{
						"key1": "value1",
					},
				},
			},
			Expectation: fmt.Errorf("unsupported %s value type", reflect.Map.String()),
		},
		{
			Name: "value is array",
			UpdateQuery: &UpdateQuery{
				Table: "table1",
				FieldsValue: map[string]interface{}{
					"field1": [3]string{"value1", "value2", "value3"},
				},
			},
			Expectation: fmt.Errorf("unsupported %s value type", reflect.Array.String()),
		},
		{
			Name: "value is slice",
			UpdateQuery: &UpdateQuery{
				Table: "table1",
				FieldsValue: map[string]interface{}{
					"field1": []string{"value1", "value2", "value3"},
				},
			},
			Expectation: fmt.Errorf("unsupported %s value type", reflect.Slice.String()),
		},
		{
			Name: "filter is empty",
			UpdateQuery: &UpdateQuery{
				Table: "table1",
				FieldsValue: map[string]interface{}{
					"field1": "value1",
				},
			},
			Expectation: errors.New("filter is required"),
		},
		{
			Name: "update query is valid",
			UpdateQuery: &UpdateQuery{
				Table: "table1",
				FieldsValue: map[string]interface{}{
					"field1": "value1",
				},
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
			},
			Expectation: nil,
		},
	}

	for i := 0; i < len(testCases); i++ {
		t.Run(testCases[i].Name, func(t *testing.T) {
			var actual error = testCases[i].UpdateQuery.validate()

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

func TestUpdateQuery_ToSQLWithArgs(t *testing.T) {
	var testCases []struct {
		Name        string
		UpdateQuery *UpdateQuery
		Dialect     Dialect
		Expectation struct {
			Query string
			Args  []interface{}
			Error error
		}
	} = []struct {
		Name        string
		UpdateQuery *UpdateQuery
		Dialect     Dialect
		Expectation struct {
			Query string
			Args  []interface{}
			Error error
		}
	}{
		{
			Name:        "table is empty",
			UpdateQuery: &UpdateQuery{},
			Dialect:     "",
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "",
				Args:  nil,
				Error: errors.New("table is required"),
			},
		},
		{
			Name: "fields and value is empty",
			UpdateQuery: &UpdateQuery{
				Table:       "table1",
				FieldsValue: map[string]interface{}{},
			},
			Dialect: "",
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "",
				Args:  nil,
				Error: errors.New("fields is required"),
			},
		},
		{
			Name: "field is empty",
			UpdateQuery: &UpdateQuery{
				Table: "table1",
				FieldsValue: map[string]interface{}{
					"": "field1",
				},
			},
			Dialect: "",
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "",
				Args:  nil,
				Error: errors.New("field is required"),
			},
		},
		{
			Name: "value is not allowed",
			UpdateQuery: &UpdateQuery{
				Table: "table1",
				FieldsValue: map[string]interface{}{
					"field1": map[string]string{
						"key1": "value1",
					},
				},
			},
			Dialect: "",
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "",
				Args:  nil,
				Error: fmt.Errorf("unsupported %s value type", reflect.Map.String()),
			},
		},
		{
			Name: "value is array",
			UpdateQuery: &UpdateQuery{
				Table: "table1",
				FieldsValue: map[string]interface{}{
					"field1": [3]string{"value1", "value2", "value3"},
				},
			},
			Dialect: "",
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "",
				Args:  nil,
				Error: fmt.Errorf("unsupported %s value type", reflect.Array.String()),
			},
		},
		{
			Name: "value is slice",
			UpdateQuery: &UpdateQuery{
				Table: "table1",
				FieldsValue: map[string]interface{}{
					"field1": []string{"value1", "value2", "value3"},
				},
			},
			Dialect: "",
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "",
				Args:  nil,
				Error: fmt.Errorf("unsupported %s value type", reflect.Slice.String()),
			},
		},
		{
			Name: "value is slice",
			UpdateQuery: &UpdateQuery{
				Table: "table1",
				FieldsValue: map[string]interface{}{
					"field1": []string{"value1", "value2", "value3"},
				},
			},
			Dialect: "",
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "",
				Args:  nil,
				Error: fmt.Errorf("unsupported %s value type", reflect.Slice.String()),
			},
		},
		{
			Name: "filter is empty",
			UpdateQuery: &UpdateQuery{
				Table: "table1",
				FieldsValue: map[string]interface{}{
					"field1": "value1",
				},
			},
			Dialect: "",
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "",
				Args:  nil,
				Error: errors.New("filter is required"),
			},
		},
		{
			Name: "filter is invalid", // don't test all invalid filter here, because it's handled in filter_test.go
			UpdateQuery: &UpdateQuery{
				Table: "table1",
				FieldsValue: map[string]interface{}{
					"field1": "value1",
				},
				Filter: &Filter{
					Logic:   FilterLogicAnd,
					Filters: []*Filter{},
				},
			},
			Dialect: "",
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "",
				Args:  nil,
				Error: errors.New("filters is required"),
			},
		},
		{
			Name: fmt.Sprintf("update with dialect %s with filter", DialectMySQL),
			UpdateQuery: &UpdateQuery{
				Table: "table1",
				FieldsValue: map[string]interface{}{
					"field1": "value1",
				},
				Filter: &Filter{
					Logic: FilterLogicAnd,
					Filters: []*Filter{
						{
							Field:    "field2",
							Operator: FilterOperatorEqual,
							Value:    "value2",
						},
					},
				},
			},
			Dialect: DialectMySQL,
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "update table1 set field1 = ? where field2 = ?",
				Args:  []interface{}{"value1", "value2"},
				Error: nil,
			},
		},
		{
			Name: fmt.Sprintf("update with dialect %s with filter", DialectPostgres),
			UpdateQuery: &UpdateQuery{
				Table: "table1",
				FieldsValue: map[string]interface{}{
					"field1": "value1",
				},
				Filter: &Filter{
					Logic: FilterLogicAnd,
					Filters: []*Filter{
						{
							Field:    "field2",
							Operator: FilterOperatorEqual,
							Value:    "value2",
						},
					},
				},
			},
			Dialect: DialectPostgres,
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "update table1 set field1 = $1 where field2 = $2",
				Args:  []interface{}{"value1", "value2"},
				Error: nil,
			},
		},
	}

	for i := 0; i < len(testCases); i++ {
		t.Run(testCases[i].Name, func(t *testing.T) {
			var (
				actualQuery string
				actualArgs  []interface{}
				actualErr   error
			)

			actualQuery, actualArgs, actualErr = testCases[i].UpdateQuery.ToSQLWithArgs(testCases[i].Dialect)

			if testCases[i].Expectation.Error != nil && actualErr == nil {
				t.Error("expectation error is not nil, got nil")
			}

			if testCases[i].Expectation.Error == nil && actualErr != nil {
				t.Error("expectation error is nil, got not nil")
			}

			if testCases[i].Expectation.Error != nil && actualErr != nil && testCases[i].Expectation.Error.Error() != actualErr.Error() {
				t.Errorf("expectation error is %s, got %s", testCases[i].Expectation.Error.Error(), actualErr.Error())
			}

			if testCases[i].Expectation.Query != actualQuery {
				t.Errorf("expectation query is %s, got %s", testCases[i].Expectation.Query, actualQuery)
			}

			if len(testCases[i].Expectation.Args) != len(actualArgs) {
				t.Errorf("expectation length of args is %d, got %d", len(testCases[i].Expectation.Args), len(actualArgs))
			}

			for j := 0; j < len(testCases[i].Expectation.Args); j++ {
				if !deepEqual(testCases[i].Expectation.Args[j], actualArgs[j]) {
					t.Errorf("expectation element of args is %v, got %v", testCases[i].Expectation.Args[j], actualArgs[j])
				}
			}
		})
	}
}
