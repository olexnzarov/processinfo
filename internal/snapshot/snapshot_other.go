//go:build !windows && !linux

package snapshot

import "github.com/olexnzarov/processinfo/internal/utility"

func Get(pid int) (*Snapshot, error) {
	return nil, utility.NewError("system is not supported")
}
