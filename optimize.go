package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/bazelbuild/buildtools/edit"
	"github.com/bazelbuild/buildtools/labels"
	"github.com/keddad/sif/bazel"
)

// Optimize performs optimization itself
// Returns true if optimizations took place
func Optimize(label, workspacePath string, params []string, verbose bool, bazelArgs []string, testLabels []string, recParams []string, recBlaclist string) (bool, error) {
	if recBlaclist != "" {
		isInBlacklist, err := regexp.Match(recBlaclist, []byte(label))

		if err != nil {
			return false, err
		}

		if isInBlacklist {
			log.Printf("Skipping label %s, blacklist", label)
			return false, nil
		}
	}

	buildFile, _, target := edit.InterpretLabelForWorkspaceLocation(workspacePath, label)

	if verbose {
		log.Printf("\nBuildfile: %s\nTarget: %s\n", (workspacePath)+buildFile, target)
	}

	err := bazel.CheckTarget(label, workspacePath, bazelArgs, verbose)

	if err != nil {
		fmt.Printf("Can't build optimization target %s\n", label)
		return false, err
	}

	for _, elem := range testLabels {
		err = bazel.CheckTarget(elem, workspacePath, bazelArgs, verbose)

		if err != nil {
			fmt.Printf("Can't build test target %s\n", label)
			return false, err
		}
	}

	removedDeps := false

	for _, param := range params {
		log.Printf("Working with %s in workspace %s, optimizing param %s.\n", label, workspacePath, param)

		removed, err := optimizeParam(buildFile, target, param, label, workspacePath, verbose, bazelArgs, testLabels)

		if err != nil {
			if err == bazel.ErrNoSuchParam {
				continue
			}

			return false, err
		}

		if len(removed) != 0 {
			log.Print("Removed following dependencies:")
			for _, elem := range removed {
				log.Print(elem)
			}

			removedDeps = true
		} else {
			log.Printf("Removed no dependencies for param.")
		}
	}

	recTarges, err := recTargets(buildFile, target, label, workspacePath, recParams, verbose)

	if err != nil {
		return false, err
	}

	for _, target := range recTarges {

		targetRemovedDeps, err := Optimize(target, workspacePath, params, verbose, bazelArgs, append(testLabels, label), recParams, recBlaclist)

		if err != nil {
			return false, err
		}

		removedDeps = removedDeps || targetRemovedDeps
	}

	return removedDeps, nil
}

func recTargets(buildFile, target, label, workspacePath string, recParams []string, verbose bool) ([]string, error) {
	content, err := ioutil.ReadFile(buildFile)

	if err != nil {
		return nil, err
	}

	targets := make([]string, 0)
	currentRuleName := labels.Parse(label)

	for _, param := range recParams {
		paramTargets, err := bazel.ExtractEntriesFromFile(content, currentRuleName.Target, param)

		if err != nil {
			if err == bazel.ErrNoSuchParam {
				continue
			}

			return targets, err
		}

		for _, target := range paramTargets {
			target = target[1 : len(target)-1] // Targets are with quotes, ":target"

			if !strings.HasPrefix(target, "//") && !strings.HasPrefix(target, "@") {
				target = labels.ParseRelative(target, currentRuleName.Package).Format()
			}

			targets = append(targets, target)
		}

		if verbose {
			log.Printf("Added labels %#v for target %s", paramTargets, param)
		}

	}

	return targets, nil
}

func optimizeParam(buildFile, target, param, label, workspacePath string, verbose bool, bazelArgs, testLabels []string) ([]string, error) {
	content, err := ioutil.ReadFile(buildFile)

	if err != nil {
		return nil, err
	}

	depsList, err := bazel.ExtractEntriesFromFile(content, target, param)

	if err != nil {
		return nil, err
	}

	removed := make([]string, 0)

	for _, elem := range depsList {
		if verbose {
			log.Printf("Trying to remove %s", elem)
		}

		contentWithoutElem, err := bazel.DeleteEntryFromFile(content, target, param, elem)

		if err != nil {
			return nil, err
		}

		err = ioutil.WriteFile(buildFile, contentWithoutElem, 0)

		if err != nil {
			return nil, err
		}

		err = bazel.CheckTarget(label, workspacePath, bazelArgs, verbose)

		if err == nil {
			for _, testLabel := range testLabels {
				err = bazel.CheckTarget(testLabel, workspacePath, bazelArgs, verbose)

				if err != nil {
					if verbose {
						log.Printf("Test label %s failed when removing dep %s", testLabel, elem)
					}

					break
				}
			}
		}

		if err == nil {
			if verbose {
				log.Printf("%s dep is redundant, removing it", elem)
			}

			content = contentWithoutElem
			removed = append(removed, elem)
		} else if verbose {
			log.Printf("%s dep is not redundant", elem)
		}
	}

	err = ioutil.WriteFile(buildFile, content, 0)

	if err != nil {
		return nil, err
	}

	return removed, err
}
