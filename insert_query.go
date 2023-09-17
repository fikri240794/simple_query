package simple_query

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type InsertQuery struct {
	Table        string
	FieldsValues map[string][]interface{}
}

func Insert() *InsertQuery {
	return &InsertQuery{
		FieldsValues: map[string][]interface{}{},
	}
}

func (i *InsertQuery) Into(table string) *InsertQuery {
	i.Table = table
	return i
}

func (i *InsertQuery) Value(field string, value interface{}) *InsertQuery {
	i.FieldsValues[field] = append(i.FieldsValues[field], value)
	return i
}

func (i *InsertQuery) getColumnsAndRowsValues() ([]string, [][]interface{}) {
	var (
		columns    []string
		rowCount   int
		rowsValues [][]interface{}
	)

	columns = []string{}
	for field, value := range i.FieldsValues {
		columns = append(columns, field)
		if rowCount < len(value) {
			rowCount = len(value)
		}
	}

	sort.Slice(columns, func(i, j int) bool {
		return columns[i] < columns[j]
	})

	rowsValues = [][]interface{}{}
	for rowIndex := 0; rowIndex < rowCount; rowIndex++ {
		var rowValues []interface{} = []interface{}{}

		for columnIndex := 0; columnIndex < len(columns); columnIndex++ {
			if rowIndex >= len(i.FieldsValues[columns[columnIndex]]) {
				continue
			}

			rowValues = append(rowValues, i.FieldsValues[columns[columnIndex]][rowIndex])
		}

		rowsValues = append(rowsValues, rowValues)
	}

	return columns, rowsValues
}

func (i *InsertQuery) validate() error {
	var (
		columns    []string
		rowsValues [][]interface{}
	)

	if i.Table == "" {
		return errors.New("table is required")
	}

	columns, rowsValues = i.getColumnsAndRowsValues()

	if len(columns) == 0 {
		return errors.New("fields is required")
	}

	for columnIndex := 0; columnIndex < len(columns); columnIndex++ {
		if columns[columnIndex] == "" {
			return errors.New("field is required")
		}
	}

	if len(rowsValues) == 0 {
		return errors.New("values is required")
	}

	for rowIndex := 0; rowIndex < len(rowsValues); rowIndex++ {
		var (
			rowValues    []interface{}
			reflectValue reflect.Value
		)

		rowValues = rowsValues[rowIndex]

		if len(rowValues) != len(columns) {
			return errors.New("value length is not equal to fields length")
		}

		for columnIndex := 0; columnIndex < len(rowValues); columnIndex++ {
			if rowValues[columnIndex] != nil {
				reflectValue = reflect.ValueOf(rowValues[columnIndex])

				if !allowedKindValue[reflectValue.Kind()] || reflectValue.Kind() == reflect.Array || reflectValue.Kind() == reflect.Slice {
					return fmt.Errorf("unsupported %s value type", reflectValue.Kind().String())
				}
			}
		}
	}

	return nil
}

func (i *InsertQuery) ToSQLWithArgs(dialect Dialect) (string, []interface{}, error) {
	var (
		columns      []string
		rowsValues   [][]interface{}
		query        string
		args         []interface{}
		placeholders []string
		err          error
	)

	err = i.validate()
	if err != nil {
		return "", nil, err
	}

	if dialect == "" {
		return "", nil, errors.New("dialect is required")
	}

	columns, rowsValues = i.getColumnsAndRowsValues()
	args = []interface{}{}

	for rowIndex := 0; rowIndex < len(rowsValues); rowIndex++ {
		var (
			placeholderStartIdx int
			placeholderEndIdx   int
			placeholder         string
		)

		args = append(args, rowsValues[rowIndex]...)
		placeholderStartIdx = len(args) - (len(rowsValues[rowIndex]) - 1)
		placeholderEndIdx = len(args)
		placeholder = fmt.Sprintf("(%s)", getPlaceholder(dialect, placeholderStartIdx, placeholderEndIdx))
		placeholders = append(placeholders, placeholder)
	}

	query = fmt.Sprintf("insert into %s(%s) values %s", i.Table, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	return query, args, nil
}
