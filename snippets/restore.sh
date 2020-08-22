
cat << EOD
Restoring files from $BACKUP to destination $DEST.
If this is correct, enter 'yes'. If this is not correct, set the environment variables similar to:
export BACKUP=$BACKUP
export DEST=$DEST
EOD
# reading answer from input
 read A

test  "$A" = "yes" || exit
echo "doing backup"
TOTAL=$(tail -1 $0)
echo $TOTAL
FINISHED=0

BUNZIP2=$(type -p bunzip2)
CHMOD=$(type -p chmod)
CHOWN=$(type -p chown)
MKDIR=$(type -p mkdir)
LN=$(type -p ln)
MKFIFO=$(type -p mkfifo)
RM=$(type -p rm)
TOUCH=$(type -p touch)

#if type busybox
#then
  #echo "using busybox"
  #BB=$(type -p busybox)
  #BUNZIP2="$BB bunzip2"
  #CHMOD="$BB chmod"
  #CHOWN="$BB chown"
  #MKDIR="$BB mkdir"
  #LN="$BB ln"
  #MKFIFO="$BB mkfifo"
  #RM="$BB rm"
  #TOUCH="$BB touch"
#fi


OD=$(date +%s)

finish(){
  FINISHED=$(( FINISHED + 1 ))
	DN=$(date +%s)
	if [ $OD != $DN ]
	then
		OD=$DN
	  DONE=$(( 100 * FINISHED / TOTAL ))
 		echo -en "${FINISHED}/${TOTAL} = ${DONE}%\r"
	fi
}

f(){
  d=${DEST}/$4
  ${BUNZIP2} <${BACKUP}/f/$3 >$d
  ${CHMOD} $2 $d
  ${CHOWN} $1 $d
  finish
}


e(){
  d=${DEST}/$3
  ${TOUCH} $d
  ${CHMOD} $2 $d
  ${CHOWN} $1 $d
  finish
}

s(){
  d=${DEST}/$2
  [ -e "$d" ] && ${RM} $d
  ${LN} -s $3 $d
  ${CHOWN} -h $1 $d
  finish
}

d(){
  d=${DEST}/$3
  ${MKDIR} -p $d
  ${CHMOD} $2 $d
  ${CHOWN} $1 $d
  finish
}

p(){
  d=${DEST}/$3
  ${MKFIFO} $d
  ${CHMOD} $2 $d
  ${CHOWN} -h $1 $d
  finish 
}

