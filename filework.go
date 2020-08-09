
package main

import (
  // "log"
  "os"
  // "sync"
)

// Info about backup file that gets send around in the program
type FileWork struct {
  // file path related to this file
  Path string
  // files are deduplicted and stored under filenames made from a hash sum
  // hash sum string is divided so to represent a file tree
  Hash string
  // Modification time
  MTime int64
  // File size
  Size int64
  // user id
  Uid int
  // group id
  Gid int
  // file mode
  Mode string
  // info from the filewalk
  FileInfo os.FileInfo
  workDone []string
}

func (fw *FileWork) record(str string){
  fw.workDone=append(fw.workDone,str)
}
