#!/usr/bin/env bash 
DB_NAME="zunkasrv.db" 
if [ -z "$ZUNKAPATH" ]
then
	printf "error: ZUNKAPATH enviorment not defined.\n" >&2
	exit 1 
else
	printf "Creating db %s/%s\n" $ZUNKAPATH $DB_NAME
fi
sqlite3 $ZUNKAPATH/db/$DB_NAME < tables.sql