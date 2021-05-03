try: build
	go run main.go -seconds 5 bash

build:
	lmt README.md Implementation.md
	go fmt


early:  build
	exit_early &
	go run main.go -seconds 30 sleep

late: build
	parent &
	go run main.go -seconds 30 sleep
