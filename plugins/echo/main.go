package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/mircodezorzi/tachyon/pkg/actions"
)

func init() {

}

type Echo struct {
	Message string `yaml:"echo"`
}

func (e *Echo) Do(step actions.Step, a interface{}) actions.Status {
	fmt.Println(a)
	args, err := actions.ParseArgs[Echo](a)
	if err != nil {
		return actions.Status{Status: actions.Fail, Err: err}
	}

	cmd := exec.Command("echo", args.Message)
	cmd.Stdout = os.Stdout
	if err := cmd.Start(); err != nil {
		return actions.Status{Status: actions.Fail, Err: err}
	}

	return actions.Status{Status: actions.Ok}
}

var Action = Echo{}
