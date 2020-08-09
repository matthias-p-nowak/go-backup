
package main

import (
  "log"
  "os"
  "path"
  "path/filepath"
  "syscall"
  "github.com/matthias-p-nowak/chancloser"
)

// discovers files to backup and sends info onto channels
func discover(f int, cfg *CFG){
  running.Add(1)
  defer running.Done()
  chancloser.Claim(fromCacheChan)
  defer chancloser.Release(fromCacheChan)
  chancloser.Claim(scriptWriterChan)
  defer chancloser.Release(scriptWriterChan)
  // setup done
  // how old do file have to be?
  cutOff:=cfg.AgeCutoff()
  // can start doing work now
  path2walk:=cfg.Include[f]
  defer log.Println("discover: done",path2walk)
  log.Println("discover: walking ",path2walk)
  // which device is path2walk on?
  st,err:=os.Lstat(path2walk)
  if err != nil {
    log.Fatal(err)
  }
  s:=st.Sys()
  ss:=s.(*syscall.Stat_t)
  // storing the device number
  dev:=ss.Dev
  // the work is done by this function
  fwf:=func(fpath string, info os.FileInfo, err error) (error){
    // log.Print("looking at ",p)
    // check if this directory should not be backed up
    if info.IsDir(){
      p2:=path.Join(fpath,".nobackup")
      _,e:=os.Lstat(p2)
      if e == nil {
        log.Print("not backing up dir ",fpath)
        // cutting off subtree
        return filepath.SkipDir
      } 
    }
    // checking the device of this file
    s:=info.Sys()
    stat_t:=s.(*syscall.Stat_t)
    if dev != stat_t.Dev {
      log.Print("skipping ",fpath," on different device")
      if info.IsDir() {
        // cutting off subtree
        return filepath.SkipDir
      } else {
        return nil
      }
    }
    // does the file match an exclusion pattern?
    b:=path.Base(fpath)
    for _,ex:=range cfg.Exclude {
      m,e:=path.Match(ex,b)
      if e!=nil {
        log.Fatal(e)
      }
      if m {
        log.Print("skipping ",fpath," ",b)
        // special treatment for cutting off subtree
        if info.IsDir(){
          return filepath.SkipDir
        } else {
          return nil
        }
      }
    }
    // don't backup files that are too new
    if info.ModTime().Unix() >= cutOff {
      return nil
    }
    // sending it to the next channel
    entry:=new(FileWork)
    entry.Path=fpath
    entry.FileInfo=info
    entry.MTime=info.ModTime().Unix()
    entry.Size=info.Size()      
    entry.record("walked "+path2walk)
    if info.Mode().IsRegular(){
      fromCacheChan <- entry            
    } else{
      scriptWriterChan <- entry      
    }
    return nil
  }
  // commence the work
  log.Print("walking ",path2walk)
  filepath.Walk(path2walk,fwf)
  log.Print("walked ",path2walk)
}
