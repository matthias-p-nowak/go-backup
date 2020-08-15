package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

var workTickets = make(chan int)

func getFields() []string {
	content, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(content), "\n")
	return strings.Fields(lines[0])
}

func checkLoad() {
  log.Print("load check started")
  defer log.Print("load check stopped")
	ticks := time.Tick(100 * time.Millisecond)
	var got [5]uint64
	var sum uint64
	var idle uint64
  fields:=getFields()
	for f := 0; f < 5; f++ {
		v, err := strconv.ParseUint(fields[f+1], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		got[f] = v
	}
	for {
		select {
		case <-ticks:
			for {
				fields = getFields()
				sum -= sum >> 4
				idle -= idle >> 4
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
				// fmt.Printf("%d - %d %v \n", sum, idle, fields)
				if sum>>3 < idle {
					break
				}
        /// fmt.Print(".")
				<-ticks
			}
		case workTickets <- 1:
		}
	}
}

func init(){
	go checkLoad()
}
