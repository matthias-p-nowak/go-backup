package main

import(
  "log"
)


// the input channel for the script writer
var errorWorkChan=make(chan *FileWork,chanLength)

type Err struct{
  path string
  cause string
}

// List of errors during the program run, can be a lot...
var allErrors []Err

func errorWork(cfg *CFG) {
  running.Add(1)
  defer running.Done()
  defer log.Println("errorWork: done")
  log.Println("errorWork: working")
  // setup done
  for entry:= range errorWorkChan{
    e:=new(Err)
    e.path=entry.Path
    l:=len(entry.workDone)
    if l> cfg.MaxErrors {
      continue
    }
    if l > 0 {
      e.cause=entry.workDone[l-1]
    }
    allErrors=append(allErrors,*e)
  }
}
