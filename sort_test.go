package simple_query

import (
	"errors"
	"testing"
)

func TestSort_NewSort(t *testing.T) {
	var (
		expectation *Sort
		actual      *Sort
	)

	expectation = &Sort{
		Field:     "field1",
		Direction: SortDirectionAscending,
	}

	actual = NewSort("field1", SortDirectionAscending)

	if expectation.Field != actual.Field {
		t.Errorf("expectation field is %s, got %s", expectation.Field, actual.Field)
	}

	if expectation.Direction != actual.Direction {
		t.Errorf("expectation direction is %s, got %s", expectation.Direction, actual.Direction)
	}
}

func TestSort_validate(t *testing.T) {
	var testCases []struct {
		Name        string
		Sort        *Sort
		Expectation error
	} = []struct {
		Name        string
		Sort        *Sort
		Expectation error
	}{
		{
			Name:        "field is empty",
			Sort:        &Sort{},
			Expectation: errors.New("field is required"),
		},
		{
			Name: "sort is valid",
			Sort: &Sort{
				Field:     "field1",
				Direction: SortDirectionDescending,
			},
			Expectation: nil,
		},
	}

	for i := range testCases {
		t.Run(testCases[i].Name, func(t *testing.T) {
			var actual error = testCases[i].Sort.validate()

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

func TestSort_ToSQL(t *testing.T) {
	var testCases []struct {
		Name        string
		Sort        *Sort
		Expectation struct {
			Query string
			Error error
		}
	} = []struct {
		Name        string
		Sort        *Sort
		Expectation struct {
			Query string
			Error error
		}
	}{
		{
			Name: "field is empty",
			Sort: &Sort{},
			Expectation: struct {
				Query string
				Error error
			}{
				Query: "",
				Error: errors.New("field is required"),
			},
		},
		{
			Name: "default direction",
			Sort: &Sort{
				Field: "field1",
			},
			Expectation: struct {
				Query string
				Error error
			}{
				Query: "field1 asc",
				Error: nil,
			},
		},
		{
			Name: "sort with direction",
			Sort: &Sort{
				Field:     "field1",
				Direction: SortDirectionDescending,
			},
			Expectation: struct {
				Query string
				Error error
			}{
				Query: "field1 desc",
				Error: nil,
			},
		},
	}

	for i := range testCases {
		t.Run(testCases[i].Name, func(t *testing.T) {
			var (
				actualQuery string
				actualErr   error
			)

			actualQuery, actualErr = testCases[i].Sort.ToSQL()

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
		})
	}
}
