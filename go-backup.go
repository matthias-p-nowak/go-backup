/*
Backup program that stores files in bzip2 form under filenames made from the hash sum of the content.
Configuration is read from go-backup.cfg. 
*/
package main

/*
 * TODO: implement a help 
 * TODO: implement command line flags
 * TODO: syslog
 * TODO: email when errors
 */

//go:generate go run scripts/go-bin.go -o snippets.go snippets

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net/smtp"
	"os"
	"runtime"
	"strings"
	"sync"
)

const (
	chanLength int =16
)

// the wait group all goroutines should take part
var running sync.WaitGroup

// dedicated function to set up the infrastructure
func setup(cfg *CFG,cache *Cache){
	workers:=cfg.NumWorkers
	if workers < 1 {
		workers=1
	}
	// worker that walk the files trees and send info out
	for i:=range cfg.Include {
		go discover(i,cfg)
	}
	// worker that retrieves from cache
	go fromCache(cache)
	for i:=0; i< workers;i++{
		go calcHash()
		go checkTarget(cfg)
		go bzip2Writer(cfg)
	}
	// TODO: calculate hash sum
	// TODO: check if target file exists
	// TODO: bzip2 workers - like 114.go
	// TODO: script writer
	go scriptWriter(cfg)
	go errorWork(cfg)
	go toCache(cache)
	// TODO: cache writer
	// temp channel to print out the structures
	// go debugSink()
}

func testMail(cfg *CFG){
	
  cl,err:=smtp.Dial(cfg.MailHost);  if err!=nil{log.Fatal(err)}
  err=cl.Mail(cfg.MailFrom);  if err!=nil{log.Fatal(err)}
  for _,rcpt:= range cfg.MailTo {
    err=cl.Rcpt(rcpt);  if err!=nil{log.Fatal(err)}
  }
  wr,err:=cl.Data();  if err!=nil{log.Fatal(err)}
 str:=`Content-type: text/html
Subject: Test mail for backup
 
Test mail
`
  str=strings.ReplaceAll(str,"\n","\r\n")
  wr2:=bufio.NewWriter(wr)
  wr2.WriteString(str)
  wr2.Flush()
  err=wr.Close();  if err!=nil{log.Fatal(err)}
  err=cl.Quit();  if err!=nil{log.Fatal(err)}

}

// Back up files
func main() {
	// setting logs
	log.SetFlags(log.LstdFlags|log.Lshortfile)
	//
	log.Print("go-backup started")
	defer log.Print("all done")
	cfgFN:=flag.String("c","go-backup.cfg","the config file to use")
	ex:=flag.Bool("e",false,"print an example config")
	tm:=flag.Bool("m",false,"sending a test mail")
	flag.Parse()
	if *ex {
		io.Copy(os.Stdout, GetStored("snippets/go-backup.cfg"))
		return
	}
	// configuration
	cfg := GetCfg(*cfgFN)
	if *tm {
		testMail(cfg)
		return
	}
	// open the bolt based cache
	cache := OpenCache(cfg.Cache)
	defer cache.Close()
	// do a setup
	setup(cfg,cache)
	// goroutines are created, but need to be running
	runtime.Gosched()
	// running waits until all goroutines are finished
	running.Wait()
	errorMail(cfg)
}
