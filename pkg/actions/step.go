package actions

import (
	"os/exec"
	"strings"

	"github.com/mircodezorzi/tachyon/pkg/playbook"
)

type Step struct {
	Name       string
	Become     bool
	BecomeUser string
	Action     Action
	Args       interface{}
	Playbook   *playbook.Playbook
}

func (s *Step) Cmd(name string, arg ...string) *exec.Cmd {
	if s.Become {
		args := append([]string{"-S"}, name)
		args = append(args, arg...)
		cmd := exec.Command("sudo", args...)
		cmd.Stdin = strings.NewReader(s.Playbook.BecomePassword)
		// u, _ := user.Lookup(s.BecomeUser)
		// uid, _ := strconv.Atoi(u.Uid)
		// gid, _ := strconv.Atoi(u.Gid)
		// cmd.SysProcAttr = &syscall.SysProcAttr{}
		// cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}
		return cmd
	}

	cmd := exec.Command(name, arg...)
	return cmd
}
