package job

import (
	"log"
	"time"

	pb "github.com/plarun/scheduler/client/data"
	"github.com/plarun/scheduler/client/model"
	"golang.org/x/net/context"
)

type EventController struct {
	client pb.SendEventClient
}

func NewEventController(client pb.SendEventClient) *EventController {
	return &EventController{
		client: client,
	}
}

func (controller EventController) Event(jobName string, event string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	eventReq := &pb.SendEventReq{
		JobName:   jobName,
		EventType: model.EventTypeConv[event],
	}

	eventRes, err := controller.client.Event(ctx, eventReq)
	if err != nil {
		return err
	}

	if eventRes.EventChanged {
		log.Println("Event performed.")
	}

	return nil
}
