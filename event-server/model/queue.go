package model

import (
	pb "github.com/plarun/scheduler/event-server/data"
)

// QueryQueue stores and iterates list of job defintions to be processed
type QueryQueue struct {
	queries []*pb.Jil
	index   int
}

// NewQueryQueue returns new instance of QueryQueue
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
