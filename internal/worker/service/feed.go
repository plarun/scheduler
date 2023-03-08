package service

import (
	"fmt"

	"github.com/plarun/scheduler/config"
	"github.com/plarun/scheduler/pkg/grpcconn"
	"github.com/plarun/scheduler/proto"
)

type TaskFeed struct {
	conn       grpcconn.GrpcConnecter
	workerPool *WorkerPool
}

func NewTaskFeed(pool *WorkerPool) *TaskFeed {
	addr := fmt.Sprintf(":%d", config.GetAppConfig().Service.EventServer.Port)
	req := &proto.ReadyTasksPullRequest{}

	conn := NewReadyTasksGrpcConnection(addr, req)

	return &TaskFeed{
		conn:       conn,
		workerPool: pool,
	}
}

func (t *TaskFeed) Feed() error {
	if err := t.conn.Connect(); err != nil {
		return fmt.Errorf("Feed: %w", err)
	}

	r, err := t.conn.Request()
	if err != nil {
		return fmt.Errorf("Feed: %w", err)
	}

	var res *proto.ReadyTasksPullResponse
	var ok bool
	if _, ok = r.(*proto.ReadyTasksPullResponse); !ok {
		panic("invalid type")
	}

	// ids of tasks which are ready for execution
	tasks := res.TaskIds

	for _, taskId := range tasks {
		cmd, fout, ferr, err := getTaskInfo(taskId)
		if err != nil {
			return fmt.Errorf("Feed: %w", err)
		}
		ex := NewExecutable(taskId, cmd, fout, ferr)
		t.workerPool.Add(ex)
	}

	// feed tasks into worker pool for execution

	if err := t.conn.Close(); err != nil {
		return fmt.Errorf("Feed: %w", err)
	}
	return nil
}

func getTaskInfo(id int64) (string, string, string, error) {
	addr := fmt.Sprintf(":%d", config.GetAppConfig().Service.EventServer.Port)
	req := &proto.TaskInfoRequest{
		TaskId: id,
	}

	conn := NewTaskInfoGrpcConnection(addr, req)
	if err := conn.Connect(); err != nil {
		return "", "", "", fmt.Errorf("getTaskInfo: %w", err)
	}

	r, err := conn.Request()
	if err != nil {
		return "", "", "", fmt.Errorf("getTaskInfo: %w", err)
	}

	var res *proto.TaskInfoResponse
	var ok bool
	if res, ok = r.(*proto.TaskInfoResponse); !ok {
		panic("invalid type")
	}

	cmd, fout, ferr := res.Command, res.OutLogFile, res.ErrLogFile

	if err := conn.Close(); err != nil {
		return cmd, fout, ferr, err
	}
	return cmd, fout, ferr, nil
}
