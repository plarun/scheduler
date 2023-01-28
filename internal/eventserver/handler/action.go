package handler

import (
	"errors"
	"fmt"
	"log"

	"github.com/plarun/scheduler/config"
	es "github.com/plarun/scheduler/internal/eventserver/service"
	"github.com/plarun/scheduler/proto"
	"golang.org/x/net/context"
)

// ParsedActionService is a grpc server for routing the parsed entity
// actions from the client to validator service for validation
type ParsedActionService struct {
	proto.UnimplementedParsedActionServiceServer
}

func NewParsedActionService() *ParsedActionService {
	return &ParsedActionService{}
}

func (s ParsedActionService) Submit(ctx context.Context, req *proto.ParsedEntitiesRequest) (*proto.EntityActionResponse, error) {
	res := &proto.EntityActionResponse{
		Status: &proto.ActionStatus{
			Success:  false,
			Errors:   make([]string, 0),
			Warnings: make([]string, 0),
		},
	}

	valRes, err := routeValidationReq(req)

	// grpc error
	if err != nil {
		for errors.Unwrap(err) != nil {
			err = errors.Unwrap(err)
		}
		log.Printf("Error: Submit: %v", err)
		return res, err
	}

	res.Status = valRes.Status

	if req.OnlyValidate {
		return res, nil
	}

	// process the validated task entities' actions
	if err := es.ProcessTaskActions(ctx, valRes.Tasks); err != nil {
		for errors.Unwrap(err) != nil {
			err = errors.Unwrap(err)
		}

		res.Status.Success = false
		res.Status.Errors = append(res.Status.Errors, err.Error())
		return res, nil
	}

	log.Println("Task processing request handled")
	return res, nil
}

func routeValidationReq(req *proto.ParsedEntitiesRequest) (*proto.ValidatedEntitiesResponse, error) {
	addr := fmt.Sprintf(":%d", config.GetAppConfig().Service.Validator.Port)
	valConn := NewValidationGrpcConnection(addr, req)

	log.Println("Routing the validation request")

	res := &proto.ValidatedEntitiesResponse{
		Tasks: make([]*proto.ValidatedTaskEntity, 0),
		Status: &proto.ActionStatus{
			Success:  false,
			Errors:   make([]string, 0),
			Warnings: make([]string, 0),
		},
	}

	if err := valConn.Connect(); err != nil {
		return res, err
	}

	valRes, err := valConn.Request()
	if err != nil {
		return res, err
	}

	ent, ok := valRes.(*proto.ValidatedEntitiesResponse)
	if !ok {
		return ent, fmt.Errorf("internal err")
	}

	if err := valConn.Close(); err != nil {
		return ent, err
	}

	log.Println("Task definitions validated by validator service")
	return ent, nil
}
