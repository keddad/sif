package main

import (
	"flag"
	"github.com/bazelbuild/buildtools/build"
	"github.com/bazelbuild/buildtools/edit"
	"github.com/keddad/sif/bazel"
	"io/ioutil"
	"log"
)

func main() {
	label := flag.String("label", "not-specified-lol", "Label of target to cleanup")
	workspacePath := flag.String("workspace", ".", "Workspace path")
	verboseFlag := flag.Bool("v", false, "")

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

	origBuildFile, err := build.ParseBuild(buildFile, content)

	if err != nil {
		panic(err)
	}

	targetRule := edit.FindRuleByName(origBuildFile, target)

	depsList, err := bazel.ExtractEntries(targetRule, "deps")

	if err != nil {
		panic(err)
	}

	print(depsList)

	// Build target with Bazel
	// This populates the cache to make further operations faster, and ensures target builds before changes
	bazel.BuildTarget(*label, *workspacePath, bazelArgs, *verboseFlag)
}
