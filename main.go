package main

import (
	"flag"
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

	_, err := Optimize(label, workspacePath, &param, *verboseFlag, bazelArgs)

	if err != nil {
		log.Fatalf("Optimization failed! %e", err)
	}
}
