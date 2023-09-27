package processinfo

import (
	"os"
	"testing"

	"github.com/olexnzarov/processinfo/internal/snapshot"
)

func TestGetNoProcess(t *testing.T) {
	if _, err := Get(-1); err != snapshot.ErrProcessNotFound {
		t.Fatalf("Get(-1) returns not ErrProcessNotFound: %s", err)
	}
}

func TestGet(t *testing.T) {
	if _, err := Get(os.Getpid()); err != nil {
		t.Fatalf("Get(os.Getpid()) returns error: %s", err)
	}
}
