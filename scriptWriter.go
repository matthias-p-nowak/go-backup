package main

import(
  "github.com/matthias-p-nowak/chancloser"
  "fmt"
  "io"
  "log"
  "os"
  "path"
  "path/filepath"
  "strings"
  "syscall"
  "time"
)

// the input channel for the script writer
var scriptWriterChan chan *FileWork=make(chan *FileWork,chanLength)

func sanitizePath(f string)(r string){
  r=strings.ReplaceAll(f,`"`,`\"`)
  r=strings.ReplaceAll(r,"`","\\`")
  return
}

/*
Writes a small script and then lines, each one related to one stored file  or a no-content element of the file system.
 */
func scriptWriter(cfg *CFG) {
  running.Add(1)
  defer running.Done()
  chancloser.Claim(errorWorkChan)
  defer chancloser.Release(errorWorkChan)
  chancloser.Claim(toCacheChan)
  defer chancloser.Release(toCacheChan)
  // setup done
  now:=time.Now()
  hostname,err:=os.Hostname()
  if err != nil {
    log.Fatal(err)
  }
  nowStr:=now.Format("2006-01-02_15-04-05")
  dest:=path.Join(cfg.Destination,"s",hostname,nowStr+".sh")
  dir,_:=filepath.Split(dest)
  // make directories
  err=os.MkdirAll(dir,0755)
  if err != nil {
    log.Fatal("can't screate script directory "+dir, err)
  }
  script,err := os.Create(dest)
  if err !=nil {
    log.Fatal(err)
  }
  defer script.Close()
  
  log.Println("writing "+dest)
  str:=`
#!/bin/bash
BACKUP=${BACKUP:-%s}
DEST=${DEST:-/}
`
  str=fmt.Sprintf(str,cfg.Destination)
  script.WriteString(str)
  io.Copy(script, GetStored("snippets/restore.sh"))
  entryCnt:=0
  // xx mode owner hash path
  for entry:=range scriptWriterChan {
    entryCnt++
    modeBits:=entry.FileInfo.Mode()
    mm:=modeBits.Perm()
    if (modeBits & os.ModeSetuid) >0 {
      mm+=04000
    }
    if (modeBits & os.ModeSetgid) >0 {
      mm+=02000
    }
    if (modeBits & os.ModeSticky) >0 {
      mm+=01000
    }
    mode:=fmt.Sprintf("0%o",mm)
    s:=entry.FileInfo.Sys()
    ss:=s.(*syscall.Stat_t)
    user:=fmt.Sprintf("%d:%d",ss.Uid,ss.Gid)
    switch {
      case len(entry.Hash)>0:
        script.WriteString("f "+user+" "+mode+" "+entry.Hash+" "+sanitizePath(entry.Path) +"\n")
        toCacheChan <- entry
      case entry.FileInfo.Mode().IsRegular():
        script.WriteString("e "+user+" "+mode+""+sanitizePath(entry.Path)+"\n")
      case entry.FileInfo.IsDir():
        script.WriteString("d "+user+" "+mode+" "+sanitizePath(entry.Path)+"\n")
      case (modeBits & os. ModeSymlink) >0:
        l,err:=os.Readlink(entry.Path)
        if err!=nil {
          entry.record("couldn't read symlink")
          errorWorkChan <- entry
          continue
        }
        script.WriteString("s "+user+" "+sanitizePath(entry.Path)+ " "+l+"\n")
      case (modeBits & os.ModeNamedPipe) >0:
        script.WriteString("p "+user+" "+mode+" "+sanitizePath(entry.Path)+"\n")
      default:
        script.WriteString("# couldn't deal with"+entry.Path+"\n")
    }
  }
  script.WriteString(fmt.Sprintf("echo all done\nexit\n##### ##### #####\n%d\n",entryCnt))
  log.Printf("done, entries: %d\n",entryCnt)
}
