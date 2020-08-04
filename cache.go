package main

import (
  "bytes"
  "encoding/gob"
  "github.com/boltdb/bolt"
  "log"
)

type FileData struct {
  MTime uint64
  Size int64
  Hash []byte
}

type Cache struct {
  Db *bolt.DB
  oldBucket []byte
  newBucket []byte
}

func (cd *Cache)Close(){
  log.Print("closing cache")
  cd.Db.Update(func(tx *bolt.Tx) error{
    return tx.DeleteBucket(cd.oldBucket)
  })
  cd.Db.Close()
}

func OpenCache(fileName string) (cd *Cache){
  cd=new(Cache)
  log.Print("opening")
  db,err := bolt.Open(fileName,0600,nil)
  if err != nil {
    log.Fatal(err.Error())
  }
  cd.Db=db
  newBucket:=0
  err= db.Update(func (tx *bolt.Tx) error {
    for i:=0;i<=1;i++{
      bb:=[]byte{byte(i)}
      b:=tx.Bucket(bb)
      if b==nil {
        newBucket=i
        _,err:=tx.CreateBucket(bb)
        if err!=nil {
          log.Fatal(err)
        }
      }
    }
    return nil
  })
  if newBucket == 0 {
    cd.newBucket = []byte{0}
    cd.oldBucket = []byte{1}
  }else{
    cd.newBucket = []byte{1}
    cd.oldBucket = []byte{0}
  }
  return
}

func (cd *Cache) Retrieve(path string)(fd *FileData, err error){
  fd=new(FileData)
  err=cd.Db.View(func (tx *bolt.Tx) (err error) {
      bb:=tx.Bucket(cd.newBucket)
      val:=bb.Get([]byte(path))
      if val==nil {
        bb:=tx.Bucket(cd.oldBucket)
        val=bb.Get([]byte(path))
      }
      if val!=nil{
        bb:=bytes.NewBuffer(val)
        dec:=gob.NewDecoder(bb)
        err=dec.Decode(fd)
        if err != nil {
          log.Fatal(err)
        }
      }
    return nil
  })
  return
}

func (cd *Cache) Store(path string, fd *FileData) (err error) {
  var bb bytes.Buffer
  enc:=gob.NewEncoder(&bb)
  err=enc.Encode(*fd)
  if err != nil {
    log.Fatal(err.Error())
  }
  err=cd.Db.Update(func (tx *bolt.Tx) (err error){
    bu:=tx.Bucket(cd.newBucket)
    err=bu.Put([]byte(path),bb.Bytes())
    return
  })
  return
}
