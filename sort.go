package simple_query

import (
	"errors"
	"fmt"
)

type SortDirection string

const (
	SortDirectionAscending  SortDirection = "asc"
	SortDirectionDescending SortDirection = "desc"
)

type Sort struct {
	Field     string
	Direction SortDirection
}

func NewSort(field string, direction SortDirection) *Sort {
	return &Sort{
		Field:     field,
		Direction: direction,
	}
}

func (s *Sort) validate() error {
	if s.Field == "" {
		return errors.New("field is required")
	}

	return nil
}

func (s *Sort) ToSQL() (string, error) {
	var (
		orderByQueryFormat string
		orderByQuery       string
		err                error
	)

	err = s.validate()
	if err != nil {
		return "", err
	}

	if s.Direction == "" {
		s.Direction = SortDirectionAscending
	}

	orderByQueryFormat = "%s %s"
	orderByQuery = fmt.Sprintf(orderByQueryFormat, s.Field, s.Direction)

	return orderByQuery, nil
}
