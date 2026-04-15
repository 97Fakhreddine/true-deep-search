package debounce

import (
	"sync"
	"time"
)

type Debouncer struct {
	mu     sync.Mutex
	timer  *time.Timer
	delay  time.Duration
	action func()
}

func New(delay time.Duration, action func()) *Debouncer {
	return &Debouncer{
		delay:  delay,
		action: action,
	}
}

func (d *Debouncer) Trigger() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}

	d.timer = time.AfterFunc(d.delay, d.action)
}
