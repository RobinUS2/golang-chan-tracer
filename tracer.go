package tracer

import (
	"log"
	"sync/atomic"
	"time"
)

const DefaultName string = "Default"

var DefaultInstance = New(DefaultName)

type Instance struct {
	name               string
	nonBlockingWrites  uint64
	blockingWrites     uint64
	totalTimeBlockedMs uint64 // in milliseconds
}

func New(name string) *Instance {
	return &Instance{
		name: name,
	}
}

func WriteBool(c chan bool, v bool) {
	DefaultInstance.WriteBool(c, v)
}

func (i *Instance) WriteBool(c chan bool, v bool) {
	select {
	case c <- v:
		// none blocking write succeeded
		atomic.AddUint64(&i.nonBlockingWrites, 1)
	default:
		start := time.Now()
		// blocking write
		c <- v
		took := time.Since(start)
		ms := took.Nanoseconds() / 1000000
		// we only count it as blocking if it really took time
		if ms > 0 {
			atomic.AddUint64(&i.blockingWrites, 1)
			atomic.AddUint64(&i.totalTimeBlockedMs, uint64(ms))
			log.Printf("[%s] write took %dms", i.name, ms)
		}
	}
}

func (i *Instance) NonBlockingWrites() uint64 {
	return atomic.LoadUint64(&i.nonBlockingWrites)
}

func (i *Instance) BlockingWrites() uint64 {
	return atomic.LoadUint64(&i.blockingWrites)
}

func (i *Instance) AvgBlockedTime() float64 {
	n := atomic.LoadUint64(&i.blockingWrites)
	total := atomic.LoadUint64(&i.totalTimeBlockedMs)
	return float64(total) / float64(n)
}

func NonBlockingWrites() uint64 {
	return DefaultInstance.NonBlockingWrites()
}

func BlockingWrites() uint64 {
	return DefaultInstance.BlockingWrites()
}

func AvgBlockedTime() float64 {
	return DefaultInstance.AvgBlockedTime()
}
