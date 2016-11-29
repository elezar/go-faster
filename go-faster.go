package main

import (
	"log"

	"time"

	fast "gopkg.in/ddo/go-fast.v0"
)

func process(speeds <-chan float64, errors <-chan error) {
	for {
		select {
		case mbps := <-speeds:
			log.Printf("SPEED (Mbps): %.2f\n", mbps)
		case err := <-errors:
			log.Println("ERROR: ", err)
		}
	}
}

func runTest(speeds chan<- float64, errors chan<- error) {
	fastCom := fast.New()

	err := fastCom.Init()
	if err != nil {
		errors <- err
	}

	// get urls
	urls, err := fastCom.GetUrls()
	if err != nil {
		errors <- err
	}

	// measure
	KbpsChan := make(chan float64)
	go func() {
		for Kbps := range KbpsChan {
			speeds <- Kbps / 1024
		}
	}()

	err = fastCom.Measure(urls, KbpsChan)
	if err != nil {
		errors <- err
	}

}

func main() {
	errors := make(chan error)
	speeds := make(chan float64)
	force := make(chan bool)

	period := 1 * time.Minute

	go func() {
		force <- true
		process(speeds, errors)
	}()

	for {
		run := false
		select {
		case <-force:
			run = true
		case <-time.After(period):
			run = true
		}
		if run {
			runTest(speeds, errors)
		}
	}

}
