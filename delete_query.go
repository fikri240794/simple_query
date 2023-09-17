package simple_query

import (
	"errors"
	"fmt"
)

type DeleteQuery struct {
	Table  string
	Filter *Filter
}

func Delete() *DeleteQuery {
	return &DeleteQuery{}
}

func (d *DeleteQuery) From(table string) *DeleteQuery {
	d.Table = table
	return d
}

func (d *DeleteQuery) Where(filter *Filter) *DeleteQuery {
	d.Filter = filter
	return d
}

func (d *DeleteQuery) validate() error {
	if d.Table == "" {
		return errors.New("table is required")
	}

	if d.Filter == nil {
		return errors.New("filter is required")
	}

	return nil
}

func (d *DeleteQuery) ToSQLWithArgs(dialect Dialect) (string, []interface{}, error) {
	var (
		query       string
		args        []interface{}
		whereClause string
		err         error
	)

	err = d.validate()
	if err != nil {
		return "", nil, err
	}

	query = fmt.Sprintf("delete from %s", d.Table)
	args = []interface{}{}

	if d.Filter != nil {
		whereClause, args, err = d.Filter.ToSQLWithArgs(dialect, args)
		if err != nil {
			return "", nil, err
		}

		if whereClause != "" {
			query = fmt.Sprintf("%s where %s", query, whereClause)
		}
	}

	return query, args, nil
}
