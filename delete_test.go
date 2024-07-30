package ququery_test

import (
	"testing"

	"github.com/adel-hadadi/ququery"
	"github.com/adel-hadadi/ququery/testutil"
)

func TestDeleteQuery_Where(t *testing.T) {
	testcases := testutil.Testcases{
		"delete query with condition": testutil.Testcase{
			Query:       ququery.Delete("users").Where("email").OrWhere("id"),
			ExpectedSQL: "DELETE FROM users WHERE email = $1 OR id = $2",
			Doc:         "delete user with this email or id",
		},
	}

	testutil.RunTests(t, testcases, nil)
}
