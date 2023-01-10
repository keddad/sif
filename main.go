package main

import (
	"flag"
	"fmt"
	"github.com/bazelbuild/buildtools/edit"
	"github.com/keddad/sif/bazel"
	"io/ioutil"
	"log"
)

func main() {
	label := flag.String("label", "not-specified-lol", "Label of target to cleanup")
	workspacePath := flag.String("workspace", ".", "Workspace path")
	verboseFlag := flag.Bool("v", false, "")
	param := "deps"

	bazelArgs := flag.Args()
	flag.Parse()

	if *label == "not-specified-lol" {
		log.Panic("-label argument is mandatory!")
	}

	buildFile, _, target := edit.InterpretLabelForWorkspaceLocation(*workspacePath, *label)

	if *verboseFlag {
		log.Printf("\nBuildfile: %s\nTarget: %s\n", *workspacePath+buildFile, target)
	}

	content, err := ioutil.ReadFile(buildFile)

	if err != nil {
		panic(err)
	}

	depsList, err := bazel.ExtractEntriesFromFile(content, target, param)

	fmt.Printf("Working with %s in workspace %s, optimizing param %s. Analyzing those options:\n", *label, *workspacePath, param)

	for _, elem := range depsList {
		println(elem)
	}

	err = bazel.BuildTarget(*label, *workspacePath, bazelArgs, *verboseFlag)

	if err != nil {
		log.Fatalf("Can't build target with it's current deps!; %e", err)
	}

	removed := make([]string, 0)

	for _, elem := range depsList {
		if *verboseFlag {
			log.Printf("Trying to remove %s", elem)
		}

		contentWithoutElem, err := bazel.DeleteEntryFromFile(content, target, param, elem)

		if err != nil {
			log.Fatalf("Can't remove %s from file; %e", elem, err)
		}

		err = ioutil.WriteFile(buildFile, contentWithoutElem, 0)

		if err != nil {
			log.Fatalf("Can't write BUILD file; %e", err)
		}

		err = bazel.BuildTarget(*label, *workspacePath, bazelArgs, *verboseFlag)

		if err == nil {
			log.Printf("%s dep is redundant, removing it", elem)
			content = contentWithoutElem
			removed = append(removed, elem)
		} else if *verboseFlag {
			log.Printf("%s dep is not redundant", elem)
		}
	}

	err = ioutil.WriteFile(buildFile, content, 0)

	if err != nil {
		log.Fatalf("Can't write BUILD file; %e", err)
	}

	if len(removed) != 0 {
		log.Print("Removed following dependencies:")
		for _, elem := range removed {
			println(elem)
		}
	} else {
		log.Print("Removed no dependencies.")
	}
}
