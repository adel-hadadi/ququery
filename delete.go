package ququery

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DeleteQuery struct {
	table      string
	conditions []whereStructure
}

func Delete(table string) DeleteQuery {
	return DeleteQuery{
		table: table,
	}
}

func (q DeleteQuery) Where(condition ...string) DeleteQuery {
	operator := "="
	if len(condition) > 1 {
		operator = condition[1]
	}

	q.conditions = append(q.conditions, whereStructure{
		column:   condition[0],
		operator: operator,
		isAnd:    true,
	})

	return q
}

func (q DeleteQuery) OrWhere(condition ...string) DeleteQuery {
	operator := "="
	if len(condition) > 1 {
		operator = condition[1]
	}

	q.conditions = append(q.conditions, whereStructure{
		column:   condition[0],
		operator: operator,
		isAnd:    false,
	})

	return q
}

func (q DeleteQuery) Query() string {
	query := fmt.Sprintf(
		`DELETE FROM %s %s`,
		q.table,
		prepareWhereQuery(q.conditions),
	)

	return sqlx.Rebind(sqlx.DOLLAR, query)
}
