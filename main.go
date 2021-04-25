package main

// sampleProc -- take a sample of N seconds of statistics from a program from /proc

import (
	"flag"
	"fmt"
	"github.com/acksin/procfs"
	"log"
	"strconv"
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

	flag.IntVar(& seconds, "seconds", 10, "length of sample in seconds")
	flag.Parse()

	period = time.Duration(seconds) * time.Second

	s := flag.Arg(0)
	pid, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("pid was not a number: %s", err)
	}

	proc, err := procfs.NewProc(pid)
	if err != nil {
		log.Fatalf("could not get process: %s", err)
	}
	sample(period, proc)
}

// sample takes a stat before and after a period of sample seconds and prints it
func sample(period time.Duration, p procfs.Proc) {
	before, err := p.NewStat()
	if err != nil {
		log.Fatalf("could not get process before: %s", err)
	}

	fmt.Printf("command:  %s\n", before.Comm)
	fmt.Printf("cpu time: %fs\n", before.CPUTime())
	fmt.Printf("vsize:    %dB\n", before.VirtualMemory())
	fmt.Printf("rss:      %dB\n", before.ResidentMemory())

	time.Sleep(period)

	after, err :=  p.NewStat()
	if err != nil {
		log.Fatalf("could not get process after: %s", err)
	}
	fmt.Printf("\nafter:\n")
	fmt.Printf("cpu time: %fs\n", after.CPUTime())
	fmt.Printf("vsize:    %dB\n", after.VirtualMemory())
	fmt.Printf("rss:      %dB\n", before.ResidentMemory())
	print(before)
}

// print does just that, in go format
func print( p procfs.ProcStat) {
	fmt.Printf("%#v\n", p)
}

// PidOf returns a pid array for a given program-name or error
func PidOf(name string) ([]int, error) {
	return nil, fmt.Errorf("not implemented yet")
}