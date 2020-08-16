package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

/*
Backup runs in addition to usual programs and should occupy the processor.
/proc/stat shows the processor loads and this ensures that 1/8th of the processor is idle.
 */

var workTickets = make(chan int)

func getFields() []string {
	// there is no easier way than reading the pseudo file
	content, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(content), "\n")
	return strings.Fields(lines[0])
}

func checkLoad() {
  log.Print("started")
	ticks := time.Tick(100 * time.Millisecond)
	var got [5]uint64
	var sum uint64
	var idle uint64
	// /proc/stat increases counter, therefore we need the first value
  fields:=getFields()
	for f := 0; f < 5; f++ {
		v, err := strconv.ParseUint(fields[f+1], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		got[f] = v
	}
	// this function never ends to run
	lastPrint:=0
	for {
		select {
		case <-ticks:
			for {
				lastPrint++
				fields = getFields()
				sum -= sum >> 5
				idle -= idle >> 5
				for f := 0; f < 5; f++ {
					v, err := strconv.ParseUint(fields[f+1], 10, 64)
					if err != nil {
						log.Fatal(err)
					}
					d := v - got[f]
					got[f] = v
					sum += d
					if f == 3 {
						idle += d
					}
				}
				if sum>>4 < idle {
					break
				}
				if lastPrint > 50{
					lastPrint=0
					log.Println("checkLoad: not enough idle")
				}
				<-ticks
			}
		case workTickets <- 1:
		}
	}
}

func init(){
	go checkLoad()
}
