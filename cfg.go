package main

/*
 * reading and preparing the configuration structure 
 */
 
import (
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "log"
)

/*
 * The structure that is read from Yaml configuration file
 */
type CFG struct {
  Destination string `yaml:"destination"`
  Cache string `yaml:"cache"`
  MinAge string `yaml:"minage"`
  Include []string `yaml:"include"`
  Exclude []string `yaml:"exclude"`
}


/*
 * reading configuration from a yaml file
 */
func GetCfg(filename string) (cfg *CFG) {
  cfg = new(CFG)
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
