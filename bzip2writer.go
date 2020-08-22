package main

import(
  "github.com/matthias-p-nowak/chancloser"
  "log"
  "os"
  "os/exec"
  "path"
  "path/filepath"
  "syscall"
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
      errorWorkChan <- &Err{entry.Path,err.Error(),E_ERROR}
      return false
    }
    // open/create file
    f,err := os.Create(dest)
    if err !=nil {
      errorWorkChan <- &Err{entry.Path,err.Error(),E_ERROR}
      return false
    }
    defer f.Close()
    // TODO: lock that file
    fd:=int(f.Fd())
    log.Printf("lock fd is %d\n",fd)
    err=syscall.Flock(fd,2)
    if err != nil {
      log.Fatal(err)
    }
    // now bzip2
    cmd:=exec.Command("bzip2","-c",entry.Path)
    cmd.Stdout=f
    err=cmd.Run()
    if err != nil {
      errorWorkChan <- &Err{ entry.Path, "bzip2 failed", E_ERROR}
      return false
    }
    stat,err := os.Lstat(entry.Path)
    if err != nil {
      errorWorkChan <- &Err{ entry.Path, "lstat failed", E_ERROR}
      return false
    }
    if stat.ModTime().Unix() != entry.MTime {
      errorWorkChan <- &Err{ entry.Path, "file time changed during bzip2", E_WARNING}
      err:=os.Remove(dest)
      if err != nil {
        errorWorkChan <- &Err{ entry.Path, "removal of the file "+dest+" failed", E_ERROR}
      }
      return false
    }
    if stat.Size() != entry.Size {
      errorWorkChan <- &Err{ entry.Path, "size changed during bzip2", E_WARNING}
      err:=os.Remove(dest)
      if err != nil {
        errorWorkChan <- &Err{ entry.Path, "removal of the file "+dest+" failed", E_ERROR}
      }
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
      scriptWriterChan <- entry
    } else {
      errorWorkChan <- &Err{entry.Path,"storing failed",E_ERROR}
    }
  }
  log.Printf("done: %d\n",worked)
}
