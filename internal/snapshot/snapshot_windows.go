package snapshot

import (
	"time"

	"github.com/elastic/go-windows"
	"github.com/olexnzarov/processinfo/internal/utility"
)

func getSystemTime() (time.Duration, error) {
	idle, kernel, user, err := windows.GetSystemTimes()
	if err != nil {
		return 0, utility.FormatError("failed to get the system time", err)
	}
	return kernel + user + idle, nil
}

func getProcessInfo(pid int64) (*ProcessInfo, error) {
	stats, err := getWinProcessStats(pid)
	if err != nil {
		return nil, err
	}
	time := (stats.kernelTime + stats.userTime)
	return &ProcessInfo{time: time, memory: stats.memory}, nil
}

func Get(pid int) (*Snapshot, error) {
	processInfo, err := getProcessInfo(pid)
	if err != nil {
		return nil, err
	}
	systemTime, err := getSystemTime()
	if err != nil {
		return nil, err
	}
	return &Snapshot{
		process:    processInfo,
		systemTime: systemTime,
	}, nil
}
