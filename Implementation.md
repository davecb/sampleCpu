This is the implementation, in mainline.md

```go "subroutines" +=
    <<<init>>>
    <<<main>>>
    <<<sample>>>
    <<<pidof>>>
    <<<sample one name>>>

```


```go "init" += 
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
```

```go "main" += 
// main -- start, parse args
func main() {
	var seconds int
	var wg sync.WaitGroup

	flag.IntVar(& seconds, "seconds", 10, "length of sample in seconds")
	flag.Parse()
	if len(flag.Args()) == 0 {
		usage()
	}
	fmt.Printf("#name, pid, seconds, cputime\n")

    <<<main loop>>>   

}
```

```go "sample" +=

// sample takes a stat before and after a specified number of seconds and prints it
func sample(period int, p procfs.Proc, wg *sync.WaitGroup) {
    var i int
	defer wg.Done()

    <<<sample collection>>>
    
 }
```

```go "pidof" +=

// PidOf returns a pid array for a given program-name or error
// An empty array is not an error here, but it usually is to
// the caller
func NewPIDs(name string, pids map[int]bool) ([]int, map[int]bool, error) {
	var newPids = make([]int,0) 
	
	procs, err := procfs.AllProcs()
	if err != nil {
		return nil, pids, fmt.Errorf("could not get list of processes")
	}
	for _, p := range procs {
		pn, err := p.Comm()
		if err != nil {
			return nil, pids, fmt.Errorf("could not get a command-name for pid %d", p.PID)
		}
		if pn == name {
		    if _, present := pids[p.PID]; !present {
		        pids[p.PID] = true
			    newPids = append(newPids, p.PID)
			}
		}
	}
	return newPids, pids, nil
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


