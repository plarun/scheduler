package task

// Type represents send event type
type SendEvent string

const (
	// start event pushes the stable or waiting task to ready state for execution
	SendEventStart SendEvent = "start"
	// abort event terminates the running task or stops the triggered task from running
	SendEventAbort SendEvent = "abort"
	// freeze event sets the stable task to be ignored one, so wont be considered during scheduling
	SendEventFreeze SendEvent = "freeze"
	// reset event sets the state of stable task to default idle state
	SendEventReset SendEvent = "reset"
	// green event changes the state of stable task to success
	SendEventGreen SendEvent = "green"
	// red event changes the state of stable task to failure
	SendEventRed SendEvent = "red"
)

func (e SendEvent) IsAbort() bool {
	return e == SendEventAbort
}

func (e SendEvent) IsFreeze() bool {
	return e == SendEventFreeze
}

func (e SendEvent) IsGreen() bool {
	return e == SendEventGreen
}

func (e SendEvent) IsRed() bool {
	return e == SendEventRed
}

func (e SendEvent) IsReset() bool {
	return e == SendEventReset
}

func (e SendEvent) IsStart() bool {
	return e == SendEventStart
}
