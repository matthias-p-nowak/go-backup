package main

import(
  "io/ioutil"
  "html/template"
  "log"
  "log/syslog"
  "net/smtp"
  "os"
  "time"
)

const (
  E_WARNING = iota
  E_ERROR
)

type Err struct{
  Path string
  Cause string
  Flag int
}

// the input channel for the script writer
var errorWorkChan=make(chan *Err,chanLength)

// List of errors during the program run, can be a lot...
var Errors []Err
var Warnings []Err
var Unknowns []Err


func errorWork(cfg *CFG) {
  running.Add(1)
  defer running.Done()
  // setup done
  worked:=0
  for entry:= range errorWorkChan{
    switch(entry.Flag){
      case E_WARNING:
        Warnings=append(Warnings,*entry)
      case E_ERROR:
        Errors=append(Errors,*entry)
      default:
        Unknowns=append(Unknowns,*entry)
    }
  }
  log.Printf("error entries: %d\n",worked)
}

type MailData struct {
  Hostname string
  Date string
  Errors []Err
  Warnings []Err
  Unknowns []Err
}

func errorMail(cfg *CFG){
  cntE:=len(Errors)
  cntW:=len(Warnings)
  cntU:=len(Unknowns)
  sl,err :=syslog.NewLogger(syslog.LOG_INFO | syslog.LOG_LOCAL0,0)
  if err != nil {log.Fatal(err)  }
  sl.Printf("got %d errors, %d warnings, %d unknowns\n",cntE,cntW,cntU)
  if cntE+cntW+cntU == 0 {
    return
  }
  sf:=GetStored("snippets/errormail.html")
  bb,err := ioutil.ReadAll(sf)
  if err != nil {log.Fatal(err)  }
  str:=string(bb)
  temp:=template.New("mail")
  temp,err=temp.Parse(str)
  if err != nil {
     log.Fatal(err)
  }
  md:=MailData{}
  md.Date=time.Now().Format("2006-01-02_15-04-05")
  md.Hostname,err=os.Hostname()
  if err != nil {
    log.Fatal(err)
  }
  md.Errors=Errors
  md.Warnings=Warnings
  md.Unknowns=Unknowns
  
  cl,err:=smtp.Dial(cfg.MailHost);  if err!=nil{log.Fatal(err)}
  err=cl.Mail(cfg.MailFrom);  if err!=nil{log.Fatal(err)}
  for _,rcpt:= range cfg.MailTo {
    err=cl.Rcpt(rcpt);  if err!=nil{log.Fatal(err)}
  }
  wr,err:=cl.Data();  if err!=nil{log.Fatal(err)}
  err=temp.Execute(wr,md);  if err!=nil{log.Fatal(err)}
  err=wr.Close();  if err!=nil{log.Fatal(err)}
  err=cl.Quit();  if err!=nil{log.Fatal(err)}

}
