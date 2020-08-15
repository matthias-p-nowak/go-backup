package main

import (
  "log"
  "github.com/matthias-p-nowak/chancloser"
)

// the input channel for the check on cached information
var fromCacheChan=make(chan *FileWork,chanLength)

// worker that checks if we have cached information about this file
func fromCache(cache *Cache){
  running.Add(1)
  defer running.Done()
  chancloser.Claim(calcHashChan)
  defer chancloser.Release(calcHashChan)  
  chancloser.Claim(checkTargetChan)
  defer chancloser.Release(checkTargetChan)
  chancloser.Claim(errorWorkChan)
  defer chancloser.Release(errorWorkChan)
  defer log.Println("fromCache: done")
  log.Println("fromCache: working")
  // setup done
  // #############################################
  for entry:=range fromCacheChan {
    <-workTickets
    if entry.Size <= 0 {
      // should be bigger than 0 - if wrong, this is a programming error
      log.Fatal("wrong programming assumption, size <0:",entry)
    }
    // avoid hash calculation if we know it from before
    fd,err:=cache.Retrieve(entry.Path)
    doHash := err == cacheEmpty // don't know
    if (err != nil) && ! doHash {
      // there was an error, but not an empty cache
      entry.record("Cache error:"+err.Error())
      errorWorkChan <- entry
      continue
    }
    // is it the same as remembered, same size and mtime? 
    doHash = doHash || fd.Size != entry.Size
    doHash = doHash || fd.MTime != entry.MTime
      // TODO is it still the same file?
    if doHash {
      entry.record("do a hash calculation")
      calcHashChan <- entry
    } else {
      checkTargetChan <- entry
    }
  }
}
