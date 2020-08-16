package main

import(
  "github.com/matthias-p-nowak/chancloser"
  "log"
  "os"
  "os/exec"
  "path"
  "path/filepath"
)


// the input channel for the script writer
var bzip2WriterChan=make(chan *FileWork,chanLength)

// defer only works inside a function, therefore a function for calling the bzip2 program
func runBzip2(entry *FileWork, dest string)(succ bool){
    // exec works excellent
    dir,_:=filepath.Split(dest)
    // make directories
    err:=os.MkdirAll(dir,0755)
    if err != nil {
      entry.record(err.Error())
      errorWorkChan <- entry
      return false
    }
    // open/create file
    f,err := os.Create(dest)
    if err !=nil {
      entry.record(err.Error())
      errorWorkChan <- entry
      return false
    }
    defer f.Close()
    // now bzip2
    cmd:=exec.Command("bzip2","-c",entry.Path)
    cmd.Stdout=f
    err=cmd.Run()
    if err != nil {
      entry.record("executing bzip2 failed")
      return false
    }
    stat,err := os.Lstat(entry.Path)
    if err != nil {
      entry.record("can't get lstat info")
      return false
    }
    if stat.ModTime().Unix() != entry.MTime {
      entry.record("file time changed during bzip2")
      return false
    }
    if stat.Size() != entry.Size {
      entry.record("size changed during bzip2")
      return false
    }
    return true
}

// worker stores files as bzip2 under hashsum made up filename
func bzip2Writer(cfg *CFG) {
  // 
  running.Add(1)
  defer running.Done()
  //
  chancloser.Claim(scriptWriterChan)
  defer chancloser.Release(scriptWriterChan)
  chancloser.Claim(errorWorkChan)
  defer chancloser.Release(errorWorkChan)
  // 
  // log.Println(" working")
  // setup done
  worked:=0
  for entry:= range bzip2WriterChan{
    <- workTickets
    worked++
    // work
    dest:=path.Join(cfg.Destination,"f",entry.Hash)
    if runBzip2(entry,dest) {
      entry.record("wrote bzip2")
      scriptWriterChan <- entry
    } else {
      entry.record("storing failed")
      errorWorkChan <- entry
    }
  }
  log.Printf("done: %d\n",worked)
}
