package snapshot

import (
	"time"

	"github.com/olexnzarov/processinfo/internal/utility"
)

var Zero = &Snapshot{process: &ProcessInfo{time: 0, memory: 0}, systemTime: 0}
var ErrProcessNotFound = utility.NewError("process not found")

type ProcessInfo struct {
	time   time.Duration
	memory uint64
}

type Snapshot struct {
	process    *ProcessInfo
	systemTime time.Duration
}

func (now *Snapshot) GetMemory() uint64 {
	return now.process.memory
}

func (now *Snapshot) GetProcessorUtilization(before *Snapshot) float64 {
	process := now.process.time - before.process.time
	system := now.systemTime - before.systemTime
	if system == 0 {
		return 0
	}
	return 100 * process.Seconds() / system.Seconds()
}
