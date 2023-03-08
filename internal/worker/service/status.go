package service

type Status struct {
	taskId   int
	exitCode int
}

func NewStatus(id, code int) Status {
	return Status{
		taskId:   id,
		exitCode: code,
	}
}
