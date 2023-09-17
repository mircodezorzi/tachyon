package main

import (
	"os/exec"
	"strings"
)

func (s *Step) updatePacman() error {
	cmd := exec.Command("pacman", "-Syy")
	err := cmd.Run()
	return err
}

func (s *Step) installYayPackages(packages ...string) error {
	arguments := []string{"-S", "--noconfirm", "--sudoflags", "-S"}
	arguments = append(arguments, packages...)
	cmd := exec.Command("yay", arguments...)
	if s.Become != nil {
		cmd = asUser(s.playbook.BecomeUser, cmd)
	}
	cmd.Stdin = strings.NewReader(s.playbook.becomePassword)
	err := cmd.Run()
	return err
}

func (s *Step) installPacmanPackages(packages ...string) error {
	arguments := []string{"-S", "--noconfirm"}
	arguments = append(arguments, packages...)
	cmd := exec.Command("pacman", arguments...)
	err := cmd.Run()
	return err
}
