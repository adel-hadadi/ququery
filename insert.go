package ququery

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type InsertQuery struct {
	table      string
	columns    []string
	returnings []string
}

func Insert(table string) InsertQuery {
	return InsertQuery{
		table: table,
	}
}

func (q InsertQuery) Into(columns ...string) InsertQuery {
	q.columns = columns

	return q
}

func (q InsertQuery) Returning(columns ...string) InsertQuery {
	q.returnings = columns

	return q
}

func (q InsertQuery) Query() string {
	query := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES (%s)`,
		q.table,
		strings.Join(q.columns, ", "),
		prepareInsertQuery(q.columns),
	)

	if len(q.returnings) > 0 {
		query += fmt.Sprintf(" RETURNING (%s)", strings.Join(q.returnings, ", "))
	}

	return sqlx.Rebind(sqlx.DOLLAR, query)
}

func prepareInsertQuery(columns []string) string {
	var query string

	for i := 0; i < len(columns); i++ {
		if i != len(columns)-1 {
			query += "?,"
			continue
		}

		query += "?"
	}

	return query
}
