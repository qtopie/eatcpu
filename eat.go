package main

import (
	"runtime"
	"time"
)

// INTERVAL is the time range to adjust cpu usage
const INTERVAL int = 50

var (
	busyTime int
)

// take up cpu
// refresh cpu usage every second
func takeUpCPU() {
	defer wg.Done()

	runtime.LockOSThread()
	busyTime = int(float64(INTERVAL) * cpuUsage)
	// TODO: quick estimation

	for {
		// reset every 1 second
		idleTime := INTERVAL - busyTime

		for i := 0; i < 1000/INTERVAL; i++ {
			// supposed calculating cpu usage is pretty fast
			// sampling
			prevTicks, _ := getCPUTicks()

			takeNap(INTERVAL, idleTime)

			// sampling again and update idleTime
			ticks, _ := getCPUTicks()
			for i := 0; i < 8; i++ {
				ticks[i] = ticks[i] - prevTicks[i]
			}

			var total uint64
			for i := 0; i < 8; i++ {
				total += ticks[i]
			}

			idle := ticks[3] + ticks[4]
			delta := float64(idleTime) * (float64(idle)/float64(total) - (1.0 - cpuUsage)) / 0.75
			idleTime -= int(delta)
		}
	}

	runtime.UnlockOSThread()
}

func takeNap(interval, idleTime int) {
	startTime := time.Now()

	time.Sleep(time.Duration(idleTime) * time.Millisecond)

	// tick time
	d := time.Duration(interval) * time.Millisecond

	for time.Now().Sub(startTime) < d {
		// just to consume cpu
	}
}
