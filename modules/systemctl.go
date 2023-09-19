package modules

import (
	"os/exec"

	"github.com/mircodezorzi/tachyon/pkg/actions"
)

type Systemctl struct {
	Service string
	Status  string
}

func (s *Systemctl) Do(step actions.Step, a interface{}) actions.Status {
	args, err := actions.ParseArgs[Systemctl](a)
	if err != nil {
		return actions.Status{Status: actions.Fail, Err: err}
	}

	cmd := exec.Command("systemctl", args.Status, args.Service)
	if step.Become {
		cmd = actions.AsUser(step.Playbook.BecomeUser, cmd)
	}
	b, err := cmd.CombinedOutput()
	if err != nil {
		return actions.Status{Status: actions.Fail, Output: b, Err: err}
	}
	return actions.Status{Status: actions.Ok}
}
