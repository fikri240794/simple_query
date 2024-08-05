package simple_query

import (
	"testing"
)

func testSort_SortEquality(t *testing.T, expectation, actual *Sort) {
	if expectation == nil && actual == nil {
		t.Skip("expectation and actual is nil")
	}

	if expectation == nil && actual != nil {
		t.Errorf("expectation is nil, got %+v", actual)
	}

	if expectation != nil && actual == nil {
		t.Errorf("expectation is %+v, got nil", expectation)
	}

	if expectation.Field != actual.Field {
		t.Errorf("expectation field is %s, got %s", expectation.Field, actual.Field)
	}

	if expectation.Direction != actual.Direction {
		t.Errorf("expectation direction is %s, got %s", expectation.Direction, actual.Direction)
	}
}

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

	testSort_SortEquality(t, expectation, actual)
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
			Expectation: ErrFieldIsRequired,
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
			Err   error
		}
	} = []struct {
		Name        string
		Sort        *Sort
		Expectation struct {
			Query string
			Err   error
		}
	}{
		{
			Name: "field is empty",
			Sort: &Sort{},
			Expectation: struct {
				Query string
				Err   error
			}{
				Query: "",
				Err:   ErrFieldIsRequired,
			},
		},
		{
			Name: "default direction",
			Sort: &Sort{
				Field: "field1",
			},
			Expectation: struct {
				Query string
				Err   error
			}{
				Query: "field1 asc",
				Err:   nil,
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
				Err   error
			}{
				Query: "field1 desc",
				Err:   nil,
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
		})
	}
}
