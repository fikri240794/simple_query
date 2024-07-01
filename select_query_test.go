package simple_query

import (
	"errors"
	"fmt"
	"testing"
)

func TestSelectQuery_Select(t *testing.T) {
	var (
		expectation *SelectQuery
		actual      *SelectQuery
	)

	expectation = &SelectQuery{
		Fields: []string{"field1", "field2", "field3"},
	}

	actual = Select("field1", "field2", "field3")

	if len(expectation.Fields) != len(actual.Fields) {
		t.Errorf("expectation length of fields is %d, got %d", len(expectation.Fields), len(actual.Fields))
	}

	for i := range expectation.Fields {
		if expectation.Fields[i] != actual.Fields[i] {
			t.Errorf("expectation element of fields is %s, got %s", expectation.Fields[i], actual.Fields[i])
		}
	}
}

func TestSelectQuery_From(t *testing.T) {
	var (
		expectation *SelectQuery
		actual      *SelectQuery
	)

	expectation = &SelectQuery{
		Fields: []string{"field1", "field2", "field3"},
		Table:  "table1",
	}

	actual = Select("field1", "field2", "field3").
		From("table1")

	if len(expectation.Fields) != len(actual.Fields) {
		t.Errorf("expectation length of fields is %d, got %d", len(expectation.Fields), len(actual.Fields))
	}

	for i := range expectation.Fields {
		if expectation.Fields[i] != actual.Fields[i] {
			t.Errorf("expectation element of fields is %s, got %s", expectation.Fields[i], actual.Fields[i])
		}
	}

	if expectation.Table != actual.Table {
		t.Errorf("expectation table is %s, got %s", expectation.Table, actual.Table)
	}
}

func TestSelectQuery_Where(t *testing.T) {
	var (
		expectation *SelectQuery
		actual      *SelectQuery
	)

	expectation = &SelectQuery{
		Fields: []string{"field1", "field2", "field3"},
		Table:  "table1",
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
	}

	actual = Select("field1", "field2", "field3").
		From("table1").
		Where(
			NewFilter().
				SetLogic(LogicAnd).
				AddFilter("field1", OperatorEqual, "value1"),
		)

	if len(expectation.Fields) != len(actual.Fields) {
		t.Errorf("expectation length of fields is %d, got %d", len(expectation.Fields), len(actual.Fields))
	}

	for i := range expectation.Fields {
		if expectation.Fields[i] != actual.Fields[i] {
			t.Errorf("expectation element of fields is %s, got %s", expectation.Fields[i], actual.Fields[i])
		}
	}

	if expectation.Table != actual.Table {
		t.Errorf("expectation table is %s, got %s", expectation.Table, actual.Table)
	}

	if !deepEqual(expectation.Filter, actual.Filter) {
		t.Errorf("expectation filter is %v, got %v", expectation.Filter, actual.Filter)
	}
}

func TestSelectQuery_OrderBy(t *testing.T) {
	var (
		expectation *SelectQuery
		actual      *SelectQuery
	)

	expectation = &SelectQuery{
		Fields: []string{"field1", "field2", "field3"},
		Table:  "table1",
		Sorts: []*Sort{
			{
				Field:     "field1",
				Direction: SortDirectionDescending,
			},
			{
				Field:     "field2",
				Direction: SortDirectionAscending,
			},
		},
	}

	actual = Select("field1", "field2", "field3").
		From("table1").
		OrderBy(
			NewSort("field1", SortDirectionDescending),
			NewSort("field2", SortDirectionAscending),
		)

	if len(expectation.Fields) != len(actual.Fields) {
		t.Errorf("expectation length of fields is %d, got %d", len(expectation.Fields), len(actual.Fields))
	}

	for i := range expectation.Fields {
		if expectation.Fields[i] != actual.Fields[i] {
			t.Errorf("expectation element of fields is %s, got %s", expectation.Fields[i], actual.Fields[i])
		}
	}

	if expectation.Table != actual.Table {
		t.Errorf("expectation table is %s, got %s", expectation.Table, actual.Table)
	}

	if len(expectation.Sorts) != len(actual.Sorts) {
		t.Errorf("expectation length of sorts is %d, got %d", len(expectation.Sorts), len(actual.Sorts))
	}

	for i := range expectation.Sorts {
		if expectation.Sorts[i] == nil && actual.Sorts[i] != nil {
			t.Errorf("expectation element of sorts is nil, got %v", actual.Sorts[i])
		}

		if expectation.Sorts[i] != nil && actual.Sorts[i] == nil {
			t.Errorf("expectation element of sorts is %v, got nil", expectation.Sorts[i])
		}

		if expectation.Sorts[i].Field != actual.Sorts[i].Field {
			t.Errorf("expectation field of sorts element is %s, got %s", expectation.Sorts[i].Field, actual.Sorts[i].Field)
		}

		if string(expectation.Sorts[i].Direction) != string(actual.Sorts[i].Direction) {
			t.Errorf("expectation direction of sorts element is %s, got %s", expectation.Sorts[i].Direction, actual.Sorts[i].Direction)
		}
	}
}

func TestSelectQuery_Limit(t *testing.T) {
	var (
		expectation *SelectQuery
		actual      *SelectQuery
	)

	expectation = &SelectQuery{
		Fields: []string{"field1", "field2", "field3"},
		Table:  "table1",
		Take:   10,
	}

	actual = Select("field1", "field2", "field3").
		From("table1").
		Limit(10)

	if len(expectation.Fields) != len(actual.Fields) {
		t.Errorf("expectation length of fields is %d, got %d", len(expectation.Fields), len(actual.Fields))
	}

	for i := range expectation.Fields {
		if expectation.Fields[i] != actual.Fields[i] {
			t.Errorf("expectation element of fields is %s, got %s", expectation.Fields[i], actual.Fields[i])
		}
	}

	if expectation.Table != actual.Table {
		t.Errorf("expectation table is %s, got %s", expectation.Table, actual.Table)
	}

	if expectation.Take != actual.Take {
		t.Errorf("expectation take is %d, got %d", expectation.Take, actual.Take)
	}
}

func TestSelectQuery_validate(t *testing.T) {
	var testCases []struct {
		Name        string
		SelectQuery *SelectQuery
		Expectation error
	} = []struct {
		Name        string
		SelectQuery *SelectQuery
		Expectation error
	}{
		{
			Name:        "fields is empty",
			SelectQuery: &SelectQuery{},
			Expectation: errors.New("fields is required"),
		},
		{
			Name: "field is empty",
			SelectQuery: &SelectQuery{
				Fields: []string{""},
			},
			Expectation: errors.New("field is required"),
		},
		{
			Name: "table is empty",
			SelectQuery: &SelectQuery{
				Fields: []string{"field1", "field2", "field3"},
			},
			Expectation: errors.New("table is required"),
		},
		{
			Name: "select query is valid",
			SelectQuery: &SelectQuery{
				Fields: []string{"field1", "field2", "field3"},
				Table:  "table1",
				Take:   10,
			},
			Expectation: nil,
		},
	}

	for i := range testCases {
		t.Run(testCases[i].Name, func(t *testing.T) {
			var actual error = testCases[i].SelectQuery.validate()

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

func TestSelectQuery_ToSQLWithArgs(t *testing.T) {
	var testCases []struct {
		Name        string
		SelectQuery *SelectQuery
		Dialect     Dialect
		Expectation struct {
			Query string
			Args  []interface{}
			Error error
		}
	} = []struct {
		Name        string
		SelectQuery *SelectQuery
		Dialect     Dialect
		Expectation struct {
			Query string
			Args  []interface{}
			Error error
		}
	}{
		{
			Name:        "fields is empty",
			SelectQuery: &SelectQuery{},
			Dialect:     "",
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
			SelectQuery: &SelectQuery{
				Fields: []string{""},
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
			Name: "table is empty",
			SelectQuery: &SelectQuery{
				Fields: []string{"field1", "field2", "field3"},
			},
			Dialect: "",
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
			Name: fmt.Sprintf("dialect %s with invalid filter", DialectMySQL),
			SelectQuery: &SelectQuery{
				Fields: []string{"field1", "field2", "field3"},
				Table:  "table1",
				Filter: &Filter{
					Logic:   LogicAnd,
					Filters: []*Filter{},
				},
				Take: 100,
			},
			Dialect: DialectMySQL,
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
			Name: fmt.Sprintf("dialect %s with invalid sort", DialectMySQL),
			SelectQuery: &SelectQuery{
				Fields: []string{"field1", "field2", "field3"},
				Table:  "table1",
				Sorts: []*Sort{
					{
						Field:     "",
						Direction: SortDirectionDescending,
					},
				},
				Take: 100,
			},
			Dialect: DialectMySQL,
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
			Name: fmt.Sprintf("dialect %s with element sorts is nil", DialectMySQL),
			SelectQuery: &SelectQuery{
				Fields: []string{"field1", "field2", "field3"},
				Table:  "table1",
				Sorts: []*Sort{
					nil,
				},
				Take: 100,
			},
			Dialect: DialectMySQL,
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "select field1, field2, field3 from table1 limit ?",
				Args:  []interface{}{100},
				Error: nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with take", DialectMySQL),
			SelectQuery: &SelectQuery{
				Fields: []string{"field1", "field2", "field3"},
				Table:  "table1",
				Take:   10,
			},
			Dialect: DialectMySQL,
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "select field1, field2, field3 from table1 limit ?",
				Args:  []interface{}{10},
				Error: nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter sort take", DialectMySQL),
			SelectQuery: &SelectQuery{
				Fields: []string{"field1", "field2", "field3"},
				Table:  "table1",
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
				Sorts: []*Sort{
					{
						Field:     "field1",
						Direction: SortDirectionDescending,
					},
				},
				Take: 10,
			},
			Dialect: DialectMySQL,
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "select field1, field2, field3 from table1 where field1 = ? order by field1 desc limit ?",
				Args:  []interface{}{"value1", 10},
				Error: nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter take", DialectMySQL),
			SelectQuery: &SelectQuery{
				Fields: []string{"field1", "field2", "field3"},
				Table:  "table1",
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
				Take: 10,
			},
			Dialect: DialectMySQL,
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "select field1, field2, field3 from table1 where field1 = ? limit ?",
				Args:  []interface{}{"value1", 10},
				Error: nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with sort take", DialectMySQL),
			SelectQuery: &SelectQuery{
				Fields: []string{"field1", "field2", "field3"},
				Table:  "table1",
				Sorts: []*Sort{
					{
						Field:     "field1",
						Direction: SortDirectionDescending,
					},
				},
				Take: 10,
			},
			Dialect: DialectMySQL,
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "select field1, field2, field3 from table1 order by field1 desc limit ?",
				Args:  []interface{}{10},
				Error: nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with invalid filter", DialectPostgres),
			SelectQuery: &SelectQuery{
				Fields: []string{"field1", "field2", "field3"},
				Table:  "table1",
				Filter: &Filter{
					Logic:   LogicAnd,
					Filters: []*Filter{},
				},
				Take: 100,
			},
			Dialect: DialectPostgres,
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
			Name: fmt.Sprintf("dialect %s with invalid sort", DialectPostgres),
			SelectQuery: &SelectQuery{
				Fields: []string{"field1", "field2", "field3"},
				Table:  "table1",
				Sorts: []*Sort{
					{
						Field:     "",
						Direction: SortDirectionDescending,
					},
				},
				Take: 100,
			},
			Dialect: DialectPostgres,
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
			Name: fmt.Sprintf("dialect %s with element sorts is nil", DialectPostgres),
			SelectQuery: &SelectQuery{
				Fields: []string{"field1", "field2", "field3"},
				Table:  "table1",
				Sorts: []*Sort{
					nil,
				},
				Take: 100,
			},
			Dialect: DialectPostgres,
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "select field1, field2, field3 from table1 limit $1",
				Args:  []interface{}{100},
				Error: nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with take", DialectPostgres),
			SelectQuery: &SelectQuery{
				Fields: []string{"field1", "field2", "field3"},
				Table:  "table1",
				Take:   10,
			},
			Dialect: DialectPostgres,
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "select field1, field2, field3 from table1 limit $1",
				Args:  []interface{}{10},
				Error: nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter sort take", DialectPostgres),
			SelectQuery: &SelectQuery{
				Fields: []string{"field1", "field2", "field3"},
				Table:  "table1",
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
				Sorts: []*Sort{
					{
						Field:     "field1",
						Direction: SortDirectionDescending,
					},
				},
				Take: 10,
			},
			Dialect: DialectPostgres,
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "select field1, field2, field3 from table1 where field1 = $1 order by field1 desc limit $2",
				Args:  []interface{}{"value1", 10},
				Error: nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter take", DialectPostgres),
			SelectQuery: &SelectQuery{
				Fields: []string{"field1", "field2", "field3"},
				Table:  "table1",
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
				Take: 10,
			},
			Dialect: DialectPostgres,
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "select field1, field2, field3 from table1 where field1 = $1 limit $2",
				Args:  []interface{}{"value1", 10},
				Error: nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with sort take", DialectPostgres),
			SelectQuery: &SelectQuery{
				Fields: []string{"field1", "field2", "field3"},
				Table:  "table1",
				Sorts: []*Sort{
					{
						Field:     "field1",
						Direction: SortDirectionDescending,
					},
				},
				Take: 10,
			},
			Dialect: DialectPostgres,
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "select field1, field2, field3 from table1 order by field1 desc limit $1",
				Args:  []interface{}{10},
				Error: nil,
			},
		},
	}

	for i := range testCases {
		t.Run(testCases[i].Name, func(t *testing.T) {
			var (
				actualQuery string
				actualArgs  []interface{}
				actualErr   error
			)

			actualQuery, actualArgs, actualErr = testCases[i].SelectQuery.ToSQLWithArgs(testCases[i].Dialect)

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

			for j := range testCases[i].Expectation.Args {
				if !deepEqual(testCases[i].Expectation.Args[j], actualArgs[j]) {
					t.Errorf("expectation element of args is %v, got %v", testCases[i].Expectation.Args[j], actualArgs[j])
				}
			}
		})
	}
}
