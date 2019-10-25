#!/usr/bin/env bash 

# ZUNKAPATH not defined.
if [ -z "$ZUNKAPATH" ]; then
	printf "error: ZUNKAPATH not defined.\n" >&2
	exit 1 
fi

# ZUNKA_SRV_DB not defined.
if [ -z "$ZUNKA_SRV_DB" ]; then
	printf "error: ZUNKA_SRV_DB not defined.\n" >&2
	exit 1 
fi

printf "Removing db %s/%s\n" $ZUNKAPATH/db/$ZUNKA_SRV_DB
rm $ZUNKAPATH/db/$ZUNKA_SRV_DB