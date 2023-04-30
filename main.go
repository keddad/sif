package main

import (
	"flag"
	"log"
)

func main() {
	// Go's param parsing is horrible

	label := flag.String("label", "not-specified-lol", "Label of target to cleanup")
	workspacePath := flag.String("workspace", ".", "Workspace path")
	testTargets := flag.String("check", "", "Labels to check when modifying target. Separate with ,")
	verboseFlag := flag.Bool("v", false, "")
	params := flag.String("params", "", "Params to optimize for. Separate with ,")
	//recParams := flag.String("recparams", "", "Params to recursivly optimize dependency graph. Separate with ,")
	//recBlacklist := flag.String("recblacklist", "", "Regexp to filter out unwanted recursive optimization targets")

	bazelArgs := flag.Args()
	flag.Parse()

	if *label == "not-specified-lol" {
		log.Panic("--label argument is mandatory!")
	}

	testTargetsList := splitArgs(*testTargets)

	optimizationParamList := splitArgs(*params)

	if len(optimizationParamList) == 0 {
		log.Panic("--params argument is mandatory!")
	}

	_, err := Optimize(*label, *workspacePath, optimizationParamList, *verboseFlag, bazelArgs, testTargetsList)

	if err != nil {
		log.Fatalf("Optimization failed! %s", err)
	}
}
