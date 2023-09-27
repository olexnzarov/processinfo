package processinfo

import (
	"context"
	"time"

	"github.com/olexnzarov/processinfo/internal/snapshot"
)

type sampler struct {
	pid      int
	channel  chan *ProcessInfo
	previous *snapshot.Snapshot
}

// NewSampler returns a channel to which process information will be sent every 'rate' moment.
// It will contain an average CPU utilization during that time.
// The channel will be closed when the context is done.
// If any error occurs during the sampling, the information sent will be zeroed.
func NewSampler(ctx context.Context, pid int, rate time.Duration) <-chan *ProcessInfo {
	sampler := &sampler{
		pid:      pid,
		channel:  make(chan *ProcessInfo, 1),
		previous: snapshot.Zero,
	}
	go sampler.start(ctx, rate)
	return sampler.channel
}

func (s *sampler) sendSample() {
	current, err := snapshot.Get(s.pid)

	// If something happened, just set everything to zero.
	// Process might not be running anymore.
	if err != nil {
		current = snapshot.Zero
		s.previous = snapshot.Zero
	}

	info := &ProcessInfo{
		CPU:    current.GetProcessorUtilization(s.previous),
		Memory: current.GetMemory(),
	}
	s.previous = current

	// Send the information non-blockingly.
	select {
	case s.channel <- info:
	default:
	}
}

func (s *sampler) start(ctx context.Context, rate time.Duration) {
	defer close(s.channel)
	s.sendSample()
	ticker := time.NewTicker(rate)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.sendSample()
		}
	}
}
