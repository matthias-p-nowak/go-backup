# go-backup

## Caveat

Backup via tape decouples the data from the reader/writer unit - whcih should be easily replaceable. Furthermore, tapes can be stored in different locations. That demands an encryption in order to secure the data. By physically removing the tapes from the library/loader/writer, one is secured against unauthorized tampering. However, using tapes requires a backup regime like full/incremental.

Backup via disks allows random access to data files. Disk space is expensive, hence deduplication is advised. The disks are physically connected to the drive, hence are receptive to local disasters. Among those disaster scenarios is the ransom virus. Hence, the backed up station shouldn't have write access to already written files. Moreover, all backup files should be mirrored to another location.

# Overview

**Go-backup** stores the backup on a disk as a medium. 
Content and overview are stored separately. 
The overview is a shell script that contains both the functions and function calls. 

The idea is to store the content of a backuped file under a filename that is made from the checksum of the content.
This way, an automatic deduplication is done.
The information, what content belongs to what file, or what to create, is stored in a shell script.

Storing the content separate from the overview enables full backups each time. 
However, this creates a minor challenge to delete the content files that are too old.
The solution to that issue is the KeepFree program from [http://github.com/matthias-p-nowak/keepfree].

## Backup

## Cache of file information

Instead of calculating the hash of the content each time, a small key-value database is used. 
Here, bolt is used with 2 buckets. 
The older bucket gets removed when closed, all has information is stored in the new one. 

## Shell script

Instead of a tar file, a certain backup consists of an overview file in form of a shell script and the folder of all content files.
The shell script will decompress the content files and store the content into the right places.
Afterwards, the file mode and ownership is corrected.

### Setup / Prolog

Restore has to fetch the content from some place and restore the files into another place. 
Both should be configurable via Environment variables

### Restore function



### Files

### Epilog

The epilog just consists of an *exit* call, with the total number of files on the last line
