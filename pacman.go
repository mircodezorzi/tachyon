package main

import (
	"errors"
	"os/exec"
	"strings"
)

func (s *Step) updatePacman() error {
	cmd := exec.Command("pacman", "-Syy")
	err := cmd.Run()
	return err
}

func (s *Step) findPacmanDelta(packages ...string) (int, error) {
	arguments := append([]string{"-Qu"}, packages...)
	cmd := exec.Command("pacman", arguments...)
	if s.Become != nil {
		cmd = asUser(s.playbook.BecomeUser, cmd)
		cmd.Stdin = strings.NewReader(s.playbook.becomePassword)
	}
	b, err := cmd.CombinedOutput()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok && exiterr.ExitCode() != 1 {
			return 0, err
		}
	} 
	if strings.Contains(string(b), "was not found") {
		return 0, errors.New(string(b))
	}
	return strings.Count(string(b), "\n"), nil
}

func (s *Step) installPacmanPackages(packages ...string) status {
	delta, err := s.findPacmanDelta(packages...);  
	if err != nil {
		return status{status: Fail, err: err}
	}
	if delta == 0 {
		return status{status: Ok}
	}
	arguments := append([]string{"-S", "--noconfirm"}, packages...)
	cmd := exec.Command("pacman", arguments...)
	b, err := cmd.CombinedOutput()
	if err != nil {
		return status{status: Fail, output: b, err: err}
	}
	return status{status: Ok, delta: delta}
}
