package main


// Caching the hash sum for files that still have the same ModTime and the size
 
import (
  "bytes"
  "encoding/gob"
  "github.com/boltdb/bolt"
  "log"
)

// Info about files that is cached, key is file path
type FileData struct {
  MTime int64 // modification time
  Size int64 // size
  Hash string // determined hash
}


// Bolt offers buckets, we use one for reading older data, writing to the new bucket
type Cache struct {
  Db *bolt.DB
  oldBucket []byte
  newBucket []byte
}

// closing the mmap bolt database, deleting the old bucket before closing
func (cd *Cache)Close(){
  log.Print("closing cache")
  cd.Db.Update(func(tx *bolt.Tx) error{
    return tx.DeleteBucket(cd.oldBucket)
  })
  cd.Db.Close()
}

// Opening a bolt database and creating the buckets
func OpenCache(fileName string) (cd *Cache){
  cd=new(Cache)
  log.Print("opening")
  db,err := bolt.Open(fileName,0600,nil)
  if err != nil {
    log.Fatal(err.Error())
  }
  cd.Db=db
  // which one is the new bucket?
  newBucket:=0
  // creating the buckets and recording which is the new one
  err= db.Update(func (tx *bolt.Tx) error {
    for i:=0;i<=1;i++{
      bb:=[]byte{byte(i)}
      b:=tx.Bucket(bb)
      if b==nil {
        // bucket did not exist
        newBucket=i
        _,err:=tx.CreateBucket(bb)
        if err!=nil {
          log.Fatal(err)
        }
      }
    }
    return nil
  })
  // recording correct bucket id
  if newBucket == 0 {
    cd.newBucket = []byte{0}
    cd.oldBucket = []byte{1}
  }else{
    cd.newBucket = []byte{1}
    cd.oldBucket = []byte{0}
  }
  return
}

// specialized Error structure
type CacheEmpty struct { 
  error // embedding error interface
}
var cacheEmpty CacheEmpty
func (c CacheEmpty) Error() string {
  return "No cache entry"
}

// Retrieves cached data for a path if found, otherwise returning cacheEmpty
func (cd *Cache) Retrieve(path string)(fd *FileData, err error){
  fd=new(FileData) // for return value
  err=cd.Db.View(func (tx *bolt.Tx) (err error) {
    // inside closure
    // getting from new bucket
    bb:=tx.Bucket(cd.newBucket)
    val:=bb.Get([]byte(path))
    if val==nil {
      // wasn't in new bucket
      bb:=tx.Bucket(cd.oldBucket)
      val=bb.Get([]byte(path))
    }
    if val!=nil{
      // got something, need to decode it
      bb:=bytes.NewBuffer(val)
      dec:=gob.NewDecoder(bb)
      err=dec.Decode(fd)
      if err != nil {
        // should have happened
        log.Fatal(err)
      }
      // the decoded is in fd
    } else {
      return cacheEmpty
    }
    return nil // no error when retrieving
  })
  // returning named results fd,err
  return
}

// Storing a gob encoded data in the new bucket
func (cd *Cache) Store(path string, fd *FileData) (err error) {
  var bb bytes.Buffer
  enc:=gob.NewEncoder(&bb)
  err=enc.Encode(*fd)
  if err != nil {
    log.Fatal(err.Error())
  }
  // do the update
  return cd.Db.Update(func (tx *bolt.Tx) (err error){
    bu:=tx.Bucket(cd.newBucket)
    err=bu.Put([]byte(path),bb.Bytes())
    return
  })
}
