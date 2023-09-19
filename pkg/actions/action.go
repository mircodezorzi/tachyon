package actions

import (
	"gopkg.in/yaml.v3"
)

type Action interface {
	Do(Step, interface{}) Status
}

func ParseArgs[T any](a interface{}) (T, error) {
	var args T
	b, err := yaml.Marshal(a)
	if err != nil {
		return args, err
	}
	if err := yaml.Unmarshal(b, &args); err != nil {
		return args, err
	}
	return args, nil
}
