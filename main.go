package main

// sampleCPU -- take a sample of N seconds of statistics from a program, using /proc.
//		Uses github.com/acksin/procfs , one of many by that name.

import (
	"flag"
	"fmt"
	"github.com/acksin/procfs"
	"log"
	"os"
	"sync"
	"time"
)

// usage prints out what this is
func usage() {
	log.Printf("sampleCPU -- take a sample of N seconds of statistics from one or more programs\n")
	log.Printf("Usage: sampleCPU [--seconds N] program ...")
	os.Exit(1)
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
	if len(flag.Args()) == 0 {
		usage()
	}
	period = time.Duration(seconds) * time.Second

	// scan /proc for processes matching names, a racy action
	fmt.Printf("#name, pid, cputime\n")
	for i := 0; i < flag.NArg(); i++ {
		name := flag.Arg(i)
		pids, err = PidOf(name)
		if len(pids) <= 0 {
			log.Printf("no process matched %s, ignored\n", name)
			continue
		}
		for _, pid := range pids {
			var proc procfs.Proc

			wg.Add(1)
			proc, err = procfs.NewProc(pid)
			if err != nil {
				log.Printf("could not get process for pid %d, ignored: %s ", pid, err)
				continue
			}
			go sample(period, proc, &wg)
		}
	}
	wg.Wait()
}

// sample takes a stat before and after a specified number of seconds and prints it
func sample(period time.Duration, p procfs.Proc, wg *sync.WaitGroup) {
	defer wg.Done()

	before, err := p.NewStat()
	if err != nil {
		log.Printf("could not get process %d, it had already exited: %s", p.PID, err)
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
// An empty array is not an error here, but it usually is to
// the caller
func PidOf(name string) ([]int, error) {
	var pids []int // make()?

	procs, err := procfs.AllProcs()
	if err != nil {
		return nil, fmt.Errorf("could not get list of processes")
	}
	for _, p := range procs {
		pn, err := p.Comm()
		if err != nil {
			return nil, fmt.Errorf("could not get a command-name for pid %d", p.PID)
		}
		if pn == name {
			pids = append(pids, p.PID)
		}
	}
	return pids, nil
}