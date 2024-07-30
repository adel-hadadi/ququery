package ququery

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UpdateQuery struct {
	table      string
	columns    []string
	conditions []whereStructure
}

func Update(table string) UpdateQuery {
	return UpdateQuery{
		table: table,
	}
}

func (q UpdateQuery) Set(columns ...string) UpdateQuery {
	q.columns = append(q.columns, columns...)

	return q
}

func (q UpdateQuery) Where(condition ...string) UpdateQuery {
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

func (q UpdateQuery) OrWhere(condition ...string) UpdateQuery {
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

func (q UpdateQuery) Query() string {
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
