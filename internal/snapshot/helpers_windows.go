package snapshot

import (
	"syscall"
	"time"

	"github.com/elastic/go-windows"
	"github.com/olexnzarov/processinfo/internal/utility"
)

type winProcStat struct {
	kernelTime time.Duration
	userTime   time.Duration
	memory     uint64
}

// Windows Server 2003 and Windows XP: This access right is not supported.
// TODO: fallback to PROCESS_QUERY_INFORMATION
const PROCESS_QUERY_LIMITED_INFORMATION = 0x00001000

func getWinProcessStats(pid int) (*winProcStat, error) {
	handle, err := syscall.OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION, true, uint32(pid))
	if err != nil {
		return nil, ErrProcessNotFound
	}
	defer syscall.CloseHandle(handle)

	// We have no use for creation and exit times, we need only kernel and user time.
	var temp, kernelTime, userTime syscall.Filetime
	if err := syscall.GetProcessTimes(handle, &temp, &temp, &kernelTime, &userTime); err != nil {
		return nil, utility.FormatError("failed to get the process time", err)
	}

	mem, err := windows.GetProcessMemoryInfo(handle)
	if err != nil {
		return nil, utility.FormatError("failed to get the process memory", err)
	}

	return &winProcStat{
		kernelTime: windows.FiletimeToDuration(&kernelTime),
		userTime:   windows.FiletimeToDuration(&userTime),
		memory:     uint64(mem.WorkingSetSize),
	}, nil
}
