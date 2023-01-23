package task

type Event string

const (
	EventStart Event = "start"
	EventAbort Event = "abort"
	EventFroze Event = "froze"
	EventReset Event = "reset"
	EventGreen Event = "green"
)
