package main

import (
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type Playbook struct {
	Quarks    []string    `yaml:"quarks"`
	Variables interface{} `yaml:"variables"`

	BecomeUser string `yaml:"become_user"`

	becomePassword string
}

func loadPlaybook(filepath string) (Playbook, error) {
	var playbook Playbook

	b, err := os.ReadFile(path.Join(filepath, "playbook.yaml"))
	if err != nil {
		return playbook, err
	}

	if err = yaml.Unmarshal(b, &playbook); err != nil {
		return playbook, err
	}

	return playbook, nil
}
