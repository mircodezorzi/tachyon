package modules

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/mircodezorzi/tachyon/pkg/actions"
)

type User struct {
	Name   string   `yaml:"name"`
	Home   string   `yaml:"home"`
	Groups []string `yaml:"groups,omitempty"`
	Shell  string   `yaml:"shell,omitempty"`
}

func (u *User) createUser(user string) actions.Status {
	cmd := exec.Command("useradd", user)
	b, err := cmd.CombinedOutput()
	if strings.Contains(string(b), "already exists") {
		return actions.Status{Status: actions.Ok, Err: err}
	}
	return actions.Status{Status: actions.Ok}
}

func (u *User) userHome(user, home string) actions.Status {
	cmd := exec.Command("usermod", "-d", home, user)
	b, err := cmd.CombinedOutput()
	if strings.Contains(string(b), "no changes") {
		return actions.Status{Status: actions.Ok}
	}
	if strings.Contains(string(b), "does not exist") {
		return actions.Status{Status: actions.Fail, Err: err}
	}
	return actions.Status{Status: actions.Changed}
}

func (u *User) userShell(user, shell string) actions.Status {
	cmd := exec.Command("chsh", user, "-s", user)
	b, err := cmd.CombinedOutput()
	if strings.Contains(string(b), "does not exist") {
		return actions.Status{Status: actions.Fail, Output: b, Err: err}
	}
	return actions.Status{Status: actions.Ok}
}

func (u *User) userGroups(user string, groups ...string) actions.Status {
	g := strings.Join(groups, ",")
	cmd := exec.Command("usermod", "-a", "-G", g, user)
	b, err := cmd.CombinedOutput()
	if strings.Contains(string(b), "does not exist") {
		return actions.Status{Status: actions.Fail, Output: b, Err: err}
	}
	return actions.Status{Status: actions.Ok}
}

func (u *User) Do(step actions.Step, a interface{}) actions.Status {
	args, err := actions.ParseArgs[User](a)
	if err != nil {
		return actions.Status{Status: actions.Fail, Err: err}
	}

	fmt.Println(args)
	return actions.Status{Status: actions.Ok}
}
