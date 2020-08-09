
package main

/*
 * It is a receiver channel goroutine for debugging
 */
 
import (
  "log"
)

/*
 a sink that prints out the received values
*/
var debugSinkChan chan(*FileWork)=make(chan *FileWork,chanLength)

func debugSink(){
  running.Add(1)
  defer running.Done()
  // setup done
  defer log.Println("debugSink: done")
  log.Println("debugSink: started")
  for fw:=range debugSinkChan {
    log.Printf("%#v\n",fw)
  }
}
