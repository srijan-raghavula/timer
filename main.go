package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const path = "./assets/alarms/sample1.wav"

func main() {
	audio := exec.Command("ffplay", "-nodisp", "-autoexit", path)
	args := os.Args
	scanner := bufio.NewScanner(os.Stdin)
	dStringArg, err := argsParser(args)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}
	for {
		d, err := time.ParseDuration(dStringArg)
		if err != nil {
			fmt.Println("\nerror parsing duration arg: ", err)
			return
		}
		spinners := []string{"|", "/", "-", "\\"}
		for i := int(d.Seconds()); i > 0; i-- {
			spin := spinners[i%len(spinners)]
			fmt.Printf("\r%s Time left: %s", spin, formatSeconds(i))
			os.Stdout.Sync()
			time.Sleep(1 * time.Second)
		}
		go func() {
			err = audio.Run()
			if err != nil {
				fmt.Println("\nError playing file: ", err)
			}
		}()
		fmt.Printf("time > ")
		if !scanner.Scan() {
			fmt.Println("\nError scanning from stdin")
			return
		}
		dStringArg, err = durationParser(strToInt(scanner.Text()))
		if err != nil {
			fmt.Println("\nerror making duration arg: ", err)
			return
		}
	}
}

func argsParser(args []string) (string, error) {
	if len(args) < 2 {
		return "", errors.New("Not enough args")
	}
	return durationParser(strToInt(args[1]))
}

func durationParser(dInt int) (string, error) {
	var seconds, mins, hours int
	seconds = dInt % 100
	if dInt > 99 {
		mins = dInt%10000 - seconds
	}
	if dInt > 9999 {
		hours = dInt%1000000 - mins
	}
	if seconds > 5900 || mins > 5900 || hours > 230000 {
		return "", errors.New("Invalid format")
	}
	dStringArg := fmt.Sprintf("%vh%vm", hours, mins)
	return dStringArg, nil
}

func strToInt(s string) int {
	dInt, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("Error convertin to string: ", err)
		return 0
	}
	return dInt
}

func formatSeconds(s int) string {
	h := s / 3600
	m := (s % 3600) / 60
	sec := s % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, sec)
}
