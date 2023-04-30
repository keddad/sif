package bazel

import (
	"errors"
	"io/ioutil"
	"strings"

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

// Buildtool's API is quite poor, since it is not really public and is mostly designed to be used by Google's CLI apps
// If it wasn't, this entire file would be redundant. Still better than writing it all from scratch
func findAssignExpr(rule *build.Rule, name string) *build.AssignExpr {
	for _, expr := range rule.Call.List {
		// Dynamic casts from interface to implementation are always nasty

		switch expr.(type) {
		case *build.AssignExpr:
			if expr.(*build.AssignExpr).LHS.(*build.Ident).Name == name {
				return expr.(*build.AssignExpr)
			}
		default:
			continue
		}
	}

	return nil
}

// extractEntriesFromRule extracts contents of param (which should be an array, like deps) from Rule.
func extractEntriesFromRule(rule *build.Rule, name string) ([]string, error) {
	expr := findAssignExpr(rule, name)

	if expr == nil {
		return nil, ErrNoSuchParam
	}

	return listStrings(expr)
}

// dumbListRemoval removes a StringExpr from ListExpr without package checks and pointer magic
// Pretty much a simplified copy of RemoveFromList
func dumbListRemoval(li *build.ListExpr, toDelete string) error {
	var all []build.Expr
	deleted := false

	for _, elem := range li.List {
		if str, ok := elem.(*build.StringExpr); ok {
			if str.Token == toDelete {
				deleted = true
				continue
			}
		}
		all = append(all, elem)
	}

	li.List = all

	if !deleted {
		return errors.New("no such list element")
	}

	return nil
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

func DeleteEntryFromFile(contents []byte, ruleName string, paramName string, optionToDelete string) ([]byte, error) {
	origBuildFile, err := build.ParseBuild("", contents)

	if err != nil {
		return nil, err
	}

	targetRule := edit.FindRuleByName(origBuildFile, ruleName)

	targetExpr := findAssignExpr(targetRule, paramName)

	if targetExpr == nil {
		return nil, ErrNoSuchParam
	}

	if targetExpr.RHS == nil {
		return nil, errors.New("wtf") // What a Terrible Failure! Probably can't happen with valid file
	}

	err = dumbListRemoval(targetExpr.RHS.(*build.ListExpr), optionToDelete)
	newFile := build.Format(origBuildFile)

	return newFile, err
}

// Returns a rule of a taget (cc_binary, etc)
func getTargetRuleKind(contents []byte, ruleName string) (string, error) {
	origBuildFile, err := build.ParseBuild("", contents)

	if err != nil {
		return "", err
	}

	targetRule := edit.FindRuleByName(origBuildFile, ruleName)

	if targetRule == nil {
		return "", ErrNoSuchRule
	}

	return targetRule.Kind(), nil
}

func IsTest(label string, workspace string) (bool, error) {
	buildFile, _, target := edit.InterpretLabelForWorkspaceLocation(workspace, label)

	content, err := ioutil.ReadFile(buildFile)

	if err != nil {
		return false, err
	}

	name, err := getTargetRuleKind(content, target)

	if err != nil {
		return false, err
	}

	// This is not ideal, but it mostly works
	// Should probably check for actual providers, but it can break with test_suite (?)
	return strings.Contains("test", name), nil
}
