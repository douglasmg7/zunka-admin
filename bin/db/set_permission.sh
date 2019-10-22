#!/usr/bin/env bash 

if [ -z $2 ]
then
  echo "Usage: $0 email permission"
  echo "Example: $0 fulango@gmail.com 1"
  exit 1
fi

DB_NAME="zunkasrv.db" 
if [ -z "$ZUNKAPATH" ]
then
	printf "error: ZUNKAPATH enviorment not defined.\n" >&2
	exit 1 
else
	printf "Updating %s/%s\n" $ZUNKAPATH $DB_NAME...
fi

# echo UPDATE user SET permission=$2 WHERE email="$1"

sqlite3 $ZUNKAPATH/db/$DB_NAME <<END_SQL
UPDATE user SET permission=$2 WHERE email="$1";
END_SQL