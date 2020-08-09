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

func checkTarget(cfg *CFG) {
  running.Add(1)
  defer running.Done()
  chancloser.Claim(bzip2WriterChan)
  defer chancloser.Release(bzip2WriterChan)
  chancloser.Claim(scriptWriterChan)
  defer chancloser.Release(scriptWriterChan)
  defer log.Println("checkTarget: done")
  log.Println("checkTarget: starting")
  // setup done
  dest:=cfg.Destination
  for entry:= range checkTargetChan{
    filepath:=path.Join(dest,entry.Path)
    _,err := os.Lstat(filepath)
    exists:=! os.IsNotExist(err)
    if exists {
      currentTime := time.Now().Local()
      err = os.Chtimes(filepath, currentTime, currentTime)
      if err != nil {
        log.Fatal(err)
        }
      entry.record("touched destination file")
      scriptWriterChan <-entry
    }else{
      // write bzip2 file
      bzip2WriterChan <- entry
    }
  }
}
