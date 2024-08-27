package ququery_test

import (
	"testing"

	"github.com/adel-hadadi/ququery"
	"github.com/adel-hadadi/ququery/testutil"
)

func TestWhereContainer_Where(t *testing.T) {
	testcases := testutil.Testcases{
		"simple where query": {
			Query:       ququery.Select("users").Where("id").Query(),
			ExpectedSQL: "SELECT * FROM users WHERE id = $1",
			Doc:         "simple select query that return user where id equal to x",
		},
		"where query with or where": {
			Query: ququery.Select("users").
				Where("email").
				OrWhere("phone_number").
				Query(),
			ExpectedSQL: "SELECT * FROM users WHERE email = $1 OR phone_number = $2",
			Doc:         "select users where email equal x or phone_number equal y",
		},
		"where query with custome operator": {
			Query: ququery.Update("users").
				Where("status", "!=").
				Set("status").
				Query(),
			ExpectedSQL: "UPDATE users SET status = $1 WHERE status != $2 ",
			Doc:         "update user statuses where status is not equal to x",
		},
	}

	testutil.RunTests(t, testcases, nil)
}

func TestWhereContainer_TextSearch(t *testing.T) {
	testcases := testutil.Testcases{
		"query with whereLike": {
			Query:       ququery.Select("users").WhereLike("name").Query(),
			ExpectedSQL: "SELECT * FROM users WHERE name LIKE $1",
			Doc:         "select users where name is like x",
		},
		"query with orWhereLike": {
			Query: ququery.Select("users").
				Where("id").
				OrWhereLike("email").
				Query(),
			ExpectedSQL: "SELECT * FROM users WHERE id = $1 OR email LIKE $2",
			Doc:         "select users where id is x or email is like y",
		},
		"query with strPos and orStrPos": {
			Query: ququery.Select("users").Strpos("name").
				OrStrpos("email").
				Query(),
			ExpectedSQL: "SELECT * FROM users WHERE (STRPOS(name, $1) > 0 OR $2 = '') OR (STRPOS(email, $3) > 0 OR $4 = '')",
			Doc:         "select users where name is like x and check x is not empty or email is like y",
		},
	}

	testutil.RunTests(t, testcases, nil)
}

func TestWhereContainer_WhereInSubquery(t *testing.T) {
	testcases := testutil.Testcases{
		"query with whereInSubquery": {
			Query: ququery.Select("users").WhereInSubquery("users.id", func(q ququery.SelectQuery) string {
				return q.Table("orders").
					Columns("user_id").
					OrderBy("total_price", ququery.DESC).
					Limit().
					Query()
			}).
				Where("id").
				Query(),
			ExpectedSQL: "SELECT * FROM users WHERE users.id IN (SELECT user_id FROM orders ORDER BY total_price DESC LIMIT $1) AND id = $2",
			Doc:         "select query with where in subquery to other table",
		},
		"query with orWhereInSubquery": {
			Query: ququery.Select("users").Where("role_id").OrWhereInSubquery("users.id", func(q ququery.SelectQuery) string {
				return q.Table("orders").
					Columns("user_id").
					OrderBy("orders.id", ququery.ASC).
					Query()
			}).Query(),
			ExpectedSQL: "SELECT * FROM users WHERE role_id = $1 OR users.id IN (SELECT user_id FROM orders ORDER BY orders.id ASC)",
		},
	}

	testutil.RunTests(t, testcases, nil)
}
