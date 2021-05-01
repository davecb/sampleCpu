try: build
	go run main.go -seconds 5 bash

build:
	lmt README.md Mainline.md
	go fmt


chmod:
	chmod a+x exit_early
fail: chmod build
	exit_early &
	go run main.go -seconds 30 sleep
ps: chmod
	exit_early &
	pgrep sleep
