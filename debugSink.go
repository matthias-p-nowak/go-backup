package main

import (
  "fmt"
  "log"
)

var debugSinkChan *FileWorkChan

func DebugSink(){
  // only a receiver
  wgStarting.Done()
  wgStarting.Wait()
  // setup done
  defer log.Print("DebugSink ended")
  for fw:=range debugSinkChan.ch {
    fmt.Printf("%#v\n",fw)
  }
}
