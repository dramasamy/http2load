package main

import (
	"math"
	"sync"
)

type Gauge struct {
	lock sync.Mutex
	sum  float64

	Count int
	Min   float64
	Max   float64
}

func NewGauge() *Gauge {
	gauge := &Gauge{}
	gauge.Reset()
	return gauge
}

func (gauge *Gauge) Reset() {
	gauge.sum, gauge.Count, gauge.Min, gauge.Max = 0, 0, math.MaxFloat64, 0
}

func (gauge *Gauge) Mean() float64 {
	return gauge.sum / float64(gauge.Count)
}

func (gauge *Gauge) Add(n float64) {
	gauge.lock.Lock()
	defer gauge.lock.Unlock()
	if n < gauge.Min {
		gauge.Min = n
	}
	if n > gauge.Max {
		gauge.Max = n
	}
	gauge.Count++
	gauge.sum = gauge.sum + n
}