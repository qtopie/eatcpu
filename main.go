package main

import (
	"flag"
	"runtime"
	"sync"
	"time"
)

var (
	wg       sync.WaitGroup
	cpuUsage float64
)

func init() {

}

func main() {

	flag.Float64Var(&cpuUsage, "percent", 0.5, "cpu usage (%)")
	flag.Parse()

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go takeUpCPU()
	}

	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			wg.Add(1)
			go recalculate()
			// case <-p.exit:
			// 	ticker.Stop()
		}
	}

	wg.Wait()
}
