package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"time"
	"flag"
	"fmt"
	"strings"
	"regexp"
	"strconv"
)

// runs the given command and writes the stdout to the given outputPath
func runner(cmdString, outputPath string) {
	cmdFields := strings.Fields(cmdString)
	cmdName := cmdFields[0]
	cmdParams := cmdFields[1:]

	f, err := os.OpenFile(outputPath, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		panic("NOOOOO")
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

}

// executes the given command periodically and writes the stdout to the specified file
func main() {
	command := flag.String("c", "speedtest-cli", "Command that will be exectued.")
	period := flag.String("p", "10m", "Period in which the command will be executed in.")
	outputPath := flag.String("o", "go-faster.log", "File in which the data should go.")

	flag.Parse()

	properPeriodRegex, _ := regexp.Compile("\\d+[s,m,h,d]")

	if ! properPeriodRegex.MatchString(*period) {
		panic("period has to match '\\d+[s,m,h,d]'")
	}

	periodUnitRegex, _ := regexp.Compile("\\d+")
	periodValueRegex, _ := regexp.Compile("[smhd]")

	periodValue, _ := strconv.ParseInt(periodValueRegex.ReplaceAllString(*period, ""),10,64)
	periodUnit := string(periodUnitRegex.ReplaceAllString(*period, ""))

	var periodDuration int64

	switch periodUnit {
	case "s" : periodDuration = periodValue * int64(time.Second)
	case "m" : periodDuration = periodValue * int64(time.Minute)
	case "h" : periodDuration = periodValue * int64(time.Hour)
	case "d" : periodDuration = periodValue * int64(time.Hour) * 24
	}

		runner(*command, *outputPath)

	for _ = range time.Tick(time.Duration(periodDuration)) {
		runner(*command, *outputPath)
	}
}