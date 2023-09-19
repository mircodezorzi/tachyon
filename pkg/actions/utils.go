package actions

import (
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
)

func AsUser(name string, cmd *exec.Cmd) *exec.Cmd {
	u, _ := user.Lookup(name)
	uid, _ := strconv.Atoi(u.Uid)
	gid, _ := strconv.Atoi(u.Gid)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}
	return cmd
}
