package ququery

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ExistsQuery struct {
	whereContainer *WhereContainer
	table          string
}

func Exists(table string) ExistsQuery {
	return ExistsQuery{
		table:          table,
		whereContainer: &WhereContainer{},
	}
}

func (q ExistsQuery) Where(condition ...string) ExistsQuery {
	q.whereContainer.Where(condition...)

	return q
}

func (q ExistsQuery) OrWhere(condition ...string) ExistsQuery {
	q.whereContainer.OrWhere(condition...)

	return q
}

func (q ExistsQuery) Query() string {
	query := fmt.Sprintf(
		"SELECT EXISTS(SELECT true FROM %s %s)",
		q.table,
		prepareWhereQuery(q.whereContainer.conditions),
	)

	return sqlx.Rebind(sqlx.DOLLAR, query)
}
