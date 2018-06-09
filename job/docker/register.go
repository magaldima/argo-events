package docker

import (
	"github.com/blackrock/axis/job"
	"go.uber.org/zap"
	"github.com/blackrock/axis/pkg/apis/sensor/v1alpha1"
)

type factory struct{}

func (f *factory) Create(abstract job.AbstractSignal) job.Signal {
	abstract.Log.Info("creating signal", zap.String("type", abstract.Docker.Type))
	return &docker{
		AbstractSignal: abstract,
	}
}

// Docker will be added to the executor session
func Docker(es *job.ExecutorSession) {
	es.AddFactory(v1alpha1.SignalTypeDocker, &factory{})
}
