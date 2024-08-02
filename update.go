package ququery

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UpdateQuery struct {
	table          string
	columns        []string
	whereContainer *whereContainer
}

func Update(table string) UpdateQuery {
	return UpdateQuery{
		table:          table,
		whereContainer: &whereContainer{},
	}
}

func (q UpdateQuery) Set(columns ...string) UpdateQuery {
	q.columns = append(q.columns, columns...)

	return q
}

// Where get rows where column
func (q UpdateQuery) Where(condition ...string) UpdateQuery {
	q.whereContainer.where(condition...)

	return q
}

func (q UpdateQuery) OrWhere(condition ...string) UpdateQuery {
	q.whereContainer.orWhere(condition...)

	return q
}

// WhereNotNull get rows where column values is not null
func (q UpdateQuery) WhereNotNull(column string) UpdateQuery {
	q.whereContainer.whereNotNull(column)

	return q
}

// OrWhereNull get rows where column values is not null or previous condition is true
func (q UpdateQuery) OrWhereNotNull(column string) UpdateQuery {
	q.whereContainer.orWhereNotNull(column)

	return q
}

// WhereNull get rows where column value is null
func (q UpdateQuery) WhereNull(column string) UpdateQuery {
	q.whereContainer.whereNull(column)

	return q
}

// OrWhereNull get rows where column values is null or previous condition is true
func (q UpdateQuery) OrWhereNull(column string) UpdateQuery {
	q.whereContainer.orWhereNull(column)

	return q
}

func (q UpdateQuery) StrPos(column string) UpdateQuery {
	q.whereContainer.strpos(column)

	return q
}

func (q UpdateQuery) OrStrPos(column string) UpdateQuery {
	q.whereContainer.orStrpos(column)

	return q
}

func (q UpdateQuery) MultiWhere(f func(subQuery MultiWhere) string) UpdateQuery {
	q.whereContainer.multiWhere(f)

	return q
}

func (q UpdateQuery) WhereInSubquery(column string, subQuery func(q SelectQuery) string) UpdateQuery {
	q.whereContainer.whereInSubquery(column, subQuery)

	return q
}

func (q UpdateQuery) Query() string {
	query := fmt.Sprintf(
		`
			UPDATE %s 
			SET %s
			%s
		`,
		q.table,
		prepareUpdateQuery(q.columns),
		prepareWhereQuery(q.whereContainer.conditions),
	)

	return sqlx.Rebind(sqlx.DOLLAR, query)
}

func prepareUpdateQuery(columns []string) string {
	var query string
	for i, column := range columns {
		query += fmt.Sprintf(" %s = ?", column)
		if i != len(columns)-1 {
			query += ","
		}
	}

	return query
}
