package main

import (
	"errors"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

// Tick contains the amount  of  time,  measured  in  units  of  USER_HZ
//  (1/100ths  of  a  second  on  most   architectures,   use
//  sysconf(_SC_CLK_TCK) to obtain the right value), that the
//  system ("cpu" line) or the  specific  CPU. See `man 5 proc` in linux manual.
type Tick struct {
	user      uint64
	nice      uint64
	system    uint64
	idle      uint64
	iowait    uint64
	irq       uint64
	softirq   uint64
	steal     uint64
	guest     uint64
	guestNice uint64
}

func (t *Tick) load() (err error) {
	ticks, err := getCPUTicks()
	if err != nil {
		return
	}

	if len(ticks) < 8 {
		return errors.New("Len is too short")
	}

	t.user = ticks[0]
	t.nice = ticks[1]
	t.system = ticks[2]
	t.idle = ticks[3]
	t.iowait = ticks[4]
	t.irq = ticks[5]
	t.softirq = ticks[6]
	t.steal = ticks[7]

	if len(ticks) >= 10 {
		t.guest = ticks[8]
		t.guestNice = ticks[9]
	}

	return
}

func getCPUTicks() (ticks []uint64, err error) {
	statData, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(statData), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "cpu ") {
			fields := strings.Fields(line)
			fields = fields[1:]

			for _, v := range fields {
				tick, err := strconv.ParseUint(v, 10, 64)
				if err != nil {
					log.Fatal(err)
				}
				ticks = append(ticks, tick)
			}
			return

		}
	}

	return
}

// totalCpu = user + nice + system + idle + iowait + irq + softirq + steal (ignore guest)
// idleCpu = idle + iowait
func getAvgCpuUsage(period int64) (cpuUsage float64, err error) {
	prevTicks, err := getCPUTicks()
	if err != nil {
		return
	}

	// sleep for a while
	time.Sleep(time.Duration(period) * time.Millisecond)

	ticks, err := getCPUTicks()
	if err != nil {
		return
	}

	for i := 0; i < 8; i++ {
		ticks[i] = ticks[i] - prevTicks[i]
	}

	var total uint64
	for i := 0; i < 8; i++ {
		total += ticks[i]
	}

	idle := ticks[3] + ticks[4]
	cpuUsage = 1.0 - float64(idle)/float64(total)
	return
}
