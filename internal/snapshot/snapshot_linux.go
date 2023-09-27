package snapshot

import (
	"time"
)

func getProcessInfo(pid int) (*ProcessInfo, error) {
	stat, err := getProcPidStat(pid)
	if err != nil {
		return nil, err
	}
	return &ProcessInfo{
		time:   (stat.userTime + stat.systemTime).asDuration(),
		memory: stat.virtualMemorySize,
	}, nil
}

func getSystemTime() (time.Duration, error) {
	times, err := getProcStat()
	if err != nil {
		return 0, err
	}
	return (times.userTime + times.lpUserTime + times.systemTime + times.idleTime).asDuration(), nil
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
