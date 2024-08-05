package simple_query

import "fmt"

type FilterValue struct {
	Value       interface{}
	SelectQuery *SelectQuery
}

func NewFilterValue(value interface{}) *FilterValue {
	return &FilterValue{
		Value: value,
	}
}

func NewSelectQueryFilterValue(selectQuery *SelectQuery) *FilterValue {
	return &FilterValue{
		SelectQuery: selectQuery,
	}
}

func (v *FilterValue) validate(dialect Dialect) error {
	if dialect == "" {
		return ErrDialectIsRequired
	}

	return nil
}

func (v *FilterValue) ToSQLWithArgs(dialect Dialect, args []interface{}) (string, []interface{}, error) {
	var (
		query string
		err   error
	)

	err = v.validate(dialect)
	if err != nil {
		return "", nil, err
	}

	if v.SelectQuery == nil {
		args = append(args, v.Value)

		return "", args, nil
	}

	query, args, err = v.SelectQuery.ToSQLWithArgsWithAlias(dialect, args)
	if err != nil {
		return "", nil, err
	}

	query = fmt.Sprintf("(%s)", query)

	return query, args, nil
}
