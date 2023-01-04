package bazel

import (
	"errors"
	"github.com/bazelbuild/buildtools/build"
)

func listStrings(expr *build.AssignExpr) ([]string, error) {
	ret := make([]string, 0)

	switch expr.RHS.(type) {
	case *build.ListExpr:
	default:
		return nil, errors.New("param is not a list")
	}

	for _, elem := range expr.RHS.(*build.ListExpr).List {
		switch elem.(type) {
		case *build.StringExpr:
			ret = append(ret, elem.(*build.StringExpr).Token)
		default:
			return nil, errors.New("list member is not a string")
		}
	}

	return ret, nil
}

// ExtractEntries extracts contents of param (which should be an array, like deps) from Rule.
func ExtractEntries(rule *build.Rule, name string) ([]string, error) {
	for _, expr := range rule.Call.List {
		// Dynamic casts from interface to implementation are always nasty

		switch expr.(type) {
		case *build.AssignExpr:
			if expr.(*build.AssignExpr).LHS.(*build.Ident).Name == name {
				return listStrings(expr.(*build.AssignExpr))
			}
		default:
			continue
		}
	}

	return nil, errors.New("no such param")
}
