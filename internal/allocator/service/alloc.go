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
	stop := make(chan error)
	var err error

	for err != nil {
		select {
		case <-ticker.C:
			err = a.poller.Stage()
			stop <- err
		case <-stop:
			ticker.Stop()
		}
	}

	if err != nil {
		ch <- fmt.Errorf("stage: %w", err)
	}
}

func (a *Allocator) poll(ch chan (error)) {
	ticker := time.NewTicker(a.pollCycle)
	stop := make(chan error)
	var err error

	for err != nil {
		select {
		case <-ticker.C:
			err = a.poller.Poll()
			stop <- err
		case <-stop:
			ticker.Stop()
		}
	}

	if err != nil {
		ch <- fmt.Errorf("poll: %w", err)
	}
}

func (a *Allocator) split(ch chan (error)) {
	ticker := time.NewTicker(a.splitCycle)
	stop := make(chan error)
	var err error

	for err != nil {
		select {
		case <-ticker.C:
			err = a.poller.Poll()
			stop <- err
		case <-stop:
			ticker.Stop()
		}
	}

	if err != nil {
		ch <- fmt.Errorf("split: %w", err)
	}
}
