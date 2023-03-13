package task

type Event string

const (
	// EventStart is an event to start the task execution
	EventStart Event = "start"
	// EventAbort is an event to abort the running task
	EventAbort Event = "abort"
	// EventFroze is an event to freeze the task
	EventFroze Event = "froze"
	// EventReset is an event to change the status of task to default state
	EventReset Event = "reset"
	// EventGreen is an event to change the status of task to success
	EventGreen Event = "green"
)
