package simple_query

import (
	"fmt"
	"testing"
)

func testSelectQuery_SelectQueryEquality(t *testing.T, expectation, actual *SelectQuery) {
	if len(expectation.Fields) != len(actual.Fields) {
		t.Errorf("expectation length of fields is %d, got %d", len(expectation.Fields), len(actual.Fields))
	} else {
		for i := range expectation.Fields {
			if !deepEqual(expectation.Fields[i], actual.Fields[i]) {
				t.Errorf("expectation element of fields is %+v, got %+v", expectation.Fields[i], actual.Fields[i])
			}
		}
	}

	if expectation.Table != nil && actual.Table == nil {
		t.Errorf("expectation table is %+v, got nil", expectation.Table)
	}
	if expectation.Table == nil && actual.Table != nil {
		t.Errorf("expectation table is nil, got %+v", actual.Table)
	}
	if !deepEqual(expectation.Table, actual.Table) {
		t.Errorf("expectation table is %+v, got %+v", expectation.Table, actual.Table)
	}

	if expectation.Filter != nil && actual.Filter == nil {
		t.Errorf("expectation filter is %+v, got nil", expectation.Filter)
	}
	if expectation.Filter == nil && actual.Filter != nil {
		t.Errorf("expectation filter is nil, got %+v", actual.Filter)
	}
	if !deepEqual(expectation.Filter, actual.Filter) {
		t.Errorf("expectation table is %+v, got %+v", expectation.Filter, actual.Filter)
	}

	if len(expectation.Sorts) != len(actual.Sorts) {
		t.Errorf("expectation length of sorts is %d, got %d", len(expectation.Sorts), len(actual.Sorts))
	} else {
		for i := range expectation.Sorts {
			if !deepEqual(expectation.Sorts[i], actual.Sorts[i]) {
				t.Errorf("expectation element of sorts is %+v, got %+v", expectation.Sorts[i], actual.Sorts[i])
			}
		}
	}

	if expectation.Take != actual.Take {
		t.Errorf("expectation take is %d, got %d", expectation.Take, actual.Take)
	}

	if expectation.Alias != actual.Alias {
		t.Errorf("expectation alias is %s, got %s", expectation.Alias, actual.Alias)
	}
}

func TestSelectQuery_Select(t *testing.T) {
	var (
		expectation *SelectQuery
		actual      *SelectQuery
	)

	expectation = &SelectQuery{
		Fields: []*Field{
			{
				Column: "field1",
			},
			{
				Column: "field2",
			},
			{
				Column: "field3",
			},
		},
	}

	actual = Select(NewField("field1"), NewField("field2"), NewField("field3"))

	testSelectQuery_SelectQueryEquality(t, expectation, actual)
}

func TestSelectQuery_From(t *testing.T) {
	var (
		expectation *SelectQuery
		actual      *SelectQuery
	)

	expectation = &SelectQuery{
		Fields: []*Field{
			{
				Column: "field1",
			},
			{
				Column: "field2",
			},
			{
				Column: "field3",
			},
		},
		Table: &Table{
			Name: "table1",
		},
	}

	actual = Select(NewField("field1"), NewField("field2"), NewField("field3")).
		From(NewTable("table1"))

	testSelectQuery_SelectQueryEquality(t, expectation, actual)
}

func TestSelectQuery_Where(t *testing.T) {
	var (
		expectation *SelectQuery
		actual      *SelectQuery
	)

	expectation = &SelectQuery{
		Fields: []*Field{
			{
				Column: "field1",
			},
			{
				Column: "field2",
			},
			{
				Column: "field3",
			},
		},
		Table: &Table{
			Name: "table1",
		},
		Filter: &Filter{
			Logic: LogicAnd,
			Filters: []*Filter{
				{
					Field: &Field{
						Column: "field1",
					},
					Operator: OperatorEqual,
					Value:    "value1",
				},
			},
		},
	}

	actual = Select(NewField("field1"), NewField("field2"), NewField("field3")).
		From(NewTable("table1")).
		Where(
			NewFilter().
				SetLogic(LogicAnd).
				AddFilter(NewField("field1"), OperatorEqual, "value1"),
		)

	testSelectQuery_SelectQueryEquality(t, expectation, actual)
}

func TestSelectQuery_OrderBy(t *testing.T) {
	var (
		expectation *SelectQuery
		actual      *SelectQuery
	)

	expectation = &SelectQuery{
		Fields: []*Field{
			{
				Column: "field1",
			},
			{
				Column: "field2",
			},
			{
				Column: "field3",
			},
		},
		Table: &Table{
			Name: "table1",
		},
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

	actual = Select(NewField("field1"), NewField("field2"), NewField("field3")).
		From(NewTable("table1")).
		OrderBy(
			NewSort("field1", SortDirectionDescending),
			NewSort("field2", SortDirectionAscending),
		)

	testSelectQuery_SelectQueryEquality(t, expectation, actual)
}

func TestSelectQuery_Limit(t *testing.T) {
	var (
		expectation *SelectQuery
		actual      *SelectQuery
	)

	expectation = &SelectQuery{
		Fields: []*Field{
			{
				Column: "field1",
			},
			{
				Column: "field2",
			},
			{
				Column: "field3",
			},
		},
		Table: &Table{
			Name: "table1",
		},
		Take: 10,
	}

	actual = Select(NewField("field1"), NewField("field2"), NewField("field3")).
		From(NewTable("table1")).
		Limit(10)

	testSelectQuery_SelectQueryEquality(t, expectation, actual)
}

func TestSelectQuery_As(t *testing.T) {
	var (
		expectation *SelectQuery
		actual      *SelectQuery
	)

	expectation = &SelectQuery{
		Fields: []*Field{
			{
				Column: "field1",
			},
			{
				Column: "field2",
			},
			{
				Column: "field3",
			},
		},
		Table: &Table{
			Name: "table1",
		},
		Alias: "alias1",
	}

	actual = Select(NewField("field1"), NewField("field2"), NewField("field3")).
		From(NewTable("table1")).
		As("alias1")

	testSelectQuery_SelectQueryEquality(t, expectation, actual)
}

func TestSelectQuery_validate(t *testing.T) {
	var testCases []struct {
		Name        string
		Dialect     Dialect
		SelectQuery *SelectQuery
		Expectation error
	} = []struct {
		Name        string
		Dialect     Dialect
		SelectQuery *SelectQuery
		Expectation error
	}{
		{
			Name:        "dialect is empty",
			Dialect:     "",
			SelectQuery: &SelectQuery{},
			Expectation: ErrDialectIsRequired,
		},
		{
			Name:        "fields is empty",
			Dialect:     DialectPostgres,
			SelectQuery: &SelectQuery{},
			Expectation: ErrFieldsIsRequired,
		},
		{
			Name:    "fields element is nil",
			Dialect: DialectPostgres,
			SelectQuery: &SelectQuery{
				Fields: []*Field{nil},
			},
			Expectation: ErrFieldIsNil,
		},
		{
			Name:    "table is nil",
			Dialect: DialectPostgres,
			SelectQuery: &SelectQuery{
				Fields: []*Field{
					{
						Column: "field1",
					},
				},
			},
			Expectation: ErrTableIsRequired,
		},
		{
			Name:    "select query is valid",
			Dialect: DialectPostgres,
			SelectQuery: &SelectQuery{
				Fields: []*Field{
					{
						Column: "field1",
					},
				},
				Table: &Table{
					Name: "table1",
				},
			},
			Expectation: nil,
		},
	}

	for i := range testCases {
		t.Run(testCases[i].Name, func(t *testing.T) {
			var actual error = testCases[i].SelectQuery.validate(testCases[i].Dialect)

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
			Dialect:     DialectPostgres,
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "",
				Args:  nil,
				Error: ErrFieldsIsRequired,
			},
		},
		{
			Name: "fields is not empty and fields element is not nil and fields element to sql with args with alias is error",
			SelectQuery: &SelectQuery{
				Fields: []*Field{
					{},
				},
				Table: &Table{},
			},
			Dialect: DialectPostgres,
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "",
				Args:  nil,
				Error: ErrColumnIsRequired,
			},
		},
		{
			Name: "fields is not empty and fields element is not nil",
			SelectQuery: &SelectQuery{
				Fields: []*Field{
					{
						Column: "field1",
					},
				},
				Table: &Table{
					Name: "table1",
				},
			},
			Dialect: DialectPostgres,
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "select field1 from table1",
				Args:  nil,
				Error: nil,
			},
		},
		{
			Name: "table is not nil and table is to sql with args with alias is error",
			SelectQuery: &SelectQuery{
				Fields: []*Field{
					{
						Column: "field1",
					},
				},
				Table: &Table{},
			},
			Dialect: DialectPostgres,
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "",
				Args:  nil,
				Error: ErrNameIsRequired,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with invalid filter", DialectPostgres),
			SelectQuery: &SelectQuery{
				Fields: []*Field{
					{
						Column: "field1",
					},
				},
				Table: &Table{
					Name: "table1",
				},
				Filter: &Filter{
					Logic:   LogicAnd,
					Filters: []*Filter{},
				},
			},
			Dialect: DialectPostgres,
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "",
				Args:  nil,
				Error: ErrFiltersIsRequired,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with filter", DialectPostgres),
			SelectQuery: &SelectQuery{
				Fields: []*Field{
					{
						Column: "field1",
					},
				},
				Table: &Table{
					Name: "table1",
				},
				Filter: &Filter{
					Logic: LogicAnd,
					Filters: []*Filter{
						{
							Field: &Field{
								Column: "field1",
							},
							Operator: OperatorEqual,
							Value:    "value1",
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
				Query: "select field1 from table1 where field1 = $1",
				Args:  []interface{}{"value1"},
				Error: nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with element sorts is nil", DialectPostgres),
			SelectQuery: &SelectQuery{
				Fields: []*Field{
					{
						Column: "field1",
					},
				},
				Table: &Table{
					Name: "table1",
				},
				Sorts: []*Sort{
					nil,
				},
			},
			Dialect: DialectPostgres,
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "select field1 from table1",
				Args:  []interface{}{},
				Error: nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with invalid sort", DialectPostgres),
			SelectQuery: &SelectQuery{
				Fields: []*Field{
					{
						Column: "field1",
					},
				},
				Table: &Table{
					Name: "table1",
				},
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
				Error: ErrFieldIsRequired,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with sort", DialectPostgres),
			SelectQuery: &SelectQuery{
				Fields: []*Field{
					{
						Column: "field1",
					},
				},
				Table: &Table{
					Name: "table1",
				},
				Sorts: []*Sort{
					{
						Field:     "field1",
						Direction: SortDirectionDescending,
					},
				},
			},
			Dialect: DialectPostgres,
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "select field1 from table1 order by field1 desc",
				Args:  []interface{}{},
				Error: nil,
			},
		},
		{
			Name: fmt.Sprintf("dialect %s with take", DialectPostgres),
			SelectQuery: &SelectQuery{
				Fields: []*Field{
					{
						Column: "field1",
					},
				},
				Table: &Table{
					Name: "table1",
				},
				Take: 10,
			},
			Dialect: DialectPostgres,
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "select field1 from table1 limit $1",
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

			actualQuery, actualArgs, actualErr = testCases[i].SelectQuery.ToSQLWithArgs(testCases[i].Dialect, []interface{}{})

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
					t.Errorf("expectation element of args is %+v, got %+v", testCases[i].Expectation.Args[j], actualArgs[j])
				}
			}
		})
	}
}

func TestSelectQuery_ToSQLWithArgsWithAlias(t *testing.T) {
	var testCases []struct {
		Name        string
		Dialect     Dialect
		SelectQuery *SelectQuery
		Expectation struct {
			Query string
			Args  []interface{}
			Err   error
		}
	} = []struct {
		Name        string
		Dialect     Dialect
		SelectQuery *SelectQuery
		Expectation struct {
			Query string
			Args  []interface{}
			Err   error
		}
	}{
		{
			Name:        "to sql with args is error",
			Dialect:     DialectPostgres,
			SelectQuery: &SelectQuery{},
			Expectation: struct {
				Query string
				Args  []interface{}
				Err   error
			}{
				Query: "",
				Args:  nil,
				Err:   ErrFieldsIsRequired,
			},
		},
		{
			Name:    "alias is not empty",
			Dialect: DialectPostgres,
			SelectQuery: &SelectQuery{
				Fields: []*Field{
					{
						Column: "field1",
					},
				},
				Table: &Table{
					Name: "table1",
				},
				Alias: "alias1",
			},
			Expectation: struct {
				Query string
				Args  []interface{}
				Err   error
			}{
				Query: "(select field1 from table1) as alias1",
				Args:  []interface{}{},
				Err:   nil,
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

			actualQuery, actualArgs, actualErr = testCases[i].SelectQuery.ToSQLWithArgsWithAlias(testCases[i].Dialect, []interface{}{})

			if testCases[i].Expectation.Err != nil && actualErr == nil {
				t.Error("expectation error is not nil, got nil")
			}

			if testCases[i].Expectation.Err == nil && actualErr != nil {
				t.Error("expectation error is nil, got not nil")
			}

			if testCases[i].Expectation.Err != nil && actualErr != nil && testCases[i].Expectation.Err.Error() != actualErr.Error() {
				t.Errorf("expectation error is %s, got %s", testCases[i].Expectation.Err.Error(), actualErr.Error())
			}

			if testCases[i].Expectation.Query != actualQuery {
				t.Errorf("expectation query is %s, got %s", testCases[i].Expectation.Query, actualQuery)
			}

			if len(testCases[i].Expectation.Args) != len(actualArgs) {
				t.Errorf("expectation length of args is %d, got %d", len(testCases[i].Expectation.Args), len(actualArgs))
			}

			for j := range testCases[i].Expectation.Args {
				if !deepEqual(testCases[i].Expectation.Args[j], actualArgs[j]) {
					t.Errorf("expectation element of args is %+v, got %+v", testCases[i].Expectation.Args[j], actualArgs[j])
				}
			}
		})
	}
}
