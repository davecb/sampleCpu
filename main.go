package main

// sampleProc -- take a sample of N seconds of statistics from a program from /proc
//		uses github.com/prometheus/procfs

import (
	"flag"
	"fmt"
	"github.com/acksin/procfs"
	"log"
	"sync"
	"time"
)

// usage prints out what this is
func usage() {
	log.Printf("sampleProc -- take a sample of N seconds of statistics from a program\n")
	log.Printf("Usage: sampleProc [--seconds N] program")
}

// main -- start, parse args
func main() {
	var seconds int
	var period time.Duration
	var pids []int
	var err error
	var wg sync.WaitGroup

	flag.IntVar(& seconds, "seconds", 10, "length of sample in seconds")
	flag.Parse()

	period = time.Duration(seconds) * time.Second

	// scan /proc for processes matching name
	name := flag.Arg(0)
	pids, err = PidOf(name)
	if len (pids) <= 0 {
		log.Fatalf("no process matched %s\n", name)
	}
	fmt.Printf("#name, pid, cputime\n")
	for _, pid := range(pids) {
		wg.Add(1)
		var proc procfs.Proc
		proc, err = procfs.NewProc(pid)
		if err != nil {
			log.Fatalf("could not get process: %s", err)
		}
		go sample(period, proc, &wg)
	}
	wg.Wait()
}

// sample takes a stat before and after a period of sample seconds and prints it
func sample(period time.Duration, p procfs.Proc, wg *sync.WaitGroup) {
	defer wg.Done()

	before, err := p.NewStat()
	if err != nil {
		log.Fatalf("could not get process, it exited: %s", err)
		return
	}

	time.Sleep(period)

	after, err :=  p.NewStat()
	if err != nil {
		fmt.Printf("%s, %d, exited\n", before.Comm, before.PID)

		return
	}

    // report results
	fmt.Printf("%s, %d, %f\n", before.Comm, before.PID, after.CPUTime() - before.CPUTime())
}

// PidOf returns a pid array for a given program-name or error
// An empty array is not an error here, but it probably is to
// the caller (;-))
func PidOf(name string) ([]int, error) {
	var pids []int // make()?

	procs, err := procfs.AllProcs()
	if err != nil {
		return nil, fmt.Errorf("could not get list of processes")
	}
	for _, p := range procs {
		pn, err := p.Comm()
		if err != nil {
			return nil, fmt.Errorf("could not a command-name from /proc")
		}
		if pn == name {
			pids = append(pids, p.PID)
		}
	}
	return pids, nil
}