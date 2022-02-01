package job

import (
	"fmt"
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
	if event != "start" && event != "abort" && event != "freeze" && event != "reset" && event != "green" {
		return fmt.Errorf("invalid event type")
	}

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
