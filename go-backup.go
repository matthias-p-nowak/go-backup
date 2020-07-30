package main

import (
  "fmt"
)

func main(){
  fmt.Println("go-backup started")
  defer fmt.Println("all done")
  cfg := GetCfg("go-backup.cfg")
  fmt.Printf("%#v\n",cfg)
  cache:=OpenCache(cfg.Cache)
  defer cache.Close()
  // testing
  fd,err:=cache.Retrieve("/xxxx")
  if err != nil {
    fmt.Printf("err is %#v\n",err)
  }
  fmt.Printf("fd is %#v\n",fd)
  fd.MTime=4711
  fd.Size=42
  fd.Hash=[]byte{7,4,3,44,55,77}
  err=cache.Store("/xxxx",fd)
  if err != nil {
    fmt.Println(err.Error())
  }
}
