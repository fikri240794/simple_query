# Simple Query
Simple dynamic SQL query builder. Intended for building simple query with basic logic. Currently, supported dialect is MySQL and Postgres only.

## Installation
```bash
go get github.com/fikri240794/simple_query
```

## Usage
### Example for SELECT:
```go
package main

import (
	"log"
	"github.com/fikri240794/simple_query"
)

func main() {
	var (
		selectQuery *simple_query.SelectQuery
		query       string
		args        []interface{}
		err         error
	)

	selectQuery = simple_query.Select("field1", "field2", "field3", "field4", "field5").
		From("table1").
		Where(
			simple_query.NewFilter().
				SetLogic(simple_query.FilterLogicAnd).
				AddFilters(
					simple_query.NewFilter().
						SetCondition("field1", simple_query.FilterOperatorEqual, "value1"),
					simple_query.NewFilter().
						SetCondition("field2", simple_query.FilterOperatorNotEqual, true),
					simple_query.NewFilter().
						SetLogic(simple_query.FilterLogicOr).
						AddFilter("field3", simple_query.FilterOperatorGreaterThan, 50).
						AddFilter("field4", simple_query.FilterOperatorGreaterThanOrEqual, 75.4),
					simple_query.NewFilter().
						SetLogic(simple_query.FilterLogicOr).
						AddFilter("field5", simple_query.FilterOperatorLessThan, "value5").
						AddFilter("field6", simple_query.FilterOperatorLessThanOrEqual, "value6"),
					simple_query.NewFilter().
						SetLogic(simple_query.FilterLogicAnd).
						AddFilter("field7", simple_query.FilterOperatorIsNull, nil).
						AddFilter("field8", simple_query.FilterOperatorIsNotNull, nil),
					simple_query.NewFilter().
						SetLogic(simple_query.FilterLogicOr).
						AddFilters(
							simple_query.NewFilter().
								SetLogic(simple_query.FilterLogicAnd).
								AddFilter("field9", simple_query.FilterOperatorIn, []string{"value9.1", "value9.2", "value9.3"}).
								AddFilter("field10", simple_query.FilterOperatorNotIn, [3]float64{10.1, 10.2, 10.3}),
							simple_query.NewFilter().SetCondition("field11", simple_query.FilterOperatorLike, "value11"),
							simple_query.NewFilter().SetCondition("field12", simple_query.FilterOperatorNotLike, "value12"),
						),
				),
		).
		OrderBy(
			simple_query.NewSort("field1", simple_query.SortDirectionAscending),
			simple_query.NewSort("field2", simple_query.SortDirectionDescending),
		).
		Limit(50)

	query, args, err = selectQuery.ToSQLWithArgs(simple_query.DialectPostgres)

	log.Printf("query: %s", query)
	/*
		-- QUERY --
		select
			field1,
			field2,
			field3,
			field4,
			field5
		from
			table1
		where
			field1 = $1
			and field2 != $2
			and (
				field3 > $3
				or field4 >= $4
			)
			and (
				field5 < $5
				or field6 <= $6
			)
			and (
				field7 is null
				and field8 is not null
			)
			and (
				(
					field9 in ($7, $8, $9)
					and field10 not in ($10, $11, $12)
				)
				or field11 ilike concat('%', $13, '%')
				or field12 not ilike concat('%', $14, '%')
			)
		order by
			field1 asc,
			field2 desc
		limit
			$15
	*/

	log.Printf("args: %v", args)
	/*
		-- ARGS --
		[
			"value1",
			true,
			50,
			75.4,
			"value5",
			"value6",
			"value9.1",
			"value9.2",
			"value9.3",
			10.1,
			10.2,
			10.3,
			"value11",
			"value12",
			50
		]
	*/

	log.Printf("err: %v", err) // nil
}
```

### Example for INSERT:
```go
package main

import (
	"log"
	"github.com/fikri240794/simple_query"
)

func main() {
	var (
		insertQuery *simple_query.InsertQuery
		query       string
		args        []interface{}
		err         error
	)

	insertQuery = simple_query.Insert().Into("table1").
		Value("field1", 1).
		Value("field2", "value2.1").
		Value("field3", 3.14).
		Value("field4", 4).
		Value("field5", false).
		Value("field1", 2).
		Value("field2", "value2.2").
		Value("field3", 3.14).
		Value("field4", 4).
		Value("field5", false).
		Value("field1", 3).
		Value("field2", "value2.1").
		Value("field3", 3.14).
		Value("field4", 4).
		Value("field5", false)

	query, args, err = insertQuery.ToSQLWithArgs(simple_query.DialectPostgres)

	log.Printf("query: %s", query)
	/*
		-- QUERY --
		insert into
			table1(field1, field2, field3, field4, field5)
		values
			($1, $2, $3, $4, $5),
			($6, $7, $8, $9, $10),
			($11, $12, $13, $14, $15)
	*/

	log.Printf("args: %v", args)
	/*
		-- ARGS --
		[
			1,
			"value2.1",
			3.14,
			4,
			false,
			2,
			"value2.2",
			3.14,
			4,
			false,
			3,
			"value2.1",
			3.14,
			4,
			false
		]
	*/

	log.Printf("err: %v", err) // nil
}
```

### Example for UPDATE:
```go
package main

import (
	"log"
	"github.com/fikri240794/simple_query"
)

func main() {
	var (
		updateQuery *simple_query.UpdateQuery
		query       string
		args        []interface{}
		err         error
	)

	updateQuery = simple_query.Update("table1").
		Set("field2", 1).
		Set("field3", "value3").
		Set("field4", 4.265).
		Set("field5", true).
		Where(
			simple_query.NewFilter().
				SetLogic(simple_query.FilterLogicAnd).
				AddFilter("field1", simple_query.FilterOperatorEqual, "value1"),
		)

	query, args, err = updateQuery.ToSQLWithArgs(simple_query.DialectPostgres)

	log.Printf("query: %s", query)
	/*
		-- QUERY --
		update
			table1
		set
			field2 = $1,
			field3 = $2,
			field4 = $3,
			field5 = $4
		where
			field1 = $5
	*/
	log.Printf("args: %v", args)
	/*
		-- ARGS --
		[
			1,
			"value3",
			4.265,
			true,
			"value1"
		]
	*/
	log.Printf("err: %v", err) // nil
}
```

### Example for DELETE:
```go
package main

import (
	"log"
	"github.com/fikri240794/simple_query"
)

func main() {
	var (
		deleteQuery *simple_query.DeleteQuery
		query       string
		args        []interface{}
		err         error
	)

	deleteQuery = simple_query.Delete().
		From("table1").
		Where(
			simple_query.NewFilter().
				SetLogic(simple_query.FilterLogicAnd).
				AddFilter("field1", simple_query.FilterOperatorEqual, "value1"),
		)

	query, args, err = deleteQuery.ToSQLWithArgs(simple_query.DialectPostgres)

	log.Printf("query: %s", query)
	/*
		-- QUERY --
		delete from
			table1
		where
			field1 = $1
	*/

	log.Printf("args: %v", args)
	/*
		-- ARGS --
		["value1"]
	*/

	log.Printf("err: %v", err) // nil
}
```