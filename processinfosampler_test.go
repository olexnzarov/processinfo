package processinfo

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestSamplerNoProcess(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1000)
	defer cancel()
	ch := NewSampler(ctx, -1, time.Millisecond)
	<-ch
	<-ch
}

func TestSampler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1000)
	defer cancel()
	ch := NewSampler(ctx, os.Getpid(), time.Millisecond)
	<-ch
	<-ch
}

func TestSamplerChannelClose(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	closed := make(chan any)
	timeout := make(chan any)

	go func() {
		time.Sleep(time.Millisecond * 600)
		timeout <- nil
	}()

	go func() {
		ch := NewSampler(ctx, os.Getpid(), time.Second*10)
		for {
			_, ok := <-ch
			if !ok {
				break
			}
		}
		closed <- nil
	}()

	select {
	case <-closed:
	case <-timeout:
		t.Fatal("sampler channel did not close in time")
		return
	}
}
