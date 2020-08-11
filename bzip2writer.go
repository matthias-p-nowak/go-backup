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
var bzip2WriterChan chan *FileWork=make(chan *FileWork,chanLength)

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
    cmd.Run()
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
  defer log.Println("bzip2Writer: done")
  log.Println("bzip2Writer: working")
  // setup done
  for entry:= range bzip2WriterChan{
    // work
    dest:=path.Join(cfg.Destination,"f",entry.Hash)
    if runBzip2(entry,dest) {
      entry.record("wrote bzip2")
      scriptWriterChan <- entry
    }
  }
}
