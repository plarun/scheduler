package job

import (
	"fmt"
	"log"

	pb "github.com/plarun/scheduler/client/data"
	"github.com/plarun/scheduler/client/model"
	"golang.org/x/net/context"
)

type EventController struct {
	client pb.SendEventClient
}

// NewEventController returns new instance of EventController
func NewEventController(client pb.SendEventClient) *EventController {
	return &EventController{
		client: client,
	}
}

// Event sends the requested event to event server
func (controller EventController) Event(jobName string, event string) error {
	if event != "start" && event != "abort" && event != "freeze" && event != "reset" && event != "green" {
		return fmt.Errorf("invalid event type")
	}

	eventReq := &pb.SendEventReq{
		JobName:   jobName,
		EventType: model.EventTypeConv[event],
	}

	eventRes, err := controller.client.Event(context.Background(), eventReq)
	if err != nil {
		return err
	}

	if eventRes.EventChanged {
		log.Println("Event performed.")
	}

	return nil
}
