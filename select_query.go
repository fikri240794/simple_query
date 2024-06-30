package simple_query

import (
	"fmt"
	"strings"
)

type SelectQuery struct {
	Fields []string
	Table  string
	Filter *Filter
	Sorts  []*Sort
	Take   uint64
}

func Select(fields ...string) *SelectQuery {
	return &SelectQuery{
		Fields: fields,
	}
}

func (s *SelectQuery) From(table string) *SelectQuery {
	s.Table = table
	return s
}

func (s *SelectQuery) Where(filter *Filter) *SelectQuery {
	s.Filter = filter
	return s
}

func (s *SelectQuery) OrderBy(sorts ...*Sort) *SelectQuery {
	s.Sorts = sorts
	return s
}

func (s *SelectQuery) Limit(take uint64) *SelectQuery {
	s.Take = take
	return s
}

func (s *SelectQuery) validate() error {
	if len(s.Fields) == 0 {
		return ErrFieldsIsRequired
	}

	for i := range s.Fields {
		if s.Fields[i] == "" {
			return ErrFieldIsRequired
		}
	}

	if s.Table == "" {
		return ErrTableIsRequired
	}

	return nil
}

func (s *SelectQuery) ToSQLWithArgs(dialect Dialect) (string, []interface{}, error) {
	var (
		query         string
		whereClause   string
		orderBy       string
		orderByClause []string
		args          []interface{}
		placeholder   string
		err           error
	)

	err = s.validate()
	if err != nil {
		return "", nil, err
	}

	query = fmt.Sprintf("select %s from %s", strings.Join(s.Fields, ", "), s.Table)
	args = []interface{}{}

	if s.Filter != nil {
		whereClause, args, err = s.Filter.ToSQLWithArgs(dialect, args)
		if err != nil {
			return "", nil, err
		}

		if whereClause != "" {
			query = fmt.Sprintf("%s where %s", query, whereClause)
		}
	}

	if len(s.Sorts) > 0 {
		orderByClause = []string{}
		for i := range s.Sorts {
			if s.Sorts[i] == nil {
				continue
			}

			orderBy, err = s.Sorts[i].ToSQL()
			if err != nil {
				return "", nil, err
			}

			orderByClause = append(orderByClause, orderBy)
		}

		if len(orderByClause) > 0 {
			query = fmt.Sprintf("%s order by %s", query, strings.Join(orderByClause, ", "))
		}
	}

	if s.Take > 0 {
		args = append(args, s.Take)
		placeholder = getPlaceholder(dialect, len(args), len(args))
		query = fmt.Sprintf("%s limit %s", query, placeholder)
	}

	return query, args, nil
}
