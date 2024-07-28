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
	sq "github.com/fikri240794/simple_query"
)

func main() {
	var (
		selectQuery *sq.SelectQuery
		query       string
		args        []interface{}
		err         error
	)

	selectQuery = sq.Select(
		sq.NewField("field1"),
		sq.NewField("field2"),
		sq.NewField("field3"),
		sq.NewField("field4"),
		sq.NewField("field5"),
	).
		From(sq.NewTable("table1")).
		Where(
			sq.NewFilter().
				SetLogic(sq.LogicAnd).
				AddFilters(
					sq.NewFilter().
						SetCondition(
							sq.NewField("field1"),
							sq.OperatorEqual,
							"value1",
						),
					sq.NewFilter().
						SetCondition(
							sq.NewField("field2"),
							sq.OperatorNotEqual,
							true,
						),
					sq.NewFilter().
						SetLogic(sq.LogicOr).
						AddFilter(
							sq.NewField("field3"),
							sq.OperatorGreaterThan,
							50,
						).
						AddFilter(
							sq.NewField("field4"),
							sq.OperatorGreaterThanOrEqual,
							75.4,
						),
					sq.NewFilter().
						SetLogic(sq.LogicOr).
						AddFilter(
							sq.NewField("field5"),
							sq.OperatorLessThan,
							"value5",
						).
						AddFilter(
							sq.NewField("field6"),
							sq.OperatorLessThanOrEqual,
							"value6",
						),
					sq.NewFilter().
						SetLogic(sq.LogicAnd).
						AddFilter(
							sq.NewField("field7"),
							sq.OperatorIsNull,
							nil,
						).
						AddFilter(
							sq.NewField("field8"),
							sq.OperatorIsNotNull,
							nil,
						),
					sq.NewFilter().
						SetLogic(sq.LogicOr).
						AddFilters(
							sq.NewFilter().
								SetLogic(sq.LogicAnd).
								AddFilter(
									sq.NewField("field9"),
									sq.OperatorIn,
									[]string{
										"value9.1",
										"value9.2",
										"value9.3",
									},
								).
								AddFilter(
									sq.NewField("field10"),
									sq.OperatorNotIn,
									[3]float64{
										10.1,
										10.2,
										10.3,
									},
								),
							sq.NewFilter().
								SetCondition(
									sq.NewField("field11"),
									sq.OperatorLike,
									"value11",
								),
							sq.NewFilter().
								SetCondition(
									sq.NewField("field12"),
									sq.OperatorNotLike,
									"value12",
								),
						),
				),
		).
		OrderBy(
			sq.NewSort(
				"field1",
				sq.SortDirectionAscending,
			),
			sq.NewSort(
				"field2",
				sq.SortDirectionDescending,
			),
		).
		Limit(50)

	query, args, err = selectQuery.ToSQLWithArgsWithAlias(sq.DialectPostgres, []interface{}{})

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
	sq "github.com/fikri240794/simple_query"
)

func main() {
	var (
		insertQuery *sq.InsertQuery
		query       string
		args        []interface{}
		err         error
	)

	insertQuery = sq.Insert().
		Into("table1").
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

	query, args, err = insertQuery.ToSQLWithArgs(sq.DialectPostgres)

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
	sq "github.com/fikri240794/simple_query"
)

func main() {
	var (
		updateQuery *sq.UpdateQuery
		query       string
		args        []interface{}
		err         error
	)

	updateQuery = sq.Update("table1").
		Set("field2", 1).
		Set("field3", "value3").
		Set("field4", 4.265).
		Set("field5", true).
		Where(
			sq.NewFilter().
				SetLogic(sq.LogicAnd).
				AddFilter(
					sq.NewField("field1"),
					sq.OperatorEqual,
					"value1",
				),
		)

	query, args, err = updateQuery.ToSQLWithArgs(sq.DialectPostgres)

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
	sq "github.com/fikri240794/simple_query"
)

func main() {
	var (
		deleteQuery *sq.DeleteQuery
		query       string
		args        []interface{}
		err         error
	)

	deleteQuery = sq.Delete().
		From("table1").
		Where(
			sq.NewFilter().
				SetLogic(sq.LogicAnd).
				AddFilter(
					sq.NewField("field1"),
					sq.OperatorEqual,
					"value1",
				),
		)

	query, args, err = deleteQuery.ToSQLWithArgs(sq.DialectPostgres)

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