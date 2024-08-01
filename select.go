package ququery

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type SelectQuery struct {
	table   string
	columns []string
	*WhereContainer
	joins     [][]string
	orderBy   []string
	hasLimit  bool
	hasOffset bool
}

func Select(table string) SelectQuery {
	return SelectQuery{
		table:          table,
		WhereContainer: &WhereContainer{},
	}
}

func (q SelectQuery) Table(table string) SelectQuery {
	q.table = table

	if q.WhereContainer == nil {
		q.WhereContainer = &WhereContainer{}
	}

	return q
}

func (q SelectQuery) Columns(columns ...string) SelectQuery {
	q.columns = columns

	return q
}

func (q SelectQuery) Where(condition ...string) SelectQuery {
	q.WhereContainer.Where(condition...)

	return q
}

func (q SelectQuery) OrWhere(condition ...string) SelectQuery {
	q.WhereContainer.OrWhere(condition...)

	return q
}

// WhereNotNull get rows where column values is not null
func (q SelectQuery) WhereNotNull(column string) SelectQuery {
	q.conditions = append(q.conditions, whereStructure{
		column:   "",
		operator: "",
		isAnd:    true,
		isRaw:    true,
		rawQuery: column + " IS NOT NULL",
	})

	return q
}

// OrWhereNull get rows where column values is not null or previous condition is true
func (q SelectQuery) OrWhereNotNull(column string) SelectQuery {
	q.conditions = append(q.conditions, whereStructure{
		column:   "",
		operator: "",
		isAnd:    false,
		isRaw:    true,
		rawQuery: column + " IS NOT NULL",
	})

	return q
}

// WhereNull get rows where column value is null
func (q SelectQuery) WhereNull(column string) SelectQuery {
	q.conditions = append(q.conditions, whereStructure{
		column:   "",
		operator: "",
		isAnd:    true,
		isRaw:    true,
		rawQuery: column + " IS NULL",
	})

	return q
}

// OrWhereNull get rows where column values is null or previous condition is true
func (q SelectQuery) OrWhereNull(column string) SelectQuery {
	q.conditions = append(q.conditions, whereStructure{
		column:   "",
		operator: "",
		isAnd:    false,
		isRaw:    true,
		rawQuery: column + " IS NULL",
	})

	return q
}

func (q SelectQuery) StrPos(column string) SelectQuery {
	q.conditions = append(q.conditions, whereStructure{
		column:   "",
		operator: "",
		isAnd:    true,
		isRaw:    true,
		rawQuery: fmt.Sprintf("(STRPOS(%s, ?) > 0 OR ? = '')", column),
	})

	return q
}

func (q SelectQuery) OrStrPos(column string) SelectQuery {
	q.conditions = append(q.conditions, whereStructure{
		column:   "",
		operator: "",
		rawQuery: fmt.Sprintf("(STRPOS(%s, ?) > 0 OR ? = '')", column),
		isAnd:    false,
		isRaw:    true,
	})

	return q
}

func (q SelectQuery) MultiWhere(f func(subQuery MultiWhere) string) SelectQuery {
	query := f(MultiWhere{})

	q.conditions = append(q.conditions, whereStructure{
		column:   "",
		operator: "",
		isAnd:    true,
		isRaw:    true,
		rawQuery: query,
	})

	return q
}

func (q SelectQuery) WhereInSubquery(column string, subQuery func(q SelectQuery) string) SelectQuery {
	q.conditions = append(q.conditions, whereStructure{
		column:   "",
		operator: "",
		rawQuery: fmt.Sprintf("%s IN (%s)", column, subQuery(SelectQuery{})),
		isAnd:    true,
		isRaw:    true,
	})

	return q
}

func (q SelectQuery) Join(table string, conditions ...string) SelectQuery {
	q.joins = append(q.joins, append([]string{table}, conditions...))

	return q
}

// With can load one-to-many relations without need to pass join column
func (q SelectQuery) With(entities ...string) SelectQuery {
	for _, entity := range entities {
		table := findTableFromEntity(entity)

		q.joins = append(q.joins, []string{
			table,
			fmt.Sprintf("%s.id = %s.%s", table, q.table, entity+"_id"),
		})
	}

	return q
}

func (q SelectQuery) OrderBy(column, direction string) SelectQuery {
	q.orderBy = []string{column, direction}

	return q
}

func (q SelectQuery) Limit() SelectQuery {
	q.hasLimit = true

	return q
}

func (q SelectQuery) Offset() SelectQuery {
	q.hasOffset = true

	return q
}

func (q SelectQuery) prepareSelectQuery() string {
	columns := strings.Join(q.columns, ", ")
	if len(q.columns) == 0 {
		columns = "*"
	}

	query := fmt.Sprintf("SELECT %s FROM %s", columns, q.table)

	if len(q.joins) > 0 {
		query += " " + q.prepareJoinQuery(q.joins)
	}

	if len(q.WhereContainer.conditions) > 0 {
		query += " " + prepareWhereQuery(q.WhereContainer.conditions)
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

func (q SelectQuery) Query() string {
	query := q.prepareSelectQuery()

	return sqlx.Rebind(sqlx.DOLLAR, query)
}

func (q SelectQuery) QueryWithoutRebinding() string {
	return q.prepareSelectQuery()
}

func (q SelectQuery) prepareJoinQuery(joins [][]string) string {
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
