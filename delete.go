package ququery

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DeleteQuery struct {
	table string
	WhereContainer[*DeleteQuery]
}

func Delete(table string) *DeleteQuery {
	q := &DeleteQuery{table: table}
	q.WhereContainer = WhereContainer[*DeleteQuery]{self: q}

	return q
}

func (q *DeleteQuery) Query() string {
	query := fmt.Sprintf(
		`DELETE FROM %s %s`,
		q.table,
		prepareWhereQuery(q.conditions),
	)

	return sqlx.Rebind(sqlx.DOLLAR, query)
}
