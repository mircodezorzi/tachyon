package main

import (
	"errors"
	"os/exec"
	"strings"
)

func (s *Step) findYayDelta(packages ...string) (int, error) {
	arguments := append([]string{"-Qu"}, packages...)
	cmd := exec.Command("yay", arguments...)
	if s.Become != nil {
		cmd = asUser(s.playbook.BecomeUser, cmd)
		cmd.Stdin = strings.NewReader(s.playbook.becomePassword)
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

func (s *Step) installYayPackages(packages ...string) status {
	delta, err := s.findYayDelta(packages...)
	if err != nil {
		return status{status: Fail, err: err}
	}
	if delta == 0 {
		return status{status: Ok}
	}
	arguments := append([]string{"-S", "--noconfirm", "--sudoflags", "-S"}, packages...)
	cmd := exec.Command("yay", arguments...)
	if s.Become != nil {
		cmd = asUser(s.playbook.BecomeUser, cmd)
		cmd.Stdin = strings.NewReader(s.playbook.becomePassword)
	}
	b, err := cmd.CombinedOutput()
	if err != nil {
		return status{status: Fail, output: b, err: err}
	}
	return status{status: Changed}
}
