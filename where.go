package ququery

import "fmt"

var allowedOpperators = []string{
	"=",
	"!=",
	">",
	"<",
	">=",
	"<=",
	"NOT",
}

type whereContainer struct {
	conditions []whereStructure
}

func (c *whereContainer) checkOperator(column []string) string {
	op := "="
	if len(column) == 1 {
		return op
	}

	opExist := false

	for _, v := range allowedOpperators {
		if v == column[1] {
			opExist = true
		}
	}

	if !opExist {
		return op
	}

	return column[1]
}

func (c *whereContainer) where(column ...string) {
	op := c.checkOperator(column)

	c.conditions = append(c.conditions, whereStructure{
		column:   column[0],
		operator: op,
		isAnd:    true,
	})
}

func (c *whereContainer) orWhere(column ...string) {
	op := c.checkOperator(column)

	c.conditions = append(c.conditions, whereStructure{
		column:   column[0],
		operator: op,
		isAnd:    false,
	})
}

func (c *whereContainer) whereNotNull(column string) {
	c.conditions = append(c.conditions, whereStructure{
		isAnd:    true,
		isRaw:    true,
		rawQuery: column + " IS NOT NULL",
	})
}

func (c *whereContainer) orWhereNotNull(column string) {
	c.conditions = append(c.conditions, whereStructure{
		isAnd:    false,
		isRaw:    true,
		rawQuery: column + " IS NOT NULL",
	})
}

func (c *whereContainer) whereNull(column string) {
	c.conditions = append(c.conditions, whereStructure{
		isAnd:    true,
		isRaw:    true,
		rawQuery: column + " IS NULL",
	})
}

func (c *whereContainer) orWhereNull(column string) {
	c.conditions = append(c.conditions, whereStructure{
		isAnd:    false,
		isRaw:    true,
		rawQuery: column + " IS NULL",
	})
}

func (c *whereContainer) strpos(column string) {
	c.conditions = append(c.conditions, whereStructure{
		isAnd:    true,
		isRaw:    true,
		rawQuery: fmt.Sprintf("(STRPOS(%s, ?) > 0 OR ? = '')", column),
	})
}

func (c *whereContainer) orStrpos(column string) {
	c.conditions = append(c.conditions, whereStructure{
		isAnd:    false,
		isRaw:    true,
		rawQuery: fmt.Sprintf("(STRPOS(%s, ?) > 0 OR ? = '')", column),
	})
}

func (c *whereContainer) multiWhere(f func(subQuery MultiWhere) string) {
	query := f(MultiWhere{})

	c.conditions = append(c.conditions, whereStructure{
		isAnd:    true,
		isRaw:    true,
		rawQuery: query,
	})
}

func (c *whereContainer) whereInSubquery(column string, subQuery func(q SelectQuery) string) {
	c.conditions = append(c.conditions, whereStructure{
		rawQuery: fmt.Sprintf("%s IN (%s)", column, subQuery(SelectQuery{})),
		isAnd:    true,
		isRaw:    true,
	})
}
