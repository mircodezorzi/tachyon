package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

type Status int

const (
	Ok      Status = 0
	Changed Status = 1
	Fail    Status = 2
)

func (s Status) String() string {
	switch s {
	case Ok:
		return "ok"
	case Changed:
		return "changed"
	case Fail:
		return "fail"
	}
	return "unimplemented"
}

type status struct {
	status Status
	output []byte
	delta  int
	err    error
}

func (s status) String() string {
	switch s.status {
	case Ok, Changed:
		return s.status.String()
	case Fail:
		return s.status.String() + ": " + string(s.output)
	}
	return "unimplemented"
}

type Quark struct {
	Name  string
	Steps []Step
}

type Packages []string

type Step struct {
	Name      string     `yaml:"name"`
	Become    *bool      `yaml:"become,omitempty"`
	User      *User      `yaml:"user,omitempty"`
	Yay       *Packages  `yaml:"yay,omitempty"`
	Pacman    *Packages  `yaml:"pacman,omitempty"`
	Command   *Command   `yaml:"command,omitempty"`
	Makepkg   *Makepkg   `yaml:"makepkg,omitempty"`
	Git       *Git       `yaml:"git,omitempty"`
	Systemctl *Systemctl `yaml:"systemctl,omitempty"`
	File      *File      `yaml:"file,omitempty"`

	playbook *Playbook
}

type File struct {
	Src string
	Path string
	State string
	Force bool
}

type Systemctl struct {
	Service string
	Status  string
}

type Git struct {
	Repo string
	Dest string
}

type Makepkg struct {
	Cwd string
}

type Command struct {
	Cmd string
	Cwd string
}

func compileTemplate(a string, b interface{}) []byte {
	tmpl := template.Must(template.New("rule").Parse(a))
	buf := &bytes.Buffer{}
	err := tmpl.Execute(buf, b)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func loadQuark(playbook Playbook, filepath, name string) (Quark, error) {
	var quark Quark

	b, err := os.ReadFile(path.Join(filepath, "quarks", name, "main.yaml"))
	if err != nil {
		return quark, err
	}

	playbook.Variables.(map[string]interface{})["role_path"] = path.Join(filepath, "quarks", name, "files")
	b = compileTemplate(string(b), playbook.Variables)

	if err = yaml.Unmarshal(b, &quark.Steps); err != nil {
		return quark, err
	}

	quark.Name = name

	return quark, nil
}

func main() {
	basepath := os.Args[1]
	p, err := loadPlaybook(basepath)
	if err != nil {
		panic(err)
	}

	fmt.Print("become password: ")
	reader := bufio.NewReader(os.Stdin)
	line, _, err := reader.ReadLine()
	p.becomePassword = string(line)

	for _, quark := range p.Quarks {
		q, err := loadQuark(p, basepath, quark)
		if err != nil {
			panic(err)
		}

		for _, step := range q.Steps {
			step.playbook = &p
			fmt.Printf("%s: %s\n", q.Name, step.Name)
			if step.Yay != nil {
				// step.updatePacman()
				fmt.Println(step.installYayPackages(*step.Yay...))
			}
			if step.Pacman != nil {
				// step.updatePacman()
				fmt.Println(step.installPacmanPackages(*step.Pacman...))
			}
			if step.User != nil {
				user := *step.User
				step.createUser(user.Name)
				step.userHome(user.Name, user.Home)
				step.userShell(user.Name, user.Shell)
				step.userGroups(user.Name, user.Groups...)
			}
			if step.Command != nil {
				command := *step.Command
				cs := strings.Split(command.Cmd, " ")
				cmd := exec.Command(cs[0], cs[1:]...)
				cmd.Dir = command.Cwd
				if step.Become != nil {
					cmd = asUser(step.playbook.BecomeUser, cmd)
				}
				b, err := cmd.CombinedOutput()
				_ = b
				_ = err
			}
			if step.Makepkg != nil {
				mkpkg := *step.Makepkg
				cmd := exec.Command("makepkg", "-si", "--noconfirm")
				cmd.Env = []string{"PACMAN=pacman -S"}
				cmd.Stdin = strings.NewReader(step.playbook.becomePassword)
				cmd.Dir = mkpkg.Cwd
				if step.Become != nil {
					cmd = asUser(step.playbook.BecomeUser, cmd)
				} else {
					panic("makepkg must become")
				}
				b, err := cmd.CombinedOutput()
				_ = b
				_ = err
			}
			if step.Git != nil {
				git := *step.Git
				cmd := exec.Command("git", "clone", git.Repo, git.Dest)
				if step.Become != nil {
					cmd = asUser(step.playbook.BecomeUser, cmd)
				}
				b, err := cmd.CombinedOutput()
				_ = b
				_ = err
			}
			if step.Systemctl != nil {
				systemctl := *step.Systemctl
				cmd := exec.Command("systemctl", systemctl.Status, systemctl.Service)
				b, err := cmd.CombinedOutput()
				_ = b
				_ = err
			}
			if step.File != nil {
				file := *step.File
				switch file.State {
				case "link":
					fmt.Println(os.Symlink(file.Src, file.Path))
				default:
					panic("unimplemented")
				}

			}
		}
	}
}
