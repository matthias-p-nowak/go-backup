
package main

import (
  "log"
  "os"
  "sync"
)

var wgChannels sync.WaitGroup
var wgStarting sync.WaitGroup

type FileWork struct {
  Path string
  Hashsum string
  MTime uint64
  Size uint64
  Uid int
  Gid int
  Mode string
  FileInfo os.FileInfo
}

type FileWorkChan struct {
  ch chan *FileWork
  wg sync.WaitGroup
}

func (ac *FileWorkChan) autoClose(){
  wgStarting.Wait()
  ac.wg.Wait()
  log.Print("closing auto channel",ac.ch)
  close(ac.ch)
  wgChannels.Done()
}


func CreateFileWorkChan()(ac *FileWorkChan){
  ac=new(FileWorkChan)
  ac.ch=make(chan *FileWork,512)
  wgChannels.Add(1)
  go ac.autoClose()
  return
}

func Wait4Channels(){
  wgChannels.Wait()
}
