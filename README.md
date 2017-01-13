## Pace [![GoDoc](https://godoc.org/github.com/xlab/pace?status.svg)](https://godoc.org/github.com/xlab/pace) [![Go Report Card](http://goreportcard.com/badge/github.com/xlab/pace)](http://goreportcard.com/report/github.com/xlab/pace)

Pace is a Go package that helps to answer one simple question:

> how fast it goes?

![pace](https://cl.ly/2P163D3Q0O03/pace.gif)

It's a threadsafe counter that measures ticks in the specified timeframe. It also has a simple and intuitive interface:

```go
func New(label string, interval time.Duration, repFn ReporterFunc) Pace

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

// ReporterFunc defines a function used to report current pace.
type ReporterFunc func(label string, timeframe time.Duration, value float64)
```

### Installation

```
$ go get github.com/xlab/pace
```

### Usage example:

```go

// initialise a pace meter
p := New("items", time.Second, nil)
go func() {
    for range items {
        wg.Done()
        p.Step(1)
    }
}()

// start pushing items:

// pushing each 1ms
push(1*time.Millisecond, 3*time.Second)
// pushing each 10ms
push(10*time.Millisecond, 3*time.Second)
// pushing each 100ms
push(100*time.Millisecond, 3*time.Second)
// pushing each 500ms
push(500*time.Millisecond, 3*time.Second)
```

Full code available at [pace_test.go](/pace_test.go).

#### Output:

```
$ go test
2017/01/13 13:29:51 items: 999/s in 1s
2017/01/13 13:29:52 items: 1001/s in 1s
2017/01/13 13:29:53 items: 1000/s in 1s
2017/01/13 13:29:54 items: 100/s in 1s
2017/01/13 13:29:55 items: 100/s in 1s
2017/01/13 13:29:56 items: 100/s in 1s
2017/01/13 13:29:57 items: 10/s in 1s
2017/01/13 13:29:58 items: 10/s in 1s
2017/01/13 13:29:59 items: 10/s in 1s
2017/01/13 13:30:00 items: 2/s in 1s
2017/01/13 13:30:01 items: 2/s in 1s
2017/01/13 13:30:02 items: 2/s in 1s
2017/01/13 13:30:02 done
PASS
ok      github.com/xlab/pace    12.006s
```

Also within **5 second** timeframe using `pace.New("items", 5 * time.Second, nil)`

```
$ go test
2017/01/13 13:32:06 3199 items in 5s (pace: 3199/s)
2017/01/13 13:32:11 133 items in 5s (pace: 133/s)
2017/01/13 13:32:13 4 items in 1.999796727s (pace: 4/s)
2017/01/13 13:32:13 done
PASS
ok      github.com/xlab/pace    12.006s
```

### License

[MIT](/LICENSE.txt)
