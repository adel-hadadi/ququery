package ququery

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type SelectQuery struct {
	table   string
	columns []string
	WhereContainer[*SelectQuery]
	joins            [][]string
	orderBy          []string
	hasLimit         bool
	hasOffset        bool
	withoutRebinding bool
}

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

func (q *SelectQuery) Join(table string, conditions ...string) *SelectQuery {
	q.joins = append(q.joins, append([]string{table}, conditions...))

	return q
}

// With can load one-to-many relations without need to pass join column
func (q *SelectQuery) With(entities ...string) *SelectQuery {
	for _, entity := range entities {
		table := findTableFromEntity(entity)

		q.joins = append(q.joins, []string{
			table,
			fmt.Sprintf("%s.id = %s.%s", table, q.table, entity+"_id"),
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

func (q *SelectQuery) prepareJoinQuery(joins [][]string) string {
	var joinQuery string

	for _, join := range joins {
		joinQuery += fmt.Sprintf(
			" LEFT JOIN %s ON %s",
			join[0],
			strings.Join(join[1:], " AND "),
		)
	}

	return joinQuery
}
