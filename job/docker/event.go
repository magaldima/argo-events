package docker

import (
	"github.com/blackrock/axis/job"
	"fmt"
	"time"
)

type event struct {
	job.AbstractEvent
	docker *docker
	timestamp time.Time
}

func (e *event) GetBody() []byte {
	return []byte(fmt.Sprintf("%v", e.docker.message))
}

func (e *event) GetTimestamp() time.Time {
	return e.timestamp
}


func (e *event) GetID() string {
	return e.docker.AbstractSignal.GetID()
}

func (e *event) GetSource() string {
	return e.docker.message.From
}

func (e *event) GetSignal() job.Signal {
	return e.docker
}
