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

func calcFromFile(fileName string)(hash string, err error){
  f, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
  sum:=h.Sum(nil)
  var bb [4]string
  bb[0]=hex.EncodeToString(sum[0:2]) 
  bb[1]=hex.EncodeToString(sum[2:4]) 
  bb[2]=hex.EncodeToString(sum[4:6]) 
  bb[3]=hex.EncodeToString(sum[6:32]) 
  hash=path.Join(bb[0],bb[1],bb[2],bb[3])
  return
}

func calcHash() {
  running.Add(1)
  defer running.Done()
  chancloser.Claim(checkTargetChan)
  defer chancloser.Release(checkTargetChan)
  chancloser.Claim(debugSinkChan)
  defer chancloser.Release(debugSinkChan)
  defer log.Println("calcHash: done")
  log.Println("calcHash: working")
  //
  for entry:= range calcHashChan{
    // TODO: calculate hash sum
    h,err:=calcFromFile(entry.Path)
    if err != nil {
      entry.record("hash failed: "+err.Error())
      debugSinkChan <- entry
      continue
    }
    entry.record("hash calc: "+h)
    entry.Hash=h
    checkTargetChan <- entry
  }
}
