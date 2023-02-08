package main

import (
	"flag"
	"log"
	"strings"
)

func main() {
	label := flag.String("label", "not-specified-lol", "Label of target to cleanup")
	workspacePath := flag.String("workspace", ".", "Workspace path")
	testTargets := flag.String("check", "", "Labels to chpatheck when modifying target. Separete with ,")
	verboseFlag := flag.Bool("v", false, "")
	param := flag.String("param", "", "Param to optimize")

	bazelArgs := flag.Args()
	flag.Parse()

	splitTargets := strings.Split(*testTargets, ",")

	if *label == "not-specified-lol" {
		log.Panic("--label argument is mandatory!")
	}

	if *param == "" {
		log.Panic("--param argument is mandatory!")
	}

	_, err := Optimize(*label, *workspacePath, *param, *verboseFlag, bazelArgs, splitTargets)

	if err != nil {
		log.Fatalf("Optimization failed! %e", err)
	}
}
