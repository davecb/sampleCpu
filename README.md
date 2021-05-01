# sampleCPU

[sampleCPU](https://github.com/davecb/sampleCPU) is a tool to take
cpu samples. Sound easy? Nope!

All I wanted to do was sure how much CPU per minute an old multi-process batch
program was using, to compare to a new, multithreaded, Go program. In an afternoon.

The only tricky thing was that the old batch program spun off 
short-lived children, each taking a few inputs, processing them, then 
shutting down and reporting the results

The program looked like 
```go "usused code" +=
for each program mentioned
    <<<get list of PIDs>>>
    for each pid
        <<<sample their CPU use>>>
```
```go main.go
<<<subroutines>>>
```


The first time I tried measuring, I got stats from one go program and 126 batch
children, 22 of which had exited before the sampling period was over.
That probably meant that 22 more had started up, for a margin of error of
44/148 or 29%. Drat!

The only reason I got an answer I could even _report_ was that the Go program
was _hugely_ better than the batch one, and the margin of error was small
by comparison.

But I wasn't pleased.

## Let's do it right

It's the weekend, and I have more than an afternoon to spend.  So let's figure out
what we really need to do to get decent stats from /proc about a collection of
running programs.

We're fine for the Go program: it just runs, and we can look at /proc/pid/status 
and get CPU. Less so for the batch one.

The first thing I need to do is make the running program report CPU even
if the program exits before the time period is up.
The initial measurement code looked like
```go "unused code" += 
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

```

In effect, it was
```go "unused code" +=
    <<<get before value>>>
    <<<sleep, get after value>>>
    fmt.Printf("%s, %d, %f\n", before.Comm, before.PID, after.CPUTime() - before.CPUTime())
```

The collection of the "before" value was still correct: it wrote 
a warning message to stderr and didn't try to do anything special.

The "after" value collector needed work, though. It needed to sample repeatedly, and 
use the last good one to compute the time spent.

That turned it into a for-loop, with a timer and a select instead of a Sleep()
```go "get after value"
    ticker := time.NewTicker(1 * time.Second)
    after := before
    loop: for i = period; i > 0; i-- {
        select {
        case <-ticker.C:
            maybe, err :=  p.NewStat()
            if err != nil {
                // use the previous value 
                //log.Printf("%s, %d, exited early\n", before.Comm, before.PID)
                break loop
            }
        after = maybe
        }
    }
```
The timer and select make it loop once a second for _period_ seconds.

If it gets a statistic into _maybe_, it assigns it to _after_ and carries on. If 
it doesn't, it breaks out of the loop and uses the value of _after_ 
from the previous second.

## I wonder if it works?
To test I clearly need a program that exits early, let's do something simple:
```bash exit_early
#!/bin/bash
# pretend to run for 30 seconds
sleep 10
# ok, we pretended! now quit.
exit
```
To test it, we uncomment the log.Printf() above and run
```
chmod a+x exit_early
exit_early &
go run main.go -seconds 30 sleep
```

and, sure enough, it replies
```
#name, pid, cputime
2021/05/01 16:06:53 sleep, 23060, exited early
sleep, 23060, 0.000000
```
So we can tell when a process has exited early, we add a new value to the final
printf, so it tells us how many seconds it collected data for.

```go "sample collection" +=
    <<<get before value>>>
    <<<get after value>>>
    fmt.Printf("%s, %d, %d, %f\n", before.Comm, before.PID, 
    period - i, after.CPUTime() - before.CPUTime())
```
Now when we run it against some long-lived processes, it reports
```bash
go run main.go -seconds 5 bash
#name, pid, seconds, cputime
bash, 15009, 5, 0.000000
bash, 12751, 5, 0.000000
bash, 10720, 5, 0.000000
```

and when we run it against exit_early, we get
```bash
exit_early &
go run main.go -seconds 30 sleep
#name, pid, seconds, cputime
sleep, 24156, 9, 0.000000
```
That tells us we caught the last 9 seconds of exit_early's life. 
More about the missing second later!

## What about processes created after we have started to measure?
That's going to cause a change to the main program: it looped through the
passed-in names, getting their process IDs, then started a goroutine for each
pid it found.

