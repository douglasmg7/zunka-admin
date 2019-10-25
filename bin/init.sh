#!/usr/bin/env bash

# ZUNKAPATH must be defined.
[[ -z "$ZUNKAPATH" ]] && printf "error: ZUNKAPATH enviorment not defined.\n" >&2 && exit 1 
[[ -z "$GS" ]] && printf "error: GS enviorment not defined.\n" >&2 && exit 1 

# Create dirs if not exist.
mkdir -p $ZUNKAPATH/log
mkdir -p $ZUNKAPATH/xml

# Create list with selected products if list not exist.
if [[ ! -f $ZUNKAPATH/list/categSel.list ]]; then
    echo "Creating list of selected products."
    mkdir -p $ZUNKAPATH/list
    cat > $ZUNKAPATH/list/categSel.list << EOF
monitor
EOF
fi

# Create dbs.
$GS/zunkasrv/bin/db/create_db.sh
$GS/aldowsc/bin/db/create_db.sh
