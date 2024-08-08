package ququery_test

import (
	"testing"

	"github.com/adel-hadadi/ququery"
	"github.com/adel-hadadi/ququery/testutil"
)

func TestInsertQuery_Insert(t *testing.T) {
	testcases := testutil.Testcases{
		"simple insert query": testutil.Testcase{
			Query:       ququery.Insert("users").Into("name", "email").Query(),
			ExpectedSQL: "INSERT INTO users (name, email) VALUES ($1,$2)",
			Doc:         "Insert a user with name and email",
		},
	}

	testutil.RunTests(t, testcases, nil)
}
