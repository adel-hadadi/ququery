package ququery

import "fmt"

type MultiWhere struct {
	wheres []whereStructure
}

func (w MultiWhere) Where(condition ...string) MultiWhere {
	operator := "="
	if len(condition) > 1 {
		operator = condition[1]
	}

	w.wheres = append(w.wheres, whereStructure{
		column:   condition[0],
		operator: operator,
		isAnd:    true,
	})

	return w
}

func (w MultiWhere) OrWhere(condition ...string) MultiWhere {
	operator := "="
	if len(condition) > 1 {
		operator = condition[1]
	}

	w.wheres = append(w.wheres, whereStructure{
		column:   condition[0],
		operator: operator,
		isAnd:    false,
	})

	return w
}

func (w MultiWhere) Query() string {
	return fmt.Sprintf("(%s)", prepareMultiWhereConditions(w.wheres))
}
