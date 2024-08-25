package ququery

import (
	"fmt"
)

var allowedOpperators = []string{
	"=",
	"!=",
	">",
	"<",
	">=",
	"<=",
	"NOT",
}

type (
	whereable interface {
		Query() string
	}

	// WhereContainer is holder that contain every condinal clauses and methods.
	WhereContainer[T whereable] struct {
		self       T
		conditions []whereStructure
	}
)

func (c *WhereContainer[T]) checkOperator(column []string) string {
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

// Where You may use the query builder's Where method to add "where" clauses to the query.
//
// Example:
//
//	query := ququery.Select("users").Where("id").Query()
//	log.Println(query) => SELECT * FROM users WHERE id = $1
//
//	query = ququery.Select("users").Where("age", ">=").Query()
//	log.Println(query) => SELECT * FROM users WHERE age >= $1
func (c *WhereContainer[T]) Where(column ...string) T {
	op := c.checkOperator(column)

	c.conditions = append(c.conditions, whereStructure{
		column:   column[0],
		operator: op,
		isAnd:    true,
	})

	return c.self
}

// OrWhere allows you to add an "or" clause to Where condition.
//
// Example:
//
//	query := ququery.Delete("users").Where("id").OrWhere("email").Query()
//	log.Println(query) => DELETE FROM users WHERE id = $1 OR email = $2
func (c *WhereContainer[T]) OrWhere(column ...string) T {
	op := c.checkOperator(column)

	c.conditions = append(c.conditions, whereStructure{
		column:   column[0],
		operator: op,
		isAnd:    false,
	})

	return c.self
}

// WhereNotNull method verifies that the column's value is not NULL:
//
// Example:
//
//	query := ququery.Delete("users").WhereNotNull("deleted_at").Query()
//	log.Println(query) => DELETE FROM users WHERE deleted_at IS NOT NULL
func (c *WhereContainer[T]) WhereNotNull(column string) T {
	c.conditions = append(c.conditions, whereStructure{
		isAnd:    true,
		isRaw:    true,
		rawQuery: column + " IS NOT NULL",
	})

	return c.self
}

// OrWhereNotNull method allows you to add an "or" clause to WhereNotNull codition
func (c *WhereContainer[T]) OrWhereNotNull(column string) T {
	c.conditions = append(c.conditions, whereStructure{
		isAnd:    false,
		isRaw:    true,
		rawQuery: column + " IS NOT NULL",
	})

	return c.self
}

// WhereNull method verifies that the value of the given column is NULL:
//
// Example:
//
//	query := ququery.Select("users").WhereNull("deleted_at").Query()
//	log.Println(query) => SELECT * FROM users WHERE deleted_at IS NULL
func (c *WhereContainer[T]) WhereNull(column string) T {
	c.conditions = append(c.conditions, whereStructure{
		isAnd:    true,
		isRaw:    true,
		rawQuery: column + " IS NULL",
	})

	return c.self
}

// OrWhereNull method allows you to add an "or" clause to WhereNull condition
//
// Example:
//
//	query := ququery.Select("users").Where("status").OrWhereNull("deleted_at").Query
//	log.Println(query) => SELECT * FROM users WHERE status = $1 OR deleted_at IS NULL
func (c *WhereContainer[T]) OrWhereNull(column string) T {
	c.conditions = append(c.conditions, whereStructure{
		isAnd:    false,
		isRaw:    true,
		rawQuery: column + " IS NULL",
	})

	return c.self
}

// WhereLike method allows you to add "LIKE" clauses to your query from pattern matchinga.
//
// Example:
//
//	query := ququery.Select("users").WhereLike("name").Query()
//	log.Println(query) => SELECT * FROM users WHERE name LIKE $1
func (c *WhereContainer[T]) WhereLike(column string) T {
	c.conditions = append(c.conditions, whereStructure{
		column:   column,
		operator: "LIKE",
		isAnd:    false,
	})

	return c.self
}

// OrWhereLike method allows you to add an "or" clause with a LIKE condition
//
// Example:
//
//	query := ququery.Select("users").Where("id").OrWhereLike("name").Query()
//	log.Println(query) => SELECT * FROM users WHERE id = $1 OR name LIKE $2
func (c *WhereContainer[T]) OrWhereLike(column string) T {
	c.conditions = append(c.conditions, whereStructure{
		column:   column,
		operator: "LIKE",
		isAnd:    false,
	})

	return c.self
}

// Strpos method is more like whereLike method,
// but the difference is that strpos method is used by postgresql users for full text search.
//
// Example:
//
//	query := ququery.Select("users").Strpos("name").Query()
//	log.Println(query) => SELECT * FROM users WHERE (STRPOS(name, $1) > 0 or $2 = '')
func (c *WhereContainer[T]) Strpos(column string) T {
	c.conditions = append(c.conditions, whereStructure{
		isAnd:    true,
		isRaw:    true,
		rawQuery: fmt.Sprintf("(STRPOS(%s, ?) > 0 OR ? = '')", column),
	})

	return c.self
}

// OrStrpos method allows you to add an "or" clause to Strpos condition.
//
// Example:
//
//	query := ququery.Select("users").Where("id").Strpos("name").Query()
//	log.Println(query) => SELECT * FROM users WHERE id = $1 OR (STRPOS(name, $2) > 0 or $3 = '')
func (c *WhereContainer[T]) OrStrpos(column string) T {
	c.conditions = append(c.conditions, whereStructure{
		isAnd:    false,
		isRaw:    true,
		rawQuery: fmt.Sprintf("(STRPOS(%s, ?) > 0 OR ? = '')", column),
	})

	return c.self
}

// WhereGroup Sometimes you may need to group several "where" clauses within
// parentheses in order to achieve your query's desired logical grouping.
// In fact, you should generally always group calls to the orWhere method in
// parentheses in order to avoid unexpected query behavior.
// To accomplish this, you may user this method:
//
// Example:
//
//		     query := ququery.Select("users").WhereGroup(func(subQuery ququery.MultiWhere) string {
//			    return subQuery.Where("email").
//				    Where("role_id").
//				    OrWhere("type").
//				    Query()
//		    }).Query(),
//	        log.Println(query) => SELECT * FROM users WHERE ( email = $1 AND role_id = $2 OR type = $3)
func (c *WhereContainer[T]) WhereGroup(f func(subQuery MultiWhere) string) T {
	query := f(MultiWhere{})

	c.conditions = append(c.conditions, whereStructure{
		isAnd:    true,
		isRaw:    true,
		rawQuery: query,
	})

	return c.self
}

// WhereInSubquery Sometimes you may need to construct a "where" clause that compares
// the results of a subquery to a given value. You may accomplish this by
// passing a closure and a value to the where method.
//
// Example:
//
//	     query := ququery.Select("users").WhereInSubquery("users.id", func(q ququery.SelectQuery) string {
//		    return q.Table("orders").
//			    Columns("user_id").
//			    OrderBy("total_price", "desc").
//			    Limit().
//			    Query()
//	        }).
//	        Query()
//
//	    log.Println(query) => SELECT * FROM users WHERE users.id IN (SELECT user_id FROM orders ORDER BY total_price DESC LIMIT $1)
func (c *WhereContainer[T]) WhereInSubquery(column string, subQuery func(q SelectQuery) string) T {
	c.conditions = append(c.conditions, whereStructure{
		rawQuery: fmt.Sprintf("%s IN (%s)", column, subQuery(SelectQuery{
			withoutRebinding: true,
		})),
		isAnd: true,
		isRaw: true,
	})

	return c.self
}
