package main

import (
  "log"
  "os"
  "path"
  "path/filepath"
  "syscall"
)

func discover(f int, cfg *CFG){
   fromCacheChan.wg.Add(1)
  defer fromCacheChan.wg.Done()
  wgStarting.Done()
  p:=cfg.Include[f]
  st,err:=os.Lstat(p)
  if err != nil {
    log.Fatal(err)
  }
  s:=st.Sys()
  ss:=s.(*syscall.Stat_t)
  dev:=ss.Dev
  fwf:=func(p string, info os.FileInfo, err error) (error){
    // log.Print("looking at ",p)
    if info.IsDir(){
      p2:=path.Join(p,".nobackup")
      _,e:=os.Lstat(p2)
      if e == nil {
        log.Print("not backing up dir ",p)
        return filepath.SkipDir
      } 
    }
    s:=info.Sys()
    ss:=s.(*syscall.Stat_t)
    if dev != ss.Dev {
      log.Print("skipping ",p," on different device")
      if info.IsDir() {
        return filepath.SkipDir
      } else {
        return nil
      }
    }
    b:=path.Base(p)
    for _,ex:=range cfg.Exclude {
      m,e:=path.Match(ex,b)
      if e!=nil {
        log.Fatal(e)
      }
      if m {
        log.Print("skipping ",p," ",b)
        if info.IsDir(){
          return filepath.SkipDir
        } else {
          return nil
        }
      }
    }
    fw:=new(FileWork)
    fw.Path=p
    fw.FileInfo=info
    fromCacheChan.ch <- fw
    return nil
  }
  log.Print("walking ",p)
  filepath.Walk(p,fwf)
}
