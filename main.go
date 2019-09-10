package main

import (
	"time"
)

func main() {
	run()
	for {
		select {
		case <-time.Tick(config().parserInterval()):
			run() // boom
		}
	}
}

func run() {
	if config().parserIsEnabled() {
		parseRSSFeed()
	}
}
