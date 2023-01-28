package main

import (
	"fmt"
	"github.com/bazelbuild/buildtools/edit"
	"github.com/keddad/sif/bazel"
	"io/ioutil"
	"log"
)

// Optimize performs optimization itself
// Returns true if optimizations took place
func Optimize(label, workspacePath, param *string, verbose bool, bazelArgs []string) (bool, error) {
	buildFile, _, target := edit.InterpretLabelForWorkspaceLocation(*workspacePath, *label)

	if verbose {
		log.Printf("\nBuildfile: %s\nTarget: %s\n", (*workspacePath)+buildFile, target)
	}

	content, err := ioutil.ReadFile(buildFile)

	if err != nil {
		return false, err
	}

	depsList, err := bazel.ExtractEntriesFromFile(content, target, *param)

	fmt.Printf("Working with %s in workspace %s, optimizing param %s. Analyzing those options:\n", *label, *workspacePath, *param)

	err = bazel.BuildTarget(*label, *workspacePath, bazelArgs, verbose)

	if err != nil {
		return false, err
	}

	removed := make([]string, 0)

	for _, elem := range depsList {
		if verbose {
			log.Printf("Trying to remove %s", elem)
		}

		contentWithoutElem, err := bazel.DeleteEntryFromFile(content, target, *param, elem)

		if err != nil {
			return false, err
		}

		err = ioutil.WriteFile(buildFile, contentWithoutElem, 0)

		if err != nil {
			return false, err
		}

		err = bazel.BuildTarget(*label, *workspacePath, bazelArgs, verbose)

		if err == nil {
			log.Printf("%s dep is redundant, removing it", elem)
			content = contentWithoutElem
			removed = append(removed, elem)
		} else if verbose {
			log.Printf("%s dep is not redundant", elem)
		}
	}

	err = ioutil.WriteFile(buildFile, content, 0)

	if err != nil {
		return false, err
	}

	if len(removed) != 0 {
		log.Print("Removed following dependencies:")
		for _, elem := range removed {
			println(elem)
		}

		return true, nil
	} else {
		log.Print("Removed no dependencies.")
		return false, nil
	}
}
