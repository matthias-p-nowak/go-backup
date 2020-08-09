package main

import(
  "github.com/matthias-p-nowak/chancloser"
  "log"
)


// the input channel for the script writer
var scriptWriterChan chan *FileWork=make(chan *FileWork,chanLength)

func scriptWriter(cfg *CFG) {
  running.Add(1)
  defer running.Done()
  chancloser.Claim(debugSinkChan)
  defer chancloser.Release(debugSinkChan)
  defer log.Println("scriptWriter: done")
  log.Println("scriptWriter: working")
  // setup done
  for entry:=range scriptWriterChan {
    debugSinkChan <- entry
  }
}
