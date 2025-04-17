package scheduler

import (
	"qbittorrent_exporter/lib/log"
	"time"
)

type taskFunc func() error

type PeriodicTaskOpts struct {
	Interval time.Duration
	IsFast   bool
}

func Default() *PeriodicTaskOpts {
	return &PeriodicTaskOpts{
		Interval: 30 * time.Second,
		IsFast:   false,
	}
}

func RunPeriodicTask(task taskFunc, o *PeriodicTaskOpts) {
	if o == nil {
		log.Info("PeriodicTaskOpts is empty, using default values")
		o = Default()
	}
	ticker := time.NewTicker(o.Interval)
	defer ticker.Stop()

	if o.IsFast {
		if err := task(); err != nil {
			log.Error(err.Error())
		}
	}
	for {
		select {
		case <-ticker.C:
			if err := task(); err != nil {
				log.Error(err.Error())
				continue
			}
		}
	}
}
