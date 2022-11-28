package main

import (
	"flag"

	"github.com/keddad/sif/bazel"
)

func main() {
	label := flag.String("label", "//main:hello-greet", "Label of target to cleanup")
	workspacePath := flag.String("workspace", "test/cppexample/", "Workspace path")
	verboseFlag := flag.Bool("v", true, "")

	bazelArgs := flag.Args()

	bazel.BuildTarget(*label, *workspacePath, bazelArgs, *verboseFlag)
}
