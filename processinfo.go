package processinfo

import (
	"github.com/olexnzarov/processinfo/internal/snapshot"
)

type ProcessInfo struct {
	CPU    float64 // Percentage of overall CPU utilization by the process.
	Memory uint64  // Memory used by the process in bytes.
}

// Get returns resources used by the process in the current instant.
//
// It will contain very inaccurate CPU utilization value as it has no time reference.
// Use 'processinfo.NewSampler' to get an accurate value on ongoing basis.
func Get(pid int) (*ProcessInfo, error) {
	data, err := snapshot.Get(pid)
	if err != nil {
		return nil, err
	}
	return &ProcessInfo{
		CPU:    data.GetProcessorUtilization(snapshot.Zero),
		Memory: data.GetMemory(),
	}, nil
}
