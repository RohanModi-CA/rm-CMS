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
	if (level >= trigger_level) {
		fmt.Printf(format, a...)
	}
}


type ExecOptions struct {
	CommandList []string
	LogCommand bool
	Dir string
	DontPanic bool
	CS *ConversionState
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


	if (err != nil && !eo.DontPanic) {
		fmt.Printf("\n\nSTDOUT:\n\n%s\n\n", stdout_buf.String())
		fmt.Printf("\n\nSTDERR:\n\n%s\n\n", stderr_buf.String())
		error_message := fmt.Sprintf("Error running command: \n%s\n, received: %s\n\n", full_command, err)

		panic(error_message)
	}

	LogF(eo.CS.Level, LogErrors, "STDERR: \n %s\n", stderr_buf.String())
	LogF(eo.CS.Level, LogAll, "STDOUT: \n %s\n", stdout_buf.String())
}
