package utils

import (
	"bytes"
	"fmt"
	"os/exec"
)

// runCommand executes a shell command and captures its output.
func RunCommand(cmd string) error {
	fmt.Printf("Executing: %s\n", cmd)

	command := exec.Command("sh", "-c", cmd)
	var out bytes.Buffer
	command.Stdout = &out
	command.Stderr = &out
	err := command.Run()
	if err != nil {
		fmt.Printf("Error: %s\n", out.String())
		return err
	}
	fmt.Printf("Output: %s\n", out.String())
	return nil
}
