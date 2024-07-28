package simple_query

import (
	"fmt"
	"testing"
)

func TestDeleteQuery_Delete(t *testing.T) {
	var (
		expectation *DeleteQuery
		actual      *DeleteQuery
	)

	expectation = &DeleteQuery{}
	actual = Delete()

	if !deepEqual(expectation, actual) {
		t.Errorf("expectation delete query is %v, got %v", expectation, actual)
	}
}

func TestDeleteQuery_From(t *testing.T) {
	var (
		expectation *DeleteQuery
		actual      *DeleteQuery
	)

	expectation = &DeleteQuery{
		Table: "table1",
	}
	actual = Delete().
		From("table1")

	if expectation.Table != actual.Table {
		t.Errorf("expectation table is %s, got %s", expectation.Table, actual.Table)
	}
}

func TestDeleteQuery_Where(t *testing.T) {
	var (
		expectation *DeleteQuery
		actual      *DeleteQuery
	)

	expectation = &DeleteQuery{
		Table: "table1",
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

	actual = Delete().
		From("table1").
		Where(
			NewFilter().
				SetLogic(LogicAnd).
				AddFilter(NewField("field1"), OperatorEqual, "value1"),
		)

	if expectation.Table != actual.Table {
		t.Errorf("expectation table is %s, got %s", expectation.Table, actual.Table)
	}

	if !deepEqual(expectation.Filter, actual.Filter) {
		t.Errorf("expectation filter is %v, got %v", expectation.Filter, actual.Filter)
	}
}

func TestDeleteQuery_validate(t *testing.T) {
	var testCases []struct {
		Name        string
		Dialect     Dialect
		DeleteQuery *DeleteQuery
		Expectation error
	} = []struct {
		Name        string
		Dialect     Dialect
		DeleteQuery *DeleteQuery
		Expectation error
	}{
		{
			Name:        "dialeg is empty",
			Dialect:     "",
			DeleteQuery: &DeleteQuery{},
			Expectation: ErrDialectIsRequired,
		},
		{
			Name:        "table is empty",
			Dialect:     DialectPostgres,
			DeleteQuery: &DeleteQuery{},
			Expectation: ErrTableIsRequired,
		},
		{
			Name:    "filter is empty",
			Dialect: DialectPostgres,
			DeleteQuery: &DeleteQuery{
				Table: "table1",
			},
			Expectation: ErrFilterIsRequired,
		},
		{
			Name:    "delete query is valid",
			Dialect: DialectPostgres,
			DeleteQuery: &DeleteQuery{
				Table: "table1",
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
			Expectation: nil,
		},
	}

	for i := range testCases {
		t.Run(testCases[i].Name, func(t *testing.T) {
			var actual error = testCases[i].DeleteQuery.validate(testCases[i].Dialect)

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

func TestDeleteQuery_ToSQLWithArgs(t *testing.T) {
	var testCases []struct {
		Name        string
		DeleteQuery *DeleteQuery
		Dialect     Dialect
		Expectation struct {
			Query string
			Args  []interface{}
			Error error
		}
	} = []struct {
		Name        string
		DeleteQuery *DeleteQuery
		Dialect     Dialect
		Expectation struct {
			Query string
			Args  []interface{}
			Error error
		}
	}{
		{
			Name:        "delete query is invalid",
			DeleteQuery: &DeleteQuery{},
			Dialect:     DialectPostgres,
			Expectation: struct {
				Query string
				Args  []interface{}
				Error error
			}{
				Query: "",
				Args:  nil,
				Error: ErrTableIsRequired,
			},
		},
		{
			Name: fmt.Sprintf("delete query with dialect %s and filter to sql args is error", DialectPostgres),
			DeleteQuery: &DeleteQuery{
				Table:  "table1",
				Filter: &Filter{},
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
			Name: fmt.Sprintf("delete query with dialect %s", DialectPostgres),
			DeleteQuery: &DeleteQuery{
				Table: "table1",
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
				Query: "delete from table1 where field1 = $1",
				Args:  []interface{}{"value1"},
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

			actualQuery, actualArgs, actualErr = testCases[i].DeleteQuery.ToSQLWithArgs(testCases[i].Dialect)

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
					t.Errorf("expectation element of args is %v, got %v", testCases[i].Expectation.Args[j], actualArgs[i])
				}
			}
		})
	}
}
