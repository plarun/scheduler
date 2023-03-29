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
