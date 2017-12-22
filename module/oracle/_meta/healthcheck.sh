#!/bin/bash
export ORACLE_HOME=/u01/app/oracle/product/11.2.0/xe

CHECK=$(/u01/app/oracle/product/11.2.0/xe/bin/sqlplus -s system/oracle@xe <<END
  set pagesize 0 feedback off verify off heading off echo off;
  select decode(open_mode, 'READ WRITE', 0, 1) open from v\$database;
  exit;
END
)

ERR=1
MAX_TRIES=5
COUNT=0

while [ ${COUNT} -lt ${MAX_TRIES} ]; do
   if [ "${CHECK}" = "${CHECK}" -a  ${CHECK} -eq 0 ] 2> /dev/null; then
     exit 0
   fi
   let COUNT=COUNT+1
   sleep 2
done

exit $ERR


