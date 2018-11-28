package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Only works on linux
func getCPUUsage() (total, idle uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range(lines) {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return
}


func main() {
	start:
		cmd := exec.Command("/go/bin/btcd", "—txindex", "—datadir=pi/media/Untiled/chaindata")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}

		ticker := time.NewTicker(1 * time.Second)
		for t := range ticker.C {
			_, cpu := getCPUUsage()
			if cpu <= 5 {
				if err := cmd.Process.Kill(); err != nil {
					log.Fatal(err)
				}
				time.Sleep(5 * time.Second) // Give time to btcd to shut down
				goto start
			}
			println(t.String())
		}
}