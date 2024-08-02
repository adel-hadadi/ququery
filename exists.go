package ququery

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ExistsQuery struct {
	whereContainer *whereContainer
	table          string
}

func Exists(table string) ExistsQuery {
	return ExistsQuery{
		table:          table,
		whereContainer: &whereContainer{},
	}
}

func (q ExistsQuery) Where(condition ...string) ExistsQuery {
	q.whereContainer.where(condition...)

	return q
}

func (q ExistsQuery) OrWhere(condition ...string) ExistsQuery {
	q.whereContainer.orWhere(condition...)

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
