package main

import (
  "log"
  "github.com/matthias-p-nowak/chancloser"
)

// the input channel for the check on cached information
var fromCacheChan chan *FileWork=make(chan *FileWork,chanLength)

// worker that checks if we have cached information about this file
func fromCache(cache *Cache){
  running.Add(1)
  defer running.Done()
  chancloser.Claim(calcHashChan)
  defer chancloser.Release(calcHashChan)  
  chancloser.Claim(checkTargetChan)
  defer chancloser.Release(checkTargetChan)
  defer log.Println("fromCache: done")
  log.Println("fromCache: working")
  // setup done
  // #############################################
  for entry:=range fromCacheChan {
    if entry.Size < 0 {
      log.Fatal("size <0:",entry)
    }
    fd,err:=cache.Retrieve(entry.Path)
    doHash := err == cacheEmpty
    doHash = doHash || fd.Size != entry.Size
    doHash = doHash || fd.MTime != entry.MTime
      // TODO is it still the same file?
    if doHash {
      calcHashChan <- entry
    }else{
      checkTargetChan <- entry
    }
  }
}
