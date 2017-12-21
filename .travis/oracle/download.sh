#!/bin/sh -e
# vim: set et sw=2 ts=2:

ORACLE_DOWNLOAD_DIR="oracle11g/xe"
ORACLE_DOWNLOAD_FILE="oracle-xe-11.2.0-1.0.x86_64.rpm.zip"

if [ -n "$ORACLE_DOWNLOAD_DIR" ]; then
  mkdir -p "$ORACLE_DOWNLOAD_DIR"
  ORACLE_DOWNLOAD_FILE="$(readlink -f "$ORACLE_DOWNLOAD_DIR")/$ORACLE_DOWNLOAD_FILE"
fi

if [ "${*#*--unless-exists}" != "$*" ] && [ -f "$ORACLE_DOWNLOAD_FILE" ]; then
  exit 0
fi

cd "$(dirname "$(readlink -f "$0")")"

echo "PhantomJS version $(phantomjs --version)"
npm install bluebird node-phantom-simple

export ORACLE_DOWNLOAD_FILE
export COOKIES='cookies.txt'
export USER_AGENT='Mozilla/5.0'

echo > "$COOKIES"
chmod 600 "$COOKIES"

exec node download.js

mv $ORACLE_DOWNLOAD_FILE dockerfiles/11.2.0.2
