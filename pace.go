// Package pace provides a threadsafe counter for measuring ticks in the specified timeframe.
package pace

import (
	"log"
	"strconv"
	"sync"
	"time"
)

// Pace is a an interface to register ticks, force reporting and pause/resume the meter.
type Pace interface {
	// Step increments the counter of pace.
	Step(n float64)
	// Pause stops reporting until resumed, all steps continue to be counted.
	Pause()
	// Resume resumes the reporting, starting a report with info since the last tick.
	// Specify a new interval or 0 if you don't want to override it.
	Resume(interval time.Duration)
	// Report manually triggers a report with time frame less than the defined interval.
	// Specify a custom reporter function just for this one report.
	Report(reporter ReporterFunc)
}

type paceImpl struct {
	mux *sync.RWMutex

	value    float64
	label    string
	paused   bool
	interval time.Duration
	lastTick time.Time
	repFn    ReporterFunc
	t        *time.Timer
}

func (p *paceImpl) Step(n float64) {
	p.mux.Lock()
	p.value += n
	p.mux.Unlock()
}

func (p *paceImpl) Pause() {
	p.t.Stop()

	p.mux.Lock()
	defer p.mux.Unlock()
	p.report(nil)

	p.paused = true
	p.value = 0
	p.lastTick = time.Now()
}

func (p *paceImpl) Resume(interval time.Duration) {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.report(nil)

	p.paused = false
	p.value = 0
	p.lastTick = time.Now()
	if interval > 0 {
		// override the interval if provided
		p.interval = interval
	}
	p.t.Reset(p.interval)
}

func (p *paceImpl) Report(reporter ReporterFunc) {
	p.t.Stop()
	p.mux.Lock()
	defer p.mux.Unlock()
	p.report(reporter)

	p.value = 0
	p.lastTick = time.Now()
	if !p.paused {
		p.t.Reset(p.interval)
	}
}

func (p *paceImpl) report(reporter ReporterFunc) {
	if reporter == nil {
		reporter = p.repFn
	}
	timeframe := time.Since(p.lastTick)
	if abs(timeframe-p.interval) < 10*time.Millisecond {
		timeframe = p.interval
	}
	label := p.label
	value := p.value
	reporter(label, timeframe, value)
}

// New creates a new pace meter with provided label and reporting function.
// All ticks (or steps) are aggregated in timeframes specified using interval.
func New(label string, interval time.Duration, repFn ReporterFunc) Pace {
	if repFn == nil {
		repFn = DefaultReporter
	}
	p := &paceImpl{
		mux: new(sync.RWMutex),

		label:    label,
		interval: interval,
		repFn:    repFn,
		lastTick: time.Now(),
		t:        time.NewTimer(interval),
	}
	go func() {
		for range p.t.C {
			func() {
				p.mux.Lock()
				defer p.mux.Unlock()
				p.report(nil)

				p.value = 0
				p.lastTick = time.Now()
				p.t.Reset(interval)
			}()
		}
	}()
	return p
}

// ReporterFunc defines a function used to report current pace.
type ReporterFunc func(label string, timeframe time.Duration, value float64)

// DefaultReporter reports using log.Printf.
func DefaultReporter(label string, timeframe time.Duration, value float64) {
	floatFmt := func(f float64) string {
		return strconv.FormatFloat(value, 'f', -1, 64)
	}
	switch timeframe {
	case time.Second:
		log.Printf("%s: %s/s in %v", label, floatFmt(value), timeframe)
	case time.Minute:
		log.Printf("%s: %s/m in %v", label, floatFmt(value), timeframe)
	case time.Hour:
		log.Printf("%s: %s/h in %v", label, floatFmt(value), timeframe)
	case 24 * time.Hour:
		log.Printf("%s: %s/day in %v", label, floatFmt(value), timeframe)
	default:
		log.Printf("%s %s in %v (pace: %s/s)", floatFmt(value), label,
			timeframe, floatFmt(value/float64(timeframe)/float64(time.Second)))
	}
}

func abs(v time.Duration) time.Duration {
	if v < 0 {
		return -v
	}
	return v
}
