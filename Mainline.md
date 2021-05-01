This is the implementation, in mainline.md

```go "subroutines" +=
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
	var proc procfs.Proc
	var pids []int
	var err error
	var wg sync.WaitGroup

	flag.IntVar(& seconds, "seconds", 10, "length of sample in seconds")
	flag.Parse()
	if len(flag.Args()) == 0 {
		usage()
	}
	fmt.Printf("#name, pid, seconds, cputime\n")


	// scan /proc for processes matching names, a racy action
	for i := 0; i < flag.NArg(); i++ {
		name := flag.Arg(i)
		pids, err = PidOf(name)
		if len(pids) <= 0 {
			log.Printf("no process matched %s, ignored\n", name)
			continue
		}
		for _, pid := range pids {

			wg.Add(1)
			proc, err = procfs.NewProc(pid)
			if err != nil {
				log.Printf("could not get process for pid %d, ignored: %s ", pid, err)
				continue
			}
			go sample(seconds, proc, &wg)
		}
	}
	wg.Wait()
}


// sample takes a stat before and after a specified number of seconds and prints it
func sample(period int, p procfs.Proc, wg *sync.WaitGroup) {
    var i int
	defer wg.Done()

    <<<sample collection>>>
	
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
```


first part
```go "get before value" +=
    before, err := p.NewStat()
	if err != nil {
		log.Printf("could not get process %d, it had already exited: %s", p.PID, err)
		return
	}
```

