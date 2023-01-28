package handler

import (
	"errors"
	"log"

	"github.com/plarun/scheduler/internal/validator/service"
	"github.com/plarun/scheduler/proto"
	"golang.org/x/net/context"
)

// taskActionService serves the requests from client
// then it will route the request to validator for
// validating the
type taskActionService struct {
	proto.UnimplementedValidatedActionServiceServer
}

func NewTaskValidationService() taskActionService {
	return taskActionService{}
}

func (svc taskActionService) Route(ctx context.Context, req *proto.ParsedEntitiesRequest) (*proto.ValidatedEntitiesResponse, error) {

	tsks := req.Tasks
	n := len(tsks)

	res := &proto.ValidatedEntitiesResponse{
		Tasks: make([]*proto.ValidatedTaskEntity, n),
		Status: &proto.ActionStatus{
			Success:  true,
			Errors:   make([]string, 0),
			Warnings: make([]string, 0),
		},
	}

	log.Println("New validation request")

	for i := 0; i < n; i++ {
		tsk, err := service.ValidateTaskAction(tsks[i])
		if err != nil {
			log.Printf("validation failure: %v", err)
			for errors.Unwrap(err) != nil {
				err = errors.Unwrap(err)
			}
			res.Status.Errors = append(res.Status.Errors, err.Error())
			res.Status.Success = false
		}
		res.Tasks[i] = tsk
	}

	if !res.Status.Success {
		return res, nil
	}

	chk := service.NewChecker()
	if err := chk.CheckExistance(res.Tasks); err != nil {
		log.Printf("Existance failure: %v", err)
		for errors.Unwrap(err) != nil {
			err = errors.Unwrap(err)
		}
		res.Status.Errors = append(res.Status.Errors, err.Error())
		res.Status.Success = false
		return res, err
	}

	log.Println("validation completed")
	return res, nil
}
