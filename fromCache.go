package main

import (
  // "fmt"
  "log"
)

var fromCacheChan *FileWorkChan

func FromCache(){
  // only a receiver
  wgStarting.Done()
  // setup done
  defer log.Print("DebugSink ended")
  for fw:=range fromCacheChan.ch {
    fw=fw
  }
}
