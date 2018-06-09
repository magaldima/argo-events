package docker

import (
	"github.com/blackrock/axis/job"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"fmt"
	"context"
	dockerEvents "github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
	"go.uber.org/zap"
	"sync"
	"time"
)

type docker struct {
	job.AbstractSignal
	stop chan struct{}
	message dockerEvents.Message
	events chan job.Event
	wg sync.WaitGroup
}

func (d *docker) validateSignal() error{
	// TODO: Validate filter and event type
	return nil
}

func (d *docker) Start(events chan job.Event) error {
	eventFilters := types.EventsOptions{}
	if len(d.Docker.Filters) > 0 {
		filter := filters.NewArgs()
		for filterKey, filterValue := range d.Docker.Filters {
			filter.Add(filterKey, filterValue)
		}
	}

	dockerClient, err := client.NewEnvClient()

	if err != nil {
		return fmt.Errorf("unable to initialize docker client. Cause %v", err)
	}

	messages, errs := dockerClient.Events(context.Background(), eventFilters)
	d.wg.Add(1)
	go d.handleEvents(messages, errs)
	return nil
}

func (d *docker) handleEvents(messages <-chan dockerEvents.Message, errs <-chan error) {
	event := &event{
		docker: d,
		timestamp: time.Now().UTC(),
	}
	for {
		select {
		case error := <-errs:
			d.Log.Error("error occurred while listening to events. Cause %v", zap.Error(error))
			event.SetError(error)
		case message := <-messages:
			if message.Type == d.Docker.Type && message.Action == d.Docker.Action {
				d.Log.Info("event occurred", zap.String("type", message.Type), zap.String("action", message.Action))
				event.timestamp = time.Now().UTC()
				d.message = message
			}
		case <-d.stop:
			d.wg.Done()
			break
		}
	}
	d.events <- event
}

func (d *docker) Stop() error{
	d.stop <- struct{}{}
	close(d.stop)
	d.wg.Wait()
	return nil
}
