package simple_query

import "testing"

func Test_getPlaceholder(t *testing.T) {
	var testCases []struct {
		Name        string
		Dialect     Dialect
		StartIdx    int
		EndIdx      int
		Expectation string
	} = []struct {
		Name        string
		Dialect     Dialect
		StartIdx    int
		EndIdx      int
		Expectation string
	}{
		{
			Name:        "unknown dialect",
			Dialect:     "unknown",
			StartIdx:    1,
			EndIdx:      1,
			Expectation: "",
		},
		{
			Name:        "zero start index",
			Dialect:     DialectMySQL,
			StartIdx:    0,
			EndIdx:      1,
			Expectation: "",
		},
		{
			Name:        "zero end index",
			Dialect:     DialectMySQL,
			StartIdx:    1,
			EndIdx:      0,
			Expectation: "",
		},
		{
			Name:        "end index less than start index",
			Dialect:     DialectMySQL,
			StartIdx:    1,
			EndIdx:      0,
			Expectation: "",
		},
		{
			Name:        "mysql with start index equal to end index",
			Dialect:     DialectMySQL,
			StartIdx:    1,
			EndIdx:      1,
			Expectation: "?",
		},
		{
			Name:        "mysql with start index less than end index",
			Dialect:     DialectMySQL,
			StartIdx:    1,
			EndIdx:      5,
			Expectation: "?, ?, ?, ?, ?",
		},
		{
			Name:        "postgres with start index equal to end index",
			Dialect:     DialectPostgres,
			StartIdx:    1,
			EndIdx:      1,
			Expectation: "$1",
		},
		{
			Name:        "postgres with start index less than end index",
			Dialect:     DialectPostgres,
			StartIdx:    1,
			EndIdx:      5,
			Expectation: "$1, $2, $3, $4, $5",
		},
	}

	for i := 0; i < len(testCases); i++ {
		t.Run(testCases[i].Name, func(t *testing.T) {
			var actual string = getPlaceholder(testCases[i].Dialect, testCases[i].StartIdx, testCases[i].EndIdx)
			if testCases[i].Expectation != actual {
				t.Errorf("expected placeholder %s, got %s", testCases[i].Expectation, actual)
			}
		})
	}
}
