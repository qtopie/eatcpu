package main

import (
	"log"
	"runtime"
	"sync"
	"time"
)

// INTERVAL is set to 200 ms
const INTERVAL int = 200

var (
	busyTime int
	l        sync.Mutex
)

// take up cpu
// refresh cpu usage every second
func takeUpCPU() {
	defer wg.Done()

	busyTime = int(float64(INTERVAL) * cpuUsage)
	tickTime := busyTime
	runtime.LockOSThread()

	for {
		idleTime := INTERVAL - tickTime
		for i := 0; i < 4; i++ {
			// consume cpu
			//			run(busyTime, idleTime) 200ms * 4
			busy(tickTime, idleTime)
		}

		// calcate cpu usage and  take up some cpu   200ms * 1
		// update busyTime to meet cpu usage
		tickTime = updateAndTick()
	}
	runtime.UnlockOSThread()
}

func busy(busyTime, idleTime int) {
	startTime := time.Now()
	d := time.Duration(busyTime) * time.Millisecond
	for time.Now().Sub(startTime) < d {
		// Just loop to consume CPU
	}

	time.Sleep(time.Duration(idleTime) * time.Millisecond)
}

func updateAndTick() (tickTime int) {

	time.Sleep(time.Duration(busyTime) * time.Millisecond)

	// update busyTime
	return busyTime
}

func recalculate() {
	defer wg.Done()

	load, err := getAvgCpuUsage(500)
	if err != nil {
		return
	}
	log.Println("Current cpu usage(%):", load*100)

	// TODO
	busy := int((cpuUsage - load) * 50)
	l.Lock()
	busyTime += busy
	if busyTime > 200 {
		busyTime = 200 - int(1-cpuUsage)*10
	}
	l.Unlock()
	log.Println("Busy time", busyTime)
}
