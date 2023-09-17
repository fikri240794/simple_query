package simple_query

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type UpdateQuery struct {
	Table       string
	FieldsValue map[string]interface{}
	Filter      *Filter
}

func Update(table string) *UpdateQuery {
	return &UpdateQuery{
		Table:       table,
		FieldsValue: map[string]interface{}{},
	}
}

func (u *UpdateQuery) Set(field string, value interface{}) *UpdateQuery {
	u.FieldsValue[field] = value
	return u
}

func (u *UpdateQuery) Where(filter *Filter) *UpdateQuery {
	u.Filter = filter
	return u
}

func (u *UpdateQuery) validate() error {
	if u.Table == "" {
		return errors.New("table is required")
	}

	if len(u.FieldsValue) == 0 {
		return errors.New("fields is required")
	}

	for field, value := range u.FieldsValue {
		if field == "" {
			return errors.New("field is required")
		}

		if value != nil {
			var reflectValue reflect.Value = reflect.ValueOf(value)

			if !allowedKindValue[reflectValue.Kind()] || reflectValue.Kind() == reflect.Array || reflectValue.Kind() == reflect.Slice {
				return fmt.Errorf("unsupported %s value type", reflectValue.Kind().String())
			}
		}
	}

	if u.Filter == nil {
		return errors.New("filter is required")
	}

	return nil
}

func (u *UpdateQuery) ToSQLWithArgs(dialect Dialect) (string, []interface{}, error) {
	var (
		query        string
		args         []interface{}
		placeholders []string
		whereClause  string
		err          error
	)

	err = u.validate()
	if err != nil {
		return "", nil, err
	}

	query = fmt.Sprintf("update %s", u.Table)
	placeholders = []string{}

	for field, value := range u.FieldsValue {
		var (
			placeholderStartIdx int
			placeholderEndIdx   int
			placeholder         string
		)

		args = append(args, value)
		placeholderStartIdx = len(args)
		placeholderEndIdx = len(args)
		placeholder = fmt.Sprintf("%s = %s", field, getPlaceholder(dialect, placeholderStartIdx, placeholderEndIdx))
		placeholders = append(placeholders, placeholder)
	}

	query = fmt.Sprintf("%s set %s", query, strings.Join(placeholders, ", "))

	if u.Filter != nil {
		whereClause, args, err = u.Filter.ToSQLWithArgs(dialect, args)
		if err != nil {
			return "", nil, err
		}

		if whereClause != "" {
			query = fmt.Sprintf("%s where %s", query, whereClause)
		}
	}

	return query, args, nil
}
