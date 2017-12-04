package main

import (
	"testing"
)

func TestGetCPUTicks(t *testing.T) {
	ticks, err := getCPUTicks()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	for _, tick := range ticks {
		t.Log(tick, "\t")
	}
}

func TestGetAvgCpuUsage(t *testing.T) {
	load, err := getAvgCpuUsage(100)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	t.Log(load)
}
