package elevate

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Reruns the current executable if we are not elevated. Call is final and will always exit
func EnsureElevatedFinal() {
	if !IsElevated() {
		RerunElevatedFinal()
	}
}

func IsElevated() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	return true
}

// Reruns the current executable. Call is final and will always exit
func RerunElevatedFinal() {
	err := RerunElevated()
	if err, ok := err.(*exec.ExitError); ok {
		os.Exit(err.ExitCode())
	}
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

// Reruns the current command with 1-to-1 arguments but elevated. May return the error from exec.Cmd's Run method
func RerunElevated() error {
	self, err := os.Executable()
	if err != nil {
		return errors.Wrap(err, "Could not get the executable")
	}
	cmdArgs := fmt.Sprintf(
		"Start-Process -Wait -Verb RunAs -WindowStyle Minimized -FilePath '%v' -ArgumentList '%v'",
		self, strings.Join(os.Args[1:], " "))
	command := exec.Command("powershell", "-Command", cmdArgs)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}
