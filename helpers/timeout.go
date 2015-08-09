package helpers

import (
	"time"
)

type Timeout struct {
	Alarm chan bool
}

func NewTimeout(d time.Duration) *Timeout {
	timeout := new(Timeout)
	timeout.Alarm = make(chan bool, 1)
	go func() {
		time.Sleep(d)
		timeout.Alarm <- true
	}()
	return timeout
}
