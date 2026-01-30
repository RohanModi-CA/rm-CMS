package misc

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type LogLevel int

const (
	LogSilent LogLevel = iota
	LogErrors
	LogSome
	LogAll
)

func LogF(level LogLevel, trigger_level LogLevel, format string, a ...interface{}) {
	if level >= trigger_level {
		fmt.Printf(format, a...)
	}
}

type ExecOptions struct {
	CommandList []string
	LogCommand  bool
	Dir         string
	DontPanic   bool
	CS          *ConversionState
	Stdin       string
}

func RunExecCommand(eo ExecOptions) {
	/* A basic wrapper function for exec.Command, to shorten function calls.
	LogSilent logs nothing, LogErrors/LogSome logs stderr, LogAll stdout.
	Panics on errors unless DontPanic is true. */

	var full_command string
	full_command = strings.Join(eo.CommandList, " ")

	cmd := exec.Command(eo.CommandList[0], eo.CommandList[1:]...)

	if eo.Dir != "" {
		cmd.Dir = eo.Dir
	}

	var stderr_buf bytes.Buffer
	var stdout_buf bytes.Buffer

	cmd.Stderr = &stderr_buf
	cmd.Stdout = &stdout_buf

	err := cmd.Run()

	if eo.LogCommand {
		fmt.Printf("\nRunning: %s ", full_command)
	}

	if err != nil && !eo.DontPanic {
		fmt.Printf("\n\nSTDOUT:\n\n%s\n\n", stdout_buf.String())
		fmt.Printf("\n\nSTDERR:\n\n%s\n\n", stderr_buf.String())
		error_message := fmt.Sprintf("Error running command: \n%s\n, received: %s\n\n", full_command, err)

		panic(error_message)
	}

	if stderr_buf.String() != "" {
		LogF(eo.CS.Level, LogErrors, "STDERR: \n %s\n", stderr_buf.String())
	}
	if stdout_buf.String() != "" {
		LogF(eo.CS.Level, LogAll, "STDOUT: \n %s\n", stdout_buf.String())
	}
}

func RunExecCommandOut(eo ExecOptions) string {
	/* RETURNS THE CONTENTS OF STDOUT.
	A basic wrapper function for exec.Command, to shorten function calls.
	LogSilent logs nothing, LogErrors/LogSome logs stderr, LogAll stdout.
	Panics on errors unless DontPanic is true. */

	var full_command string
	full_command = strings.Join(eo.CommandList, " ")

	cmd := exec.Command(eo.CommandList[0], eo.CommandList[1:]...)

	if eo.Dir != "" {
		cmd.Dir = eo.Dir
	}

	if eo.Stdin != "" {
		cmd.Stdin = strings.NewReader(eo.Stdin)
	}

	var stderr_buf bytes.Buffer
	var stdout_buf bytes.Buffer

	cmd.Stderr = &stderr_buf
	cmd.Stdout = &stdout_buf

	err := cmd.Run()

	if eo.LogCommand {
		fmt.Printf("\nRunning: %s ", full_command)
	}

	if err != nil && !eo.DontPanic {
		fmt.Printf("\n\nSTDOUT:\n\n%s\n\n", stdout_buf.String())
		fmt.Printf("\n\nSTDERR:\n\n%s\n\n", stderr_buf.String())
		error_message := fmt.Sprintf("Error running command: \n%s\n, received: %s\n\n", full_command, err)

		panic(error_message)
	}

	if stderr_buf.String() != "" {
		LogF(eo.CS.Level, LogErrors, "STDERR: \n %s\n", stderr_buf.String())
	}

	// This will be "" if nothing.
	return stdout_buf.String()
}

func ErrorHandlePanic(err error) {
	if err != nil {
		print("Panicking.")
		panic(err)
	}
}

func GetRMCMSVersion(cs *ConversionState) string {
	commandlist := []string{"git", "describe", "--tags", "--always"}
	eo := ExecOptions{
		CommandList: commandlist,
		LogCommand:  false,
		Dir:         "",
		DontPanic:   false,
		CS:          cs,
	}
	version := RunExecCommandOut(eo)
	return strings.TrimSpace(version)
}
