package ququery

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ExistsQuery struct {
	table string
	WhereContainer[*ExistsQuery]
}

func Exists(table string) *ExistsQuery {
	e := &ExistsQuery{table: table}
	e.WhereContainer = WhereContainer[*ExistsQuery]{self: e}

	return e
}

func (q *ExistsQuery) Query() string {
	query := fmt.Sprintf(
		"SELECT EXISTS(SELECT true FROM %s %s)",
		q.table,
		prepareWhereQuery(q.conditions),
	)

	return sqlx.Rebind(sqlx.DOLLAR, query)
}
