package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

// runs the given command and writes the stdout to the given outputPath
func runner(cmdString, outputPath string) error {
	cmdFields := strings.Fields(cmdString)
	cmdName := cmdFields[0]
	cmdParams := cmdFields[1:]

	f, err := os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file: ", outputPath)
		fmt.Println(err)
		return err
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("### Starting:")
	defer log.Println("### Done:")
	cmd := exec.Command(cmdName, cmdParams...)
	output := io.MultiWriter(f, os.Stdout)
	cmd.Stdout = output
	cmd.Stderr = output
	err = cmd.Run()
	if err != nil {
		fmt.Fprintln(output, "Error running: ", cmdString)
		fmt.Fprintln(output, err, cmdString)
	}

	return nil

}

// executes the given command periodically and writes the stdout to the specified file
func main() {
	command := flag.String("c", "speedtest-cli", "Command that will be exectued.")
	period := flag.Int("p", 600, "Period in which the command will be executed in seconds.")
	outputPath := flag.String("o", "go-faster.log", "File in which the data should go.")

	flag.Parse()

	*period = *period * int(time.Second)

	err := runner(*command, *outputPath)
	if err != nil {
		return
	}

	for _ = range time.Tick(time.Duration(*period)) {
		runner(*command, *outputPath)
	}
}
