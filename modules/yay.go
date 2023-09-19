package modules

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/mircodezorzi/tachyon/pkg/actions"
)

type Yay struct {
}

func (y *Yay) findYayDelta(step actions.Step, packages ...string) (int, error) {
	arguments := append([]string{"-Qu"}, packages...)
	cmd := exec.Command("yay", arguments...)
	if step.Become {
		cmd = actions.AsUser(step.Playbook.BecomeUser, cmd)
		cmd.Stdin = strings.NewReader(step.Playbook.BecomePassword)
	}
	b, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}
	if strings.Contains(string(b), "was not found") {
		return 0, errors.New(string(b))
	}
	return strings.Count(string(b), "\n"), nil
}

func (y *Yay) Do(step actions.Step, a interface{}) actions.Status {
	args, err := actions.ParseArgs[[]string](a)
	delta, err := y.findYayDelta(step, args...)
	if err != nil {
		return actions.Status{Status: actions.Fail, Err: err}
	}
	if delta == 0 {
		return actions.Status{Status: actions.Ok}
	}
	arguments := append([]string{"-S", "--noconfirm", "--sudoflags", "-S"}, args...)
	cmd := exec.Command("yay", arguments...)
	if step.Become {
		cmd = actions.AsUser(step.Playbook.BecomeUser, cmd)
		cmd.Stdin = strings.NewReader(step.Playbook.BecomePassword)
	}
	b, err := cmd.CombinedOutput()
	if err != nil {
		return actions.Status{Status: actions.Fail, Output: b, Err: err}
	}
	return actions.Status{Status: actions.Changed}
}
