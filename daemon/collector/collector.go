package collector

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type (
	Stats struct {
		Host    Resources
		Runtime RuntimeStats
	}

	RuntimeStats struct {
		RunningTasks uint32
	}

	Resources struct {
		Load float64
		RAM  uint64
	}

	Thresholds struct {
		Load float64
		RAM  uint64
	}

	Collector interface {
		Collect(context.Context) (Stats, error)
	}

	collector struct {
	}
)

func (t Thresholds) Upholds(s Stats) bool {
	return s.Host.Load < t.Load && s.Host.RAM < t.RAM
}

func NewCollector() Collector {
	return &collector{}
}

func (c *collector) Collect(ctx context.Context) (Stats, error) {
	log.Debug("collecting runtime stats")

	log.Debug("collecting host stats")
	// todo: collect host stats

	return Stats{
		Host: Resources{},
	}, nil
}
