package service

import (
	"context"
	"fmt"

	"github.com/plarun/scheduler/config"
	"github.com/plarun/scheduler/proto"
)

func AwakeWaitingDependentTasks(ctx context.Context, id int64) error {
	// todo: call rpc to allocator for awake
	addr := fmt.Sprintf(":%d", config.GetAppConfig().Service.Allocator.Port)
	req := &proto.DependentTaskAwakeRequest{}

	conn := NewWaitTaskGrpcConnection(addr, req)

	if err := conn.Connect(); err != nil {
		return fmt.Errorf("AwakeWaitingDependentTasks: %w", err)
	}

	r, err := conn.Request()
	if err != nil {
		return fmt.Errorf("AwakeWaitingDependentTasks: %w", err)
	}

	var ok bool
	if _, ok = r.(*proto.DependentTaskAwakeResponse); !ok {
		panic("invalid type")
	}

	return nil
}
