#!/bin/bash

CMDNAME=`basename $0`
if [ $# -ne 2 ]; then
  echo "Usage: $CMDNAME  <input-file> <pubring-key>" 1>&2
  exit 1
fi

if [ ! -e $2 ]; then
  echo "$2 does not exist"
  exit 1
fi

INPUT_FILE=$1
PUBRING_KEY=$2

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    base='base64 -w 0'
    base_decode='base64 -d'
elif [[ "$OSTYPE" == "darwin"* ]]; then
    base='base64'
    base_decode='base64 -D'
fi


msg=$(yq r -d0 ${INPUT_FILE} 'metadata.annotations.message')
sign=$(yq r -d0 ${INPUT_FILE} 'metadata.annotations.signature')

cat $INPUT_FILE > /tmp/input

yq d /tmp/input metadata.annotations.message -i
yq d /tmp/input metadata.annotations.signature -i


msg_body=`cat /tmp/input | $base`

if [ "${msg}" != "${msg_body}" ]; then
   echo Input file content has been changed.
   exit 0
fi

if [ -z ${msg} ] || [ -z ${sign} ] ; then
   echo "Input file is not yet signed."
else
   echo $msg | ${base_decode} >  /tmp/msg
   echo $sign | ${base_decode} > /tmp/sign.sig

   status=$(gpg --no-default-keyring --keyring ${PUBRING_KEY} --dry-run --verify /tmp/sign.sig /tmp/msg 2>&1)
   result=$(echo $status | grep "Good" | wc -c)
   echo ----------------------------------------------
   if [ ${result} -gt 0 ]; then
      echo $status
      echo "Signature is successfully verified."
      exit 0
   else
      echo $status
      echo "Signature is invalid"
      exit 1
   fi
   echo --------------------------------------------------
   if [ -f /tmp/msg ]; then
     rm /tmp/msg
   fi

   if [ -f /tmp/sign.sig ]; then
     rm /tmp/sign.sig
   fi

   if [ -f /tmp/input ]; then
     rm  /tmp/input
   fi
fi


