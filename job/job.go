package job

import (
	"github.com/kitchen-delivery/config"
	"github.com/kitchen-delivery/entity"
	"github.com/kitchen-delivery/service"
)

// Jobs holds both event-driven asynchronous jobs and scheduled jobs.
type Jobs struct {
	Order OrderJob
}

// InitializeJobs creates a new jobs instance.
func InitializeJobs(cfg config.AppConfig, services service.Services, queues entity.Queues) Jobs {
	orderJob := NewOrderJob(cfg, services, queues)

	return Jobs{
		Order: orderJob,
	}
}
