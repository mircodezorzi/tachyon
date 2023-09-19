package modules

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mircodezorzi/tachyon/pkg/actions"
)

type Pacman struct {
}

func (p *Pacman) updatePacman() error {
	cmd := exec.Command("pacman", "-Syy")
	err := cmd.Run()
	return err
}

func (p *Pacman) findPacmanDelta(step actions.Step, packages ...string) (int, error) {
	arguments := append([]string{"-Qu"}, packages...)
	cmd := step.Cmd("pacman", arguments...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) && exit.ExitCode() != 1 {
			return 0, err
		}
	}
	if strings.Contains(string(b), "was not found") {
		return 0, errors.New(string(b))
	}
	return strings.Count(string(b), "\n"), nil
}

func (p *Pacman) Do(step actions.Step, a interface{}) actions.Status {
	args, err := actions.ParseArgs[[]string](a)
	if err != nil {
		return actions.Status{Status: actions.Fail, Err: err}
	}

	delta, err := p.findPacmanDelta(step, args...)
	if err != nil {
		return actions.Status{Status: actions.Fail, Err: err}
	}
	if delta == 0 {
		return actions.Status{Status: actions.Ok}
	}
	arguments := append([]string{"-S", "--noconfirm"}, args...)
	cmd := step.Cmd("pacman", arguments...)
	b, err := cmd.CombinedOutput()
	fmt.Println(string(b))
	if err != nil {
		return actions.Status{Status: actions.Fail, Output: b, Err: err}
	}
	return actions.Status{Status: actions.Ok, Delta: delta}
}
