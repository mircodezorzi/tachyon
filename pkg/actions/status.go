package actions

type StatusCode int

const (
	Ok      StatusCode = 0
	Changed StatusCode = 1
	Fail    StatusCode = 2
)

func (s StatusCode) String() string {
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

type Status struct {
	Status StatusCode
	Output []byte
	Delta  int
	Err    error
}

func (s Status) String() string {
	switch s.Status {
	case Ok, Changed:
		return s.Status.String()
	case Fail:
		return s.Status.String() + ": " + string(s.Output)
	}
	return "unimplemented"
}
