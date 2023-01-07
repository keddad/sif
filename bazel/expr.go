package bazel

import (
	"errors"
	"github.com/bazelbuild/buildtools/build"
	"github.com/bazelbuild/buildtools/edit"
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

// extractEntriesFromRule extracts contents of param (which should be an array, like deps) from Rule.
func extractEntriesFromRule(rule *build.Rule, name string) ([]string, error) {
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

// ExtractEntriesFromFile extracts contents of param (which should be an array, like deps) of Rule from file contents.
func ExtractEntriesFromFile(contents []byte, ruleName string, paramName string) ([]string, error) {
	origBuildFile, err := build.ParseBuild("", contents)

	if err != nil {
		return nil, err
	}

	targetRule := edit.FindRuleByName(origBuildFile, ruleName)

	depsList, err := extractEntriesFromRule(targetRule, paramName)

	if err != nil {
		return nil, err
	}

	return depsList, nil
}
