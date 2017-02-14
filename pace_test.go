package pace

import (
	"log"
	"sync"
	"testing"
	"time"
)

const timeframe = time.Second

func TestSimple(t *testing.T) {
	items := make(chan struct{}, 100)
	wg := new(sync.WaitGroup)

	p := New("items", timeframe, nil)
	go func() {
		for range items {
			wg.Done()
			p.Step(1)
		}
	}()

	push := func(interval, duration time.Duration) {
		tick := time.NewTicker(interval)
		start := time.Now()
		for range tick.C {
			wg.Add(1)
			items <- struct{}{}
			if time.Since(start) > duration {
				break
			}
		}
	}
	push(1*time.Millisecond, 3*time.Second)
	push(10*time.Millisecond, 3*time.Second)
	push(100*time.Millisecond, 3*time.Second)
	push(500*time.Millisecond, 3*time.Second)

	wg.Wait()
	time.Sleep(3 * time.Second)
	push(10*time.Millisecond, 3*time.Second)
	time.Sleep(3 * time.Second)
	p.Report(nil)
	log.Println("done")
}
