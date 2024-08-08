package testutil

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type Testcases map[string]Testcase

// Also used to generate documentation
type Testcase struct {
	ExpectedSQL  string
	Doc          string
	Query        string
	ExpectedArgs []any
}

var (
	oneOrMoreSpace      = regexp.MustCompile(`\s+`)
	spaceAroundBrackets = regexp.MustCompile(`\s*([\(|\)])\s*`)
)

func Clean(s string) string {
	s = strings.TrimSpace(s)
	s = oneOrMoreSpace.ReplaceAllLiteralString(s, " ")
	s = spaceAroundBrackets.ReplaceAllString(s, " $1 ")

	return s
}

type FormatFunc = func(string) (string, error)

func QueryDiff(a, b string, clean FormatFunc) (string, error) {
	if clean == nil {
		clean = func(s string) (string, error) { return Clean(s), nil }
	}

	cleanA, err := clean(a)
	if err != nil {
		return "", fmt.Errorf("%s\n%w", a, err)
	}

	cleanB, err := clean(b)
	if err != nil {
		return "", fmt.Errorf("%s\n%w", b, err)
	}

	return cmp.Diff(cleanA, cleanB), nil
}

func ArgsDiff(a, b []any) string {
	return cmp.Diff(a, b)
}

func ErrDiff(a, b error) string {
	return cmp.Diff(a, b)
}

func RunTests(t *testing.T, cases Testcases, format FormatFunc) {
	t.Helper()

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			diff, err := QueryDiff(tc.ExpectedSQL, tc.Query, format)
			if err != nil {
				t.Fatalf("error: %v", err)
			}

			if diff != "" {
				fmt.Println(tc.Query)
				t.Fatalf("diff: %s", diff)
			}
		})
	}
}
