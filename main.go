package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"plugin"
	// "strings"
	"text/template"

	"github.com/mircodezorzi/tachyon/pkg/playbook"
	"gopkg.in/yaml.v3"

	"github.com/mircodezorzi/tachyon/modules"
	"github.com/mircodezorzi/tachyon/pkg/actions"
)

type Quark struct {
	Name  string
	Steps []S
}

type Packages []string

type S map[string]interface{}

type File struct {
	Src   string
	Path  string
	State string
	Force bool
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

func loadQuark(playbook playbook.Playbook, filepath, name string) (Quark, error) {
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

func registerAction(name string, action actions.Action) error {
	if _, ok := supportedActions[name]; ok {
		return errors.New("action already registered")
	}
	supportedActions[name] = action
	return nil
}

var supportedActions = map[string]actions.Action{}

func load(name string) actions.Action {
	p, err := plugin.Open(path.Join("plugins", name, name))
	if err != nil {
		panic(err)
	}

	a, err := p.Lookup("Action")
	if err != nil {
		panic(err)
	}

	return a.(actions.Action)
}

func main() {
	registerAction("git", &modules.Git{})
	registerAction("systemctl", &modules.Systemctl{})
	registerAction("yay", &modules.Yay{})
	registerAction("pacman", &modules.Pacman{})

	// registerAction("echo", load("echo"))

	basepath := os.Args[1]
	p, err := playbook.LoadPlaybook(basepath)
	if err != nil {
		panic(err)
	}

	fmt.Print("become password: ")
	reader := bufio.NewReader(os.Stdin)
	line, _, err := reader.ReadLine()
	p.BecomePassword = string(line)

	for _, quark := range p.Quarks {
		q, err := loadQuark(p, basepath, quark)
		if err != nil {
			panic(err)
		}

		for _, s := range q.Steps {
			var action actions.Action
			var args interface{}
			name, ok := s["name"].(string)
			if !ok {
			}
			become, ok := s["become"].(bool)
			for k, v := range s {
				a, ok := supportedActions[k]
				if !ok {
					continue
				}
				action = a
				args = v
			}
			if action == nil {
				continue
			}
			step := actions.Step{
				Name:       name,
				Become:     become,
				BecomeUser: "root",
				Action:     action,
				Args:       args,
				Playbook:   &p,
			}

			fmt.Printf("%s: %s\n", q.Name, step.Name)
			step.Action.Do(step, args)
			/*
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
				if step.File != nil {
					file := *step.File
					switch file.State {
					case "link":
						fmt.Println(os.Symlink(file.Src, file.Path))
					default:
						panic("unimplemented")
					}
				}
			*/
		}
	}
}
