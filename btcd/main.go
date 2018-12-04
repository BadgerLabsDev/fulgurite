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
		args := []string{"--txindex", "--datadir=../../../../media/pi/Picard/chaindata"}
		cmd := exec.Command("./btcd", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		for _, arg := range cmd.Args {
			println(arg)
		}

		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}

		i := 0
		dropcount := 0		// How many seconds has cpu usage been down?
		_, cpuaccum := getCPUUsage()
		ticker := time.NewTicker(1 * time.Second)
		for t := range ticker.C {
			_, cpu   := getCPUUsage()
			cpudiff  := cpu - cpuaccum
			cpuaccum = cpu
			if cpudiff >= 380 {
				dropcount += 1
			}

			if dropcount >= 150 {
				if err := cmd.Process.Kill(); err != nil {
					log.Fatal(err)
				}
				time.Sleep(5 * time.Second) // Give time to btcd to shut down
				goto start
			}
			if (i % 30) == 0 {
				println(t.String())
				println("CPU USAGE:: ", (float32(cpudiff) / 4.0), "%")
			}
			i++
		}
}
