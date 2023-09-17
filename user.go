package main

import (
	"os/exec"
	"strings"
)

type User struct {
	Name   string   `yaml:"name"`
	Home   string   `yaml:"home"`
	Groups []string `yaml:"groups,omitempty"`
	Shell  string   `yaml:"shell,omitempty"`
}

func (s *Step) createUser(user string) status {
	cmd := exec.Command("useradd", user)
	b, err := cmd.CombinedOutput()
	if strings.Contains(string(b), "already exists") {
		return status{status: Ok, err: err}
	}
	return status{status: Ok}
}

func (s *Step) userHome(user, home string) status {
	cmd := exec.Command("usermod", "-d", home, user)
	b, err := cmd.CombinedOutput()
	if strings.Contains(string(b), "no changes") {
		return status{status: Ok}
	}
	if strings.Contains(string(b), "does not exist") {
		return status{status: Fail, err: err}
	}
	return status{status: Changed}
}

func (s *Step) userShell(user, shell string) status {
	cmd := exec.Command("chsh", user, "-s", user)
	b, err := cmd.CombinedOutput()
	if strings.Contains(string(b), "does not exist") {
		return status{status: Fail, output: b, err: err}
	}
	return status{status: Ok}
}

func (s *Step) userGroups(user string, groups ...string) status {
	g := strings.Join(groups, ",")
	cmd := exec.Command("usermod", "-a", "-G", g, user)
	b, err := cmd.CombinedOutput()
	if strings.Contains(string(b), "does not exist") {
		return status{status: Fail, output: b, err: err}
	}
	return status{status: Ok}
}
