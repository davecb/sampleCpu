package main

// sampleProc -- take a sample of N seconds of statistics from a program from /proc
//		uses github.com/prometheus/procfs

import (
	"flag"
	"fmt"
	"github.com/acksin/procfs"
	"log"
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

	flag.IntVar(& seconds, "seconds", 10, "length of sample in seconds")
	flag.Parse()

	period = time.Duration(seconds) * time.Second

	// scan /proc for processes matching name
	name := flag.Arg(0)
	pids, err = PidOf(name)
	if len (pids) <= 0 {
		log.Fatalf("no process matched %s\n", name)
	}

	// do exactly and only the first
	proc, err := procfs.NewProc(pids[0])
	if err != nil {
		log.Fatalf("could not get process: %s", err)
	}
	sample(period, proc)
}

// sample takes a stat before and after a period of sample seconds and prints it
func sample(period time.Duration, p procfs.Proc) {
	before, err := p.NewStat()
	if err != nil {
		log.Fatalf("could not get process at beginning: %s", err)
	}

	time.Sleep(period)

	after, err :=  p.NewStat()
	if err != nil {
		log.Fatalf("could not get process at end: %s", err)
	}

    // report results
	fmt.Printf("#name, cputime\n")
	fmt.Printf("%s, %f\n", before.Comm, after.CPUTime() - before.CPUTime())
}

// print does just that, in go format
func print( p procfs.ProcStat) {
	fmt.Printf("%#v\n", p)
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
	//return nil, fmt.Errorf("not implemented yet")
}