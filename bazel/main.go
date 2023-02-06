package bazel

import (
	"log"
	"os/exec"
	"path/filepath"
)

// If target is of test rule, runs bazel test
// Otherwise build target
func CheckTarget(label string, workspacePath string, envArgs []string, verbose bool) error {
	cmdArgs := []string{"build", label}
	cmdArgs = append(cmdArgs, envArgs...)

	cmd := exec.Command("bazel", cmdArgs...)
	cmd.Dir, _ = filepath.Abs(workspacePath)

	out, err := cmd.CombinedOutput()

	if verbose {
		log.Printf("Building %s in workspace %s", label, workspacePath)
		log.Println(string(out))
	}

	return err
}
