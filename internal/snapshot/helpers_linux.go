package snapshot

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olexnzarov/processinfo/internal/utility"
	"github.com/tklauser/go-sysconf"
)

var clockTicks *float64

type userHz float64

type procPidStat struct {
	userTime          userHz
	systemTime        userHz
	virtualMemorySize uint64
}

type procStat struct {
	userTime   userHz
	lpUserTime userHz
	systemTime userHz
	idleTime   userHz
}

func getProcPidStat(pid int) (*procPidStat, error) {
	bytes, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return nil, ErrProcessNotFound
	}
	line := string(bytes)
	if stat, ok := tryParseProcPidStat(line); ok {
		return stat, nil
	}
	return nil, utility.NewError(fmt.Sprintf("failed to parse/proc/%d/stat", pid))
}

func getProcStat() (*procStat, error) {
	bytes, err := os.ReadFile("/proc/stat")
	if err != nil {
		return nil, err
	}
	if stat, ok := tryParseProcStat(string(bytes)); ok {
		return stat, nil
	}
	return nil, utility.NewError("failed to parse/proc/stat")
}

func getClockTicks() float64 {
	if clockTicks == nil {
		value, err := sysconf.Sysconf(sysconf.SC_CLK_TCK)
		if err != nil {
			return 100
		}
		floatValue := float64(value)
		clockTicks = &floatValue
	}
	return *clockTicks
}

func (hz userHz) asDuration() time.Duration {
	seconds := float64(hz) / getClockTicks()
	return time.Duration(seconds * float64(time.Second))
}

func parseUserHz(value string) (userHz, error) {
	v, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return userHz(v), nil
}

func tryParseProcPidStat(line string) (*procPidStat, bool) {
	pidEnd := strings.LastIndexByte(line, ')')
	if pidEnd == -1 {
		return nil, false
	}

	properties := strings.Split(line[pidEnd+2:], " ")
	if len(properties) < 21 {
		return nil, false
	}

	utime, err := parseUserHz(properties[11])
	if err != nil {
		return nil, false
	}
	stime, err := parseUserHz(properties[12])
	if err != nil {
		return nil, false
	}
	vsize, err := strconv.ParseUint(properties[20], 10, 64)
	if err != nil {
		return nil, false
	}

	return &procPidStat{
		userTime:          utime,
		systemTime:        stime,
		virtualMemorySize: vsize,
	}, true
}

func tryParseProcStat(data string) (*procStat, bool) {
	line, _, ok := strings.Cut(data, "\n")
	if !ok {
		return nil, false
	}
	properties := strings.Fields(line)

	user, err := parseUserHz(properties[1])
	if err != nil {
		return nil, false
	}
	nice, err := parseUserHz(properties[2])
	if err != nil {
		return nil, false
	}
	system, err := parseUserHz(properties[3])
	if err != nil {
		return nil, false
	}
	idle, err := parseUserHz(properties[4])
	if err != nil {
		return nil, false
	}

	return &procStat{
		userTime:   user,
		lpUserTime: nice,
		systemTime: system,
		idleTime:   idle,
	}, true
}
