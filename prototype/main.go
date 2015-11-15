package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

func runner(what string) {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic("NOOOOO")
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("### Starting:")
	cmd := exec.Command(what)
	cmd.Stdout = io.MultiWriter(f, os.Stdout)
	cmd.Stderr = os.Stderr
	cmd.Run()
	log.Println("### Done:")

}

func main() {
	runner("speedtest-cli")

	period := 10 * time.Minute
	for _ = range time.Tick(period) {
		runner("speedtest-cli")
	}

}
