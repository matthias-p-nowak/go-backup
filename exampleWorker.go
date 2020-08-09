package main

import(
  "github.com/matthias-p-nowak/chancloser"
  "log"
)


// the input channel for the script writer
var exampleChan chan *FileWork=make(chan *FileWork,chanLength)

func example() {
  running.Add(1)
  defer running.Done()
  chancloser.Claim(debugSinkChan)
  defer chancloser.Release(debugSinkChan)
  defer log.Println("example: done")
  log.Println("example: working")
  // setup done
  for entry:= range exampleChan{
    debugSinkChan <- entry
  }
}
