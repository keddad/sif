package bazel

import (
	"log"
	"os/exec"
	"path/filepath"
)

type Target struct {
	label, rel_label string
}

func BuildTarget(label string, workspacePath string, envArgs []string, verbose bool) error {
	cmdArgs := []string{"build", workspacePath}
	cmdArgs = append(cmdArgs, envArgs...)

	cmd := exec.Command("bazel", cmdArgs...)
	cmd.Dir, _ = filepath.Abs(workspacePath)

	out, err := cmd.CombinedOutput()

	if err != nil && verbose {
		log.Printf("Cannot build target %s\n", label)
		log.Println(out)
	}

	return err
}

func ParseDeps(label string, workspacePath string, envArgs []string, verbose bool) ([]Target, error) {
	return nil, nil
}
