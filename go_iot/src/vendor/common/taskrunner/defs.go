package taskrunner

const (
	READY_TO_DISPATCH = "d"
	READY_TO_EXECUTE = "e"
	CLOSE = "c"
)

type ControlChan chan string

type DataChan chan interface{}

type Fn func(dc DataChan) error
