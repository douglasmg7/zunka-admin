#!/usr/bin/env bash

# ZUNKAPATH must be defined.
[[ -z "$ZUNKAPATH" ]] && printf "error: ZUNKAPATH enviorment not defined.\n" >&2 && exit 1 
[[ -z "$GS" ]] && printf "error: GS enviorment not defined.\n" >&2 && exit 1 

# Create dirs if not exist.
mkdir -p $ZUNKAPATH/log
mkdir -p $ZUNKAPATH/xml

# Create dbs.
$GS/zunkasrv/bin/db/create_db.sh
$GS/aldowsc/bin/db/create_db.sh
$GS/allnations/bin/db/create_db.sh
