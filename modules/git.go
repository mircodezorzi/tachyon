package modules

import (
	"os/exec"

	"github.com/mircodezorzi/tachyon/pkg/actions"
)

type Git struct {
	Repo string `yaml:"repo,omitempty"`
	Dest string `yaml:"dest,omitempty"`
}

func (g *Git) Do(step actions.Step, a interface{}) actions.Status {
	args, err := actions.ParseArgs[Git](a)
	if err != nil {
		return actions.Status{Status: actions.Fail, Err: err}
	}

	cmd := exec.Command("git", "clone", args.Repo, args.Dest)
	if step.Become {
		cmd = actions.AsUser(step.Playbook.BecomeUser, cmd)
	}
	b, err := cmd.CombinedOutput()
	if err != nil {
		return actions.Status{Status: actions.Fail, Output: b, Err: err}
	}
	return actions.Status{Status: actions.Ok}
}
