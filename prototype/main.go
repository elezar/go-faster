package main

import (
	"fmt"
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
	defer log.Println("### Done:")

	cmd := exec.Command(what)

	output := io.MultiWriter(f, os.Stdout)
	cmd.Stdout = output
	cmd.Stderr = output

	err = cmd.Run()
	if err != nil {
		fmt.Fprintln(output, "Error running: ", what)
		fmt.Fprintln(output, err)
	}

}

func main() {
	runner("speedtest-cli")

	period := 10 * time.Minute
	for _ = range time.Tick(period) {
		runner("speedtest-cli")
	}

}
