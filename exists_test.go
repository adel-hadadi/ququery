package ququery_test

import (
	"testing"

	"github.com/adel-hadadi/ququery"
	"github.com/adel-hadadi/ququery/testutil"
)

func TestExistsQuery_Exists(t *testing.T) {
	testcases := testutil.Testcases{
		"simple check exists query with condition": testutil.Testcase{
			Query:       ququery.Exists("users").Where("email").OrWhere("id"),
			ExpectedSQL: "SELECT EXISTS(SELECT true FROM users WHERE email = $1 OR id = $2)",
			Doc:         "check user with this email exists or not",
		},
	}

	testutil.RunTests(t, testcases, nil)
}
