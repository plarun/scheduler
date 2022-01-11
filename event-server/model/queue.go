package model

import (
	pb "github.com/plarun/scheduler/event-server/data"
)

type QueryQueue struct {
	queries []*pb.Jil
	index   int
}

func NewQueryQueue() *QueryQueue {
	return &QueryQueue{
		queries: make([]*pb.Jil, 0),
		index:   0,
	}
}

func (que *QueryQueue) Add(queryData *pb.Jil) error {
	que.queries = append(que.queries, queryData)
	return nil
}

func (que *QueryQueue) HasNext() bool {
	return que.index < len(que.queries)
}

func (que *QueryQueue) Next() *pb.Jil {
	if que.HasNext() {
		query := que.queries[que.index]
		que.index++
		return query
	}
	return nil
}
