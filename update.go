package ququery

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UpdateQuery struct {
	table   string
	columns []string
	WhereContainer[*UpdateQuery]
}

func Update(table string) *UpdateQuery {
	u := &UpdateQuery{table: table}
	u.WhereContainer = WhereContainer[*UpdateQuery]{self: u}

	return u
}

func (q *UpdateQuery) Set(columns ...string) *UpdateQuery {
	q.columns = append(q.columns, columns...)

	return q
}

func (q *UpdateQuery) Query() string {
	query := fmt.Sprintf(
		`
			UPDATE %s 
			SET %s
			%s
		`,
		q.table,
		prepareUpdateQuery(q.columns),
		prepareWhereQuery(q.conditions),
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
