/*
Backup program that stores files in bzip2 form under filenames made from the hash sum of the content.
Configuration is read from go-backup.cfg. 
*/
package main

/*
 * TODO: implement a help 
 * TODO: implement command line flags
 */

import (
	"log"
	"runtime"
	"sync"
)

const (
	chanLength int =16
)

// the wait group all goroutines should take part
var running sync.WaitGroup

// dedicated function to set up the infrastructure
func setup(cfg *CFG,cache *Cache){
	workers:=cfg.NumWorkers
	if workers < 1 {
		workers=1
	}
	// worker that walk the files trees and send info out
	for i:=range cfg.Include {
		go discover(i,cfg)
	}
	// worker that retrieves from cache
	go fromCache(cache)
	for i:=0; i< workers;i++{
		go calcHash()
		go checkTarget(cfg)
		go bzip2Writer(cfg)
	}
	// TODO: calculate hash sum
	// TODO: check if target file exists
	// TODO: bzip2 workers - like 114.go
	// TODO: script writer
	go scriptWriter(cfg)
	// TODO: cache writer
	// temp channel to print out the structures
	go debugSink()
}

// Back up files
func main() {
	// setting logs
	log.SetFlags(log.LstdFlags|log.Lshortfile)
	//
	log.Print("go-backup started")
	defer log.Print("all done")
	// configuration
	cfg := GetCfg("go-backup.cfg")
	// open the bolt based cache
	cache := OpenCache(cfg.Cache)
	defer cache.Close()
	// do a setup
	setup(cfg,cache)
	// goroutines are created, but need to be running
	runtime.Gosched()
	// running waits until all goroutines are finished
	running.Wait()
}
