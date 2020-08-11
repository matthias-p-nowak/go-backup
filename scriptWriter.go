package main

import(
  "github.com/matthias-p-nowak/chancloser"
  "log"
)

// the input channel for the script writer
var scriptWriterChan chan *FileWork=make(chan *FileWork,chanLength)

/*
Writes a small script and then lines, each one related to one stored file  or a no-content element of the file system.
 */
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
