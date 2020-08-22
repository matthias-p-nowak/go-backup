package main

import(
  "github.com/matthias-p-nowak/chancloser"
  "log"
  "io"
  "crypto/sha256"
  "os"
  "encoding/hex"
  "path"
)


// the input channel for the script writer
var calcHashChan chan *FileWork=make(chan *FileWork,chanLength)

// due to defer, it is a separate function
func calcFromFile(fileName string)(hash string, err error){
  // file 
  f, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer f.Close()
  // hash calculation
	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return
	}
  sum:=h.Sum(nil) // has sum as 32 byte array
  // split up to get folders with up to 64k entries
  var bb [4]string
  // folders
  bb[0]=hex.EncodeToString(sum[0:2]) 
  bb[1]=hex.EncodeToString(sum[2:4]) 
  bb[2]=hex.EncodeToString(sum[4:6]) 
  // remaining is file name
  bb[3]=hex.EncodeToString(sum[6:32]) 
  hash=path.Join(bb[0],bb[1],bb[2],bb[3])
  return
}

// worker function
func calcHash() {
  // running
  running.Add(1)
  defer running.Done()
  // channel
  chancloser.Claim(checkTargetChan)
  defer chancloser.Release(checkTargetChan)
  chancloser.Claim(errorWorkChan)
  defer chancloser.Release(errorWorkChan)
  // log.Println(" working")
  // setup done
  worked:=0
  for entry:= range calcHashChan{
    worked++
    <-workTickets
    h,err:=calcFromFile(entry.Path)
    if err != nil {
      // sending to error channel
      errorWorkChan <- &Err{entry.Path,"hash failed "+err.Error(),E_ERROR}
      continue
    }
    // entry.record("hash calc: "+h)
    entry.Hash=h
    checkTargetChan <- entry
  }
  log.Printf("done: %d\n",worked)
}
