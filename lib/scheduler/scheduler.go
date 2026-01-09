package scheduler

import (
	"qbittorrent_exporter/lib/log"
	"sync"
	"time"
)

var (
	lock           sync.Mutex
	singleInstance *Scheduler
)

type (
	Scheduler struct {
		wg *sync.WaitGroup
	}

	PeriodicTaskOpts struct {
		Interval time.Duration
		IsFast   bool
	}

	taskFunc func() error
)

func Run(task taskFunc, po *PeriodicTaskOpts) {
	scheduler := Get()
	scheduler.wg.Add(1)

	go func() {
		defer scheduler.wg.Done()
		if po != nil {
			scheduler.RunPeriodicTask(task, po)
		} else if err := task(); err != nil {
			log.Error(err.Error())
		}
	}()
}

func Get() *Scheduler {
	if singleInstance == nil {
		func() {
			lock.Lock()
			defer lock.Unlock()
			var wg sync.WaitGroup
			singleInstance = &Scheduler{
				wg: &wg,
			}
		}()
	}
	return singleInstance
}

func (s *Scheduler) Wait() {
	s.wg.Wait()
}

func (s *Scheduler) RunPeriodicTask(task taskFunc, o *PeriodicTaskOpts) {
	opts := o
	if o == nil {
		log.Warn("PeriodicTaskOpts is nil, using default values")
		opts = &PeriodicTaskOpts{
			Interval: 30 * time.Second,
			IsFast:   false,
		}
	}
	ticker := time.NewTicker(opts.Interval)
	defer ticker.Stop()

	if opts.IsFast {
		if err := task(); err != nil {
			log.Error(err.Error())
		}
	}
	for range ticker.C {
		if err := task(); err != nil {
			log.Error(err.Error())
		}
	}
}
