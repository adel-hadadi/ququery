package ququery

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type (
	SelectQuery struct {
		table   string
		columns []string
		WhereContainer[*SelectQuery]
		joins            []join
		orderBy          []string
		hasLimit         bool
		hasOffset        bool
		withoutRebinding bool
	}

	joinType string

	join struct {
		table       string
		constraints string
		jType       joinType
	}
)

const (
	rightJoin joinType = "RIGHT"
	leftJoin           = "LEFT"
	innerJoin          = "INNER"
)

func Select(table string) *SelectQuery {
	s := &SelectQuery{table: table}
	s.WhereContainer = WhereContainer[*SelectQuery]{self: s}

	return s
}

func (q *SelectQuery) Table(table string) *SelectQuery {
	q.table = table

	q.WhereContainer = WhereContainer[*SelectQuery]{self: q}

	return q
}

func (q *SelectQuery) Columns(columns ...string) *SelectQuery {
	q.columns = columns

	return q
}

// Join method used to add inner join to your queries
//
// Example:
//
//	query := ququery.Select("users").Join("posts", "posts.user_id = users.id").Query()
//	log.Println(query) => SELECT * FROM users INNER JOIN posts ON posts.user_id = users.id
func (q *SelectQuery) Join(table, constraints string) *SelectQuery {
	q.joins = append(q.joins, join{
		table:       table,
		constraints: constraints,
		jType:       innerJoin,
	})

	return q
}

// LeftJoin method used to add left join to your queries
//
// Example:
//
//	query := ququery.Select("users").LeftJoin("posts", "posts.user_id = users.id").Query()
//	log.Println(query) => SELECT * FROM users LEFT JOIN posts ON posts.user_id = users.id
func (q *SelectQuery) LeftJoin(table, constraints string) *SelectQuery {
	q.joins = append(q.joins, join{
		table:       table,
		constraints: constraints,
		jType:       leftJoin,
	})

	return q
}

// RightJoin method used to add right join to your queries
//
// Example:
//
//	query := ququery.Select("users").RightJoin("posts", "posts.user_id = users.id").Query()
//	log.Println(query) => SELECT * FROM users RIGHT JOIN posts ON posts.user_id = users.id
func (q *SelectQuery) RightJoin(table, constraints string) *SelectQuery {
	q.joins = append(q.joins, join{
		table:       table,
		constraints: constraints,
		jType:       rightJoin,
	})

	return q
}

// With can load one-to-many relations without need to pass join column
func (q *SelectQuery) With(entities ...string) *SelectQuery {
	for _, entity := range entities {
		table := findTableFromEntity(entity)

		q.joins = append(q.joins, join{
			table:       table,
			constraints: fmt.Sprintf("%s.id = %s.%s", table, q.table, entity+"_id"),
			jType:       leftJoin,
		})
	}

	return q
}

func (q *SelectQuery) OrderBy(column, direction string) *SelectQuery {
	q.orderBy = []string{column, direction}

	return q
}

func (q *SelectQuery) Limit() *SelectQuery {
	q.hasLimit = true

	return q
}

func (q *SelectQuery) Offset() *SelectQuery {
	q.hasOffset = true

	return q
}

func (q *SelectQuery) prepareSelectQuery() string {
	columns := strings.Join(q.columns, ", ")
	if len(q.columns) == 0 {
		columns = "*"
	}

	query := fmt.Sprintf("SELECT %s FROM %s", columns, q.table)

	if len(q.joins) > 0 {
		query += " " + q.prepareJoinQuery(q.joins)
	}

	if len(q.conditions) > 0 {
		query += " " + prepareWhereQuery(q.conditions)
	}

	if len(q.orderBy) > 0 {
		query += fmt.Sprintf(" ORDER BY %s %s", q.orderBy[0], strings.ToUpper(q.orderBy[1]))
	}

	if q.hasLimit {
		query += " LIMIT ?"
	}

	if q.hasOffset {
		query += " OFFSET ?"
	}

	return strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(query, "\n", ""), "\t", ""))
}

func (q *SelectQuery) Query() string {
	query := q.prepareSelectQuery()

	if q.withoutRebinding {
		return query
	}

	return sqlx.Rebind(sqlx.DOLLAR, query)
}

func (q *SelectQuery) prepareJoinQuery(joins []join) string {
	var joinQuery string

	for _, join := range joins {
		joinQuery += fmt.Sprintf(
			" %s JOIN %s ON %s",
			join.jType,
			join.table,
			join.constraints,
		)
	}

	return joinQuery
}
