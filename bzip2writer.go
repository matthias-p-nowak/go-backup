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

func runBzip2(src string, dest string){
    cmd:=exec.Command("bzip2","-c",src)
    dir,_:=filepath.Split(dest)
    err:=os.MkdirAll(dir,0755)
    if err != nil {
      log.Fatal(err)
    }
    f,err := os.Create(dest)
    if err !=nil {
      log.Fatal(err)
    }
    defer f.Close()
    cmd.Stdout=f
    cmd.Run()
}

func bzip2Writer(cfg *CFG) {
  running.Add(1)
  defer running.Done()
  chancloser.Claim(scriptWriterChan)
  defer chancloser.Release(scriptWriterChan)
  defer log.Println("bzip2Writer: done")
  log.Println("bzip2Writer: working")
  // setup done
  for entry:= range bzip2WriterChan{
    dest:=path.Join(cfg.Destination,entry.Hash)
    runBzip2(entry.Path,dest)
    entry.record("wrote bzip2")
    scriptWriterChan <- entry
  }
}
