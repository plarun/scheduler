package service

type EventResult struct {
	IsSuccess bool
	Msg       error
}

func NewEventResult(succ bool, msg error) *EventResult {
	return &EventResult{
		IsSuccess: succ,
		Msg:       msg,
	}
}
