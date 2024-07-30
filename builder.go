package ququery

type Query interface {
	Query() string
}

type whereStructure struct {
	column   string
	operator string
	rawQuery string
	isAnd    bool
	isRaw    bool
}

func prepareMultiWhereConditions(conditions []whereStructure) string {
	var query string

	for i, condition := range conditions {
		if i > 0 {
			if condition.isAnd {
				query += " AND"
			} else {
				query += " OR"
			}
		}

		if condition.isRaw {
			query += " " + condition.rawQuery
		} else {
			query += " " + condition.column + " " + condition.operator + " ?"
		}
	}

	return query
}

func prepareWhereQuery(wheres []whereStructure) string {
	var conditions string

	for i, where := range wheres {
		if i == 0 {
			conditions += "WHERE"
		}

		if i > 0 {
			if where.isAnd {
				conditions += " AND"
			} else {
				conditions += " OR"
			}
		}

		if where.isRaw {
			conditions += " " + where.rawQuery
		} else {
			conditions += " " + where.column + " " + where.operator + " ?"
		}
	}

	return conditions
}

func CountOver() string {
	return "COUNT(*) OVER()"
}

func findTableFromEntity(entity string) string {
	var table string

	if entity[len(entity)-1:] == "y" {
		table = entity[:len(entity)-1] + "ies"
	} else {
		table = entity + "s"
	}

	return table
}
