package main

import (
	"flag"
	"fmt"
	"github.com/chaseisabelle/pipeline"
	"sync"
	"time"
)

type data struct {
	Iteration int
	Pipe      int
	Start     time.Time
}

var mux sync.Mutex
var min time.Duration
var max time.Duration
var tot int64

func main() {
	iterations := flag.Int("iterations", 1, "number of things to pump into the pipeline")
	pipes := flag.Int("pipes", 1, "number of pipes in the pipeline")
	handlers := flag.Int("handlers", 1, "number of parallel handlers per pipe")

	flag.Parse()

	pl := pipeline.Pipeline{}

	for p := 0; p < *pipes; p++ {
		fuckup(pl.Append(&pipeline.Pipe{
			Handler:  handler,
			Handlers: uint(*handlers),
			Retries:  0,
		}))
	}

	fuckup(pl.Append(&pipeline.Pipe{
		Handler:  finisher,
		Handlers: uint(*handlers),
		Retries:  0,
	}))

	fuckup(pl.Open())

	now := time.Now()

	for i := 1; i <= *iterations; i++ {
		fuckup(pl.Feed(data{
			Iteration: i,
			Pipe:      0,
			Start:     time.Now(),
		}))
	}

	fuckup(pl.Close())

	dur := time.Since(now)

	mux.Lock()

	defer mux.Unlock()

	fmt.Printf("min: %+v\nmax: %+v\navg: %+v\nall: %+v\n", min, max, time.Duration(float64(tot)/float64(*iterations)), dur)
}

func fuckup(err error) {
	if err != nil {
		panic(err)
	}
}

func handler(i interface{}) (interface{}, error) {
	d, ok := i.(data)

	if !ok {
		return nil, fmt.Errorf("invalid handler input %+v", i)
	}

	d.Pipe++

	t := time.Since(d.Start)

	fmt.Printf("%+v: delta: %+v\n", d, t)

	return d, nil
}

func finisher(i interface{}) (interface{}, error) {
	d, ok := i.(data)

	if !ok {
		return nil, fmt.Errorf("invalid finisher input %+v", i)
	}

	t := time.Since(d.Start)

	mux.Lock()

	defer mux.Unlock()

	if min > t || min == 0 {
		min = t
	}

	if max < t {
		max = t
	}

	tot += t.Nanoseconds()

	return nil, nil
}
