package main

import(
  "github.com/matthias-p-nowak/chancloser"
  "log"
)


// the input channel for the script writer
var bzip2WriterChan chan *FileWork=make(chan *FileWork,chanLength)

func bzip2Writer(cfg *CFG) {
  running.Add(1)
  defer running.Done()
  chancloser.Claim(debugSinkChan)
  defer chancloser.Release(debugSinkChan)
  defer log.Println("bzip2Writer: done")
  log.Println("bzip2Writer: working")
  // setup done
  for entry:= range bzip2WriterChan{
    debugSinkChan <- entry
  }
}
