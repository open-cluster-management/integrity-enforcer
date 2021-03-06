#!/bin/bash
CMDNAME=`basename $0`
if [ $# -ne 2 ]; then
  echo "Usage: $CMDNAME <YAML files directory> <pubring-key>" 1>&2
  exit 1
fi
if [ ! -e $3 ]; then
  echo "$3 does not exist"
  exit 1
fi
SCRIPT_DIR=$(cd $(dirname $0); pwd)
TARGET_DIR=$1
PUBRING_KEY=$2

find ${TARGET_DIR} -type f -name "*.yaml" | while read file;
do
  echo "Verifying signature annotation in ${file}"

  $SCRIPT_DIR/../gpg-annotation-verify.sh  "$file" ${PUBRING_KEY}

done
