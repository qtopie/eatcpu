package main

import (
	"flag"
	"runtime"
	"sync"
)

var (
	wg       sync.WaitGroup
	cpuUsage float64
)

func main() {

	flag.Float64Var(&cpuUsage, "percent", 0.5, "cpu usage (%)")
	flag.Parse()

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go takeUpCPU()
	}

	wg.Wait()
}
