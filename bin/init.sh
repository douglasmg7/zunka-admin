#!/usr/bin/env bash

# ZUNKAPATH must be defined.
[[ -z "$ZUNKAPATH" ]] && printf "error: ZUNKAPATH enviorment not defined.\n" >&2 && exit 1 
[[ -z "$GS" ]] && printf "error: GS enviorment not defined.\n" >&2 && exit 1 

# Create dirs if not exist.
mkdir -p $ZUNKAPATH/log
mkdir -p $ZUNKAPATH/xml

# Create configuration file if not exist.
if [[ ! -f $ZUNKAPATH/config.toml ]]; then
    mkdir -p $ZUNKAPATH
    printf "Creating configuration file.\n"
    cat > $ZUNKAPATH/config.toml << EOF
# Configuration file.

[all]
env = "production"
logDir = "log"
dbDir = "db"
listDir = "list"
xmlDir = "xml"

[zunkasrv]
logFileName = "zunkasrv.log"
dbFileName = "zunkasrv.db"
port = "8080"

[aldowsc]
logFileName = "aldowsc.log"
dbFileName = "aldowsc.db"
minPrice = "2.000,00"
maxPrice = "100.000,00"
EOF
fi

# Create list with selected products if list not exist.
if [[ ! -d $ZUNKAPATH/list ]]; then
    echo "Creating list of selected products."
    mkdir -p $ZUNKAPATH/list
    cat > $ZUNKAPATH/list/categSel.list << EOF
monitor
EOF
fi

# Create dbs.
$GS/zunkasrv/bin/db/create_db.sh
$GS/aldowsc/bin/db/create_db.sh
