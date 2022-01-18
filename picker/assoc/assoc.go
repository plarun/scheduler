package assoc

import (
	pb "github.com/plarun/scheduler/picker/data"
)

var dependency *Dependent = nil

type Dependent struct {
	link map[string][]string
}

func NewDependent() *Dependent {
	if dependency == nil {
		dependency = &Dependent{
			link: make(map[string][]string),
		}
	}
	return dependency
}

func (dep *Dependent) Add(job *pb.ReadyJob) error {
	for _, preceder := range job.GetPreceders() {
		if _, ok := dep.link[preceder]; !ok {
			dep.link[preceder] = make([]string, 0)
		}
		dep.link[preceder] = append(dep.link[preceder], job.GetJobName())
	}

	return nil
}
