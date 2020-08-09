package main

/*
 * reading and preparing the configuration structure 
 */
 
import (
  "fmt"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "log"
  "time"
)

/*
 * The structure that is read from Yaml configuration file
 */
type CFG struct {
  // where the files and scripts will be stored
  Destination string `yaml:"destination"`
  // name of the bolt database file
  Cache string `yaml:"cache"`
  // only take backup of files younger than this
  MinAge string `yaml:"minage"`
  NumWorkers int  `yaml:"workers"`
  // walk all files that are on this device starting with this path
  Include []string `yaml:"include"`
  // list all matches that will be excluded
  Exclude []string `yaml:"exclude"`
}


/*
 * reading configuration from a yaml file
 */
func GetCfg(filename string) (cfg *CFG) {
  cfg = new(CFG)
  cfg.NumWorkers = 1
  // file is small enough
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    log.Fatal(err)
  }
  err = yaml.Unmarshal(data, cfg)
  if err != nil {
    log.Fatal(err)
  }
  return
}

/*
 * giving back the unix time stamp (seconds after 1970) for comparing with files
 * new files will not be backed up
 */
func (cfg *CFG)AgeCutoff()(cutOff int64){
  var i int64
  var s string
  fmt.Sscanf(cfg.MinAge,"%d%s",&i,&s)
  // turning unit into a multiplier
  for _,val := range s {
    switch val {
      case 's': // is in seconds
      case 'm': i*=60
      case 'h': i*=3600
      case 'd': i*=3600*24
      case 'w': i*=3600*24*7
    }
  }
  t:=time.Now().Unix()
  cutOff=t-i
  return
}
