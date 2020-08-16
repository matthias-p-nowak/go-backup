package main

import(
  "github.com/matthias-p-nowak/chancloser"
  "log"
  "os"
  "path"
  "time"
)


// the input channel for the script writer
var checkTargetChan chan *FileWork=make(chan *FileWork,chanLength)

// worker that checks if the backup target exists
// if it exists, we can skip the bzip2 step
func checkTarget(cfg *CFG) {
  // running
  running.Add(1)
  defer running.Done()
  // channels
  chancloser.Claim(bzip2WriterChan)
  defer chancloser.Release(bzip2WriterChan)
  chancloser.Claim(scriptWriterChan)
  defer chancloser.Release(scriptWriterChan)
  //
  // setup done
  dest:=cfg.Destination
  worked:=0
  for entry:= range checkTargetChan{
    worked++
    filepath:=path.Join(dest,"f",entry.Hash)
    _,err := os.Lstat(filepath)
    if os.IsNotExist(err){
      // write bzip2 file
      bzip2WriterChan <- entry
    } else {
      // existing target needs new timestamp for keepfree to work
      // keepfree deletes the oldest (mtime) files to make space
      currentTime := time.Now().Local()
      err = os.Chtimes(filepath, currentTime, currentTime)
      if err != nil {
        // this might indicate a bigger problem, therefore fatal log
        log.Fatal(err)
        }
      entry.record("touched destination file")
      scriptWriterChan <-entry
    }
  }
  log.Printf("done: %d\n",worked)
}
