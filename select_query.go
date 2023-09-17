package simple_query

import (
	"errors"
	"fmt"
	"strings"
)

type SelectQuery struct {
	Fields  []string
	Table   string
	Filter  *Filter
	Sorts   []*Sort
	maxTake uint64
	Take    uint64
}

func Select(fields ...string) *SelectQuery {
	return &SelectQuery{
		Fields:  fields,
		maxTake: 100,
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
		return errors.New("fields is required")
	}

	for i := 0; i < len(s.Fields); i++ {
		if s.Fields[i] == "" {
			return errors.New("field is required")
		}
	}

	if s.Table == "" {
		return errors.New("table is required")
	}

	if s.Take == 0 {
		return errors.New("take is required")
	}

	if s.Take > s.maxTake {
		return fmt.Errorf("maximum take is %d", s.maxTake)
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
		for i := 0; i < len(s.Sorts); i++ {
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

	args = append(args, s.Take)
	placeholder = getPlaceholder(dialect, len(args), len(args))
	query = fmt.Sprintf("%s limit %s", query, placeholder)

	return query, args, nil
}
