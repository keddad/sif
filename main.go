package main

import (
	"log"

	"github.com/jessevdk/go-flags"
	"github.com/keddad/sif/bazel"
)

type cmdArgs struct {
	label         string   `required:"true" short:"l"`
	workspacePath string   `description:"Path to Bazel workspace" default:"." short:"wp"`
	isVebrose     bool     `default:"false" short:"v"`
	bazelArgs     []string `long:"Additonal Bazel arguments" positional-args:`
}

func main() {
	var args cmdArgs
	parser := flags.NewParser(&args, flags.Default)
	_, err := parser.Parse()

	if err != nil {
		log.Fatal(err)
	}

	bazel.BuildTarget(args.label, args.workspacePath, args.bazelArgs, args.isVebrose)
}
