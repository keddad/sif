package main

import (
	"flag"
	"log"
	"strings"
)

func main() {
	// Go's param parsing is horrible

	label := flag.String("label", "not-specified-lol", "Label of target to cleanup")
	workspacePath := flag.String("workspace", ".", "Workspace path")
	testTargets := flag.String("check", "", "Labels to check when modifying target. Separate with ,")
	verboseFlag := flag.Bool("v", false, "")
	param := flag.String("param", "", "Param to optimize")

	bazelArgs := flag.Args()
	flag.Parse()

	var testTargetsList []string

	if *testTargets != "" {
		testTargetsList = strings.Split(*testTargets, ",")
	} else {
		testTargetsList = make([]string, 0)
	}

	if *label == "not-specified-lol" {
		log.Panic("--label argument is mandatory!")
	}

	if *param == "" {
		log.Panic("--param argument is mandatory!")
	}

	_, err := Optimize(*label, *workspacePath, *param, *verboseFlag, bazelArgs, testTargetsList)

	if err != nil {
		log.Fatalf("Optimization failed! %s", err)
	}
}
