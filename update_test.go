package ququery_test

import (
	"testing"

	"github.com/adel-hadadi/ququery"
	"github.com/adel-hadadi/ququery/testutil"
)

func TestUpdateQuery_Update(t *testing.T) {
	testcases := testutil.Testcases{
		"simple update query": testutil.Testcase{
			Query: ququery.Update("users").
				Where("id").
				OrWhere("email").
				Set("first_name", "last_name").
				Query(),
			ExpectedSQL: "UPDATE users SET first_name = $1, last_name = $2 WHERE id = $3 OR email = $4",
			Doc:         "simple update query",
		},
	}

	testutil.RunTests(t, testcases, nil)
}
