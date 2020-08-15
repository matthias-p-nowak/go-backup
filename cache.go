package main


// Caching the hash sum for files that still have the same ModTime and the size
 
import (
  "bytes"
  "encoding/gob"
  "github.com/syndtr/goleveldb/leveldb"
  "log"
  "os"
)

// Info about files that is cached, key is file path
type FileData struct {
  MTime int64 // modification time
  Size int64 // size
  Hash string // determined hash
}

type Cache struct {
  DbOld *leveldb.DB
  DbNew *leveldb.DB
  fileNameOld string
  fileNameNew string
}

func (cd *Cache)Close(){
  log.Println("closing cache")
  cd.DbOld.Close()
  cd.DbNew.Close()
  err:=os.RemoveAll(cd.fileNameOld)
  if err != nil {
    log.Fatal(err.Error())
  }
  err=os.Rename(cd.fileNameNew,cd.fileNameOld)
  if err != nil {
    log.Fatal(err.Error())
  }
}

func OpenCache(fileName string) (cd *Cache){
  cd=new(Cache)
  log.Print("opening")
  cd.fileNameOld=fileName
  cd.fileNameNew=fileName+".new"
  var err error
  cd.DbOld,err=leveldb.OpenFile(cd.fileNameOld,nil)
  if err != nil {
    log.Fatal(err.Error())
  }
  cd.DbNew,err=leveldb.OpenFile(cd.fileNameNew,nil)
  if err != nil {
    log.Fatal(err.Error())
  }
  return
}


// specialized Error structure
type CacheEmpty struct { 
  error // embedding error interface
}

func (cd *Cache) Retrieve(filename string)(fd *FileData, err error){
  fd=new(FileData)
  val,err:=cd.DbNew.Get([]byte(filename),nil)
  if err == leveldb.ErrNotFound {
    val,err = cd.DbOld.Get([]byte(filename),nil)
    if err == leveldb.ErrNotFound {
      return
    }
  }
  // fmt.Printf("val is %#v\n",val)
  bb:=bytes.NewBuffer(val)
  dec:=gob.NewDecoder(bb)
  err=dec.Decode(fd)
  if err != nil {
    log.Fatal(err.Error())
  }
  return
}

func (cd *Cache) Store(filename string, fd *FileData) (err error) {
  var bb bytes.Buffer
  enc:=gob.NewEncoder(&bb)
  err=enc.Encode(*fd)
  if err != nil {
    log.Fatal(err.Error())
  }
  err=cd.DbNew.Put([]byte(filename),bb.Bytes(),nil)
  if err != nil {
    log.Fatal(err.Error())
  }
  return
}


