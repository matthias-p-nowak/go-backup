package main

import (
	"fmt"
	"log"
)

func t1(cache *Cache){
	// testing
	fd, err := cache.Retrieve("/xxxx")
	if err != nil {
		fmt.Printf("err is %#v\n", err)
	}
	fmt.Printf("fd is %#v\n", fd)
	fd.MTime = 4711
	fd.Size = 42
	fd.Hash = []byte{7, 4, 3, 44, 55, 77}
	err = cache.Store("/xxxx", fd)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func setup(cfg *CFG,cache *Cache){
	wgStarting.Add(1)
	debugSinkChan=CreateFileWorkChan()
	go DebugSink()
	wgStarting.Add(1)
	fromCacheChan=CreateFileWorkChan()
	go FromCache()
	for i:=range cfg.Include {
		wgStarting.Add(1)
		go discover(i,cfg)
	}
}

func main() {
	log.SetFlags(log.LstdFlags|log.Lshortfile)
	log.Print("go-backup started")
	defer log.Print("all done")
	cfg := GetCfg("go-backup.cfg")
	fmt.Printf("%#v\n", cfg)
	cache := OpenCache(cfg.Cache)
	defer cache.Close()
	wgStarting.Add(1)
	setup(cfg,cache)
	wgStarting.Done()
	Wait4Channels()
}
