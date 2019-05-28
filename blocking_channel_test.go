package tracer_test

import (
	"github.com/RobinUS2/golang-chan-tracer"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	wg := sync.WaitGroup{}
	a := make(chan bool)

	// stats
	var written, read uint64

	trace := tracer.New("test1")

	// producer
	write := func() {
		wg.Add(1)
		go func() {
			trace.WriteBool(a, true)
			atomic.AddUint64(&written, 1)
		}()
	}

	// consumer
	go func() {
		for {
			<-a
			wg.Done()
			atomic.AddUint64(&read, 1)
		}
	}()
	write()
	wg.Wait()

	if written != 1 {
		t.Error(written)
	}
	if read != 1 {
		t.Error(read)
	}

	if trace.BlockingWrites() != 0 {
		t.Error(trace.BlockingWrites())
	}

	if trace.NonBlockingWrites() != 1 {
		t.Error(trace.NonBlockingWrites())
	}
	if trace.AvgBlockedTime() != 0 {
		t.Error()
	}
}

func TestNewBlocking(t *testing.T) {
	wg := sync.WaitGroup{}
	a := make(chan bool)

	// stats
	var written, read uint64

	trace := tracer.New("test2")

	// producer
	write := func() {
		wg.Add(1)
		go func() {
			trace.WriteBool(a, true)
			atomic.AddUint64(&written, 1)
		}()
	}

	// consumer
	go func() {
		for {
			<-a
			time.Sleep(50 * time.Millisecond)
			wg.Done()
			atomic.AddUint64(&read, 1)
		}
	}()
	time.Sleep(10 * time.Millisecond)
	write()
	go write()

	wg.Wait()

	if written != 2 {
		t.Error(written)
	}
	if read != 2 {
		t.Error(read)
	}

	if trace.BlockingWrites() != 1 {
		t.Error(trace.BlockingWrites())
	}

	if trace.NonBlockingWrites() != 1 {
		t.Error(trace.NonBlockingWrites())
	}
	if trace.AvgBlockedTime() < 30 || trace.AvgBlockedTime() > 70 {
		t.Error(trace.AvgBlockedTime())
	}
}
