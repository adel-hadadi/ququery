package ququery

var allowedOpperators = []string{
	"=",
	"!=",
	">",
	"<",
	">=",
	"<=",
	"NOT",
}

type WhereContainer struct {
	conditions []whereStructure
}

func (c *WhereContainer) checkOperator(column []string) string {
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

func (c *WhereContainer) Where(column ...string) {
	op := c.checkOperator(column)

	c.conditions = append(c.conditions, whereStructure{
		column:   column[0],
		operator: op,
		rawQuery: "",
		isAnd:    true,
		isRaw:    false,
	})
}

func (c *WhereContainer) OrWhere(column ...string) {
	op := c.checkOperator(column)

	c.conditions = append(c.conditions, whereStructure{
		column:   column[0],
		operator: op,
		rawQuery: "",
		isAnd:    false,
		isRaw:    false,
	})
}
