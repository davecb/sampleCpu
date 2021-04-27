sampleProc.1
====
## NAME
sampleProc - Capture a sample of activity from a program
## SYNOPSIS
sampleProc [--seconds N] program [program] ...
## DESCRIPTION
sampleProc is used to capture samples of CPU activity from
programs via /proc.

Output is in comma-separated-value format

### --seconds N
Sets the length of the sample in seconds.

## BUGS
This is inherently a "racy" process. Programs can exit before the
sample ends, and new ones can come into existance just as easily.

Be quite careful when comparing multithereaded programs and ones
which fork/exec children. This programs was written to help
with just that case, but it's still not easy.



