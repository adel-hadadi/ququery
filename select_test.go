package ququery_test

import (
	"testing"

	"github.com/adel-hadadi/ququery"
	"github.com/adel-hadadi/ququery/testutil"
)

func TestSelectQuery_Columns(t *testing.T) {
	testcases := testutil.Testcases{
		"select columns": testutil.Testcase{
			Query:        ququery.Select("users").Columns("id", "name").Query(),
			ExpectedSQL:  "SELECT id, name FROM users",
			Doc:          "select columns query",
			ExpectedArgs: nil,
		},

		"select all columns when not determined them": testutil.Testcase{
			Query:       ququery.Select("users").Query(),
			ExpectedSQL: "SELECT * FROM users",
			Doc:         "select all columns of table when columns not defineded",
		},
	}

	testutil.RunTests(t, testcases, nil)
}

func TestSelectQuery_Where(t *testing.T) {
	testcases := testutil.Testcases{
		"simple select with where": testutil.Testcase{
			Query:        ququery.Select("users").Where("id").Query(),
			ExpectedSQL:  "SELECT * FROM users WHERE id = $1",
			ExpectedArgs: nil,
			Doc:          "Simple select query with where",
		},

		"where and orWhere": testutil.Testcase{
			Query:        ququery.Select("users").Where("status").OrWhere("role_id").Query(),
			ExpectedSQL:  "SELECT * FROM users WHERE status = $1 OR role_id = $2",
			ExpectedArgs: nil,
			Doc:          "Select query with where and or where",
		},
	}

	testutil.RunTests(t, testcases, nil)
}

func TestSelectQuery_Join(t *testing.T) {
	testcases := testutil.Testcases{
		"inner join": testutil.Testcase{
			Query: ququery.Select("users").
				Columns("id").
				Join("wallets", "wallets.id = users.id").
				Query(),
			ExpectedSQL:  "SELECT id FROM users INNER JOIN wallets ON wallets.id = users.id",
			ExpectedArgs: nil,
			Doc:          "simple select on users table and load users wallet",
		},

		"right join": testutil.Testcase{
			Query: ququery.Select("users").
				Columns("id").
				RightJoin("wallets", "wallets.id = users.id").
				Query(),
			ExpectedSQL:  "SELECT id FROM users RIGHT JOIN wallets ON wallets.id = users.id",
			ExpectedArgs: nil,
			Doc:          "simple select on users table and load users wallet with right join",
		},

		"left join": testutil.Testcase{
			Query: ququery.Select("users").
				Columns("id").
				LeftJoin("wallets", "wallets.id = users.id").
				Query(),
			ExpectedSQL:  "SELECT id FROM users LEFT JOIN wallets ON wallets.id = users.id",
			ExpectedArgs: nil,
			Doc:          "simple select on users table and load users wallet with right join",
		},
	}

	testutil.RunTests(t, testcases, nil)
}

func TestSelectQuery_Strpos(t *testing.T) {
	testcases := testutil.Testcases{
		"select query with check string posititon": testutil.Testcase{
			Query:       ququery.Select("users").Strpos("name").Query(),
			ExpectedSQL: "SELECT * FROM users WHERE (STRPOS(name, $1) > 0 OR $2 = '')",
			Doc:         "Select query with check string position",
		},

		"select query with conditions and or string position": testutil.Testcase{
			Query:       ququery.Select("users").Where("id").OrStrpos("email").Query(),
			ExpectedSQL: "SELECT * FROM users WHERE id = $1 OR (STRPOS(email, $2) > 0 OR $3 = '')",
			Doc:         "Select query with where and or string position check",
		},
	}

	testutil.RunTests(t, testcases, nil)
}

func TestSelectQuery_WhereGroup(t *testing.T) {
	testcases := testutil.Testcases{
		"select query with multi where": testutil.Testcase{
			Query: ququery.Select("users").WhereGroup(func(subQuery ququery.MultiWhere) string {
				return subQuery.Where("email").
					Where("role_id").
					OrWhere("type").
					Query()
			}).Query(),
			ExpectedSQL: "SELECT * FROM users WHERE ( email = $1 AND role_id = $2 OR type = $3)",
			Doc:         "select query with group of where conditions",
		},
	}

	testutil.RunTests(t, testcases, nil)
}

func TestSelectQuery_WhereInSubQuery(t *testing.T) {
	testcases := testutil.Testcases{
		"select query with where in subquery": testutil.Testcase{
			Query: ququery.Select("users").WhereInSubquery("users.id", func(q ququery.SelectQuery) string {
				return q.Table("orders").
					Columns("user_id").
					OrderBy("total_price", "desc").
					Limit().
					Query()
			}).Query(),
			ExpectedSQL: "SELECT * FROM users WHERE users.id IN (SELECT user_id FROM orders ORDER BY total_price DESC LIMIT $1)",
			Doc:         "Select query with where in subquery to other table",
		},
	}

	testutil.RunTests(t, testcases, nil)
}

func TestSelectQuery_With(t *testing.T) {
	testcases := testutil.Testcases{
		"select query with load realations": testutil.Testcase{
			Query:       ququery.Select("users").With("role").Query(),
			ExpectedSQL: "SELECT * FROM users LEFT JOIN roles ON roles.id = users.role_id",
			Doc:         "Select all users with role",
		},
		"select query with load relations that contain y in name": testutil.Testcase{
			Query:       ququery.Select("users").With("city").Query(),
			ExpectedSQL: "SELECT * FROM users LEFT JOIN cities ON cities.id = users.city_id",
			Doc:         "Select all users with city",
		},
	}

	testutil.RunTests(t, testcases, nil)
}

func TestSelectQuery_LimitAndOffset(t *testing.T) {
	testcases := testutil.Testcases{
		"select query with limit and offset": testutil.Testcase{
			Query:       ququery.Select("users").Limit().Offset().Query(),
			ExpectedSQL: "SELECT * FROM users LIMIT $1 OFFSET $2",
			Doc:         "select query by limit and offset",
		},
	}

	testutil.RunTests(t, testcases, nil)
}

func TestSelectQuery_WhereNotNull(t *testing.T) {
	testcases := testutil.Testcases{
		"select query with where not null": testutil.Testcase{
			Query:       ququery.Select("users").WhereNotNull("email").Query(),
			ExpectedSQL: "SELECT * FROM users WHERE email IS NOT NULL",
			Doc:         "select users where email value is not null",
		},
		"select query with or where not null": testutil.Testcase{
			Query:       ququery.Select("users").Where("email").OrWhereNotNull("role_id").Query(),
			ExpectedSQL: "SELECT * FROM users WHERE email = $1 OR role_id IS NOT NULL",
			Doc:         "select users where email value is not null",
		},
	}

	testutil.RunTests(t, testcases, nil)
}

func TestSelectQuery_WhereNull(t *testing.T) {
	testcases := testutil.Testcases{
		"select query with where null": testutil.Testcase{
			Query:       ququery.Select("users").WhereNull("deleted_at").Query(),
			ExpectedSQL: "SELECT * FROM users WHERE deleted_at IS NULL",
			Doc:         "select users where deleted_at is null",
		},
		"select query with or where null": testutil.Testcase{
			Query:       ququery.Select("users").Where("status").OrWhereNull("deleted_at").Query(),
			ExpectedSQL: "SELECT * FROM users WHERE status = $1 OR deleted_at IS NULL",
			Doc:         "select users where status is specific value or  deleted_at is null",
		},
	}

	testutil.RunTests(t, testcases, nil)
}
