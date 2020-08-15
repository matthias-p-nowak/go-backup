package main

import (
  "log"
  "github.com/matthias-p-nowak/chancloser"
)

// the input channel for the check on cached information
var toCacheChan=make(chan *FileWork,chanLength)

// worker that checks if we have cached information about this file
func toCache(cache *Cache){
  running.Add(1)
  defer running.Done()
  chancloser.Claim(errorWorkChan)
  defer chancloser.Release(errorWorkChan)
  defer log.Println("toCache: done")
  log.Println("toCache: working")
  //
  for entry:=range toCacheChan{
    fd:=new(FileData)
    fd.MTime=entry.MTime
    fd.Size=entry.Size
    fd.Hash=entry.Hash
    log.Println("storing "+entry.Path)
    //err:=cache.Store(entry.Path,fd)
    //if err!=nil {
      //log.Fatal(err)
    //}
  }
} 
