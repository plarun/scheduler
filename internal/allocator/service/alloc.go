package service

import (
	"fmt"
	"time"
)

const (
	SPLIT_CYCLE time.Duration = 5 * time.Second
	STAGE_CYCLE time.Duration = 5 * time.Second
	POLL_CYCLE  time.Duration = 10 * time.Second
)

type Allocator struct {
	splitCycle time.Duration
	stageCycle time.Duration
	pollCycle  time.Duration
	poller     *TaskPoller
	splitter   *TaskSplitter
}

func NewAllocator() *Allocator {
	return &Allocator{
		splitCycle: SPLIT_CYCLE,
		stageCycle: STAGE_CYCLE,
		pollCycle:  POLL_CYCLE,
		poller:     NewTaskPoller(POLL_CYCLE),
		splitter:   NewTaskSplitter(),
	}
}

func (a *Allocator) Start() error {
	fail := make(chan error)

	go a.stage(fail)
	go a.poll(fail)
	go a.split(fail)

	if err := <-fail; err != nil {
		return fmt.Errorf("Start: %w", err)
	}
	return nil
}

func (a *Allocator) stage(ch chan (error)) {
	ticker := time.NewTicker(a.stageCycle)

	for range ticker.C {
		if err := a.poller.Stage(); err != nil {
			ch <- err
			break
		}
	}
}

func (a *Allocator) poll(ch chan (error)) {
	ticker := time.NewTicker(a.pollCycle)

	for range ticker.C {
		if err := a.poller.Poll(); err != nil {
			ch <- err
			break
		}
	}
}

func (a *Allocator) split(ch chan (error)) {
	ticker := time.NewTicker(a.splitCycle)

	for range ticker.C {
		if err := a.splitter.Split(); err != nil {
			ch <- err
			break
		}
	}
}
