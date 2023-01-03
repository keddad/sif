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
	label := flag.String("label", "//main:hello-greet", "Label of target to cleanup")
	workspacePath := flag.String("workspace", "test/cppexample/", "Workspace path")
	verboseFlag := flag.Bool("v", false, "")

	bazelArgs := flag.Args()

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

	print(origBuildFile)

	// Build target with Bazel
	// This populates the cache to make further operations faster, and ensures target builds before changes
	bazel.BuildTarget(*label, *workspacePath, bazelArgs, *verboseFlag)
}
