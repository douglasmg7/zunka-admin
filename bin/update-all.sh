#!/usr/bin/env bash

[[ -z "$GS" ]] && printf "error: GS - Go Source enviorment variable not defined.\n" >&2 && exit 1 
[[ -z "$ZUNKAPATH" ]] && printf "error: ZUNKAPATH enviorment variable not defined.\n" >&2 && exit 1 

echo :: Checking bluetang repository...
cd $GS/bluetang
RES=`git pull`
if [[ ! $RES == Already* ]]; then 
    INSTALL_ZUNKASRV=true
    echo bluetang repository updated.
fi

echo :: Checking currency repository...
cd $GS/currency
RES=`git pull`
if [[ ! $RES == Already* ]]; then
    INSTALL_ZUNKASRV=true
    INSTALL_ALDOWSC=true
    echo currency repository updated.
fi

echo :: Checking aldoutil repository...
cd $GS/aldoutil
RES=`git pull`
if [[ ! $RES == Already* ]]; then
    INSTALL_ZUNKASRV=true
    INSTALL_ALDOWSC=true
    echo aldoutil repository updated.
fi

echo :: Checking aldowsc repository...
cd $GS/aldowsc
RES=`git pull`
if [[ ! $RES == Already* ]]; then
    INSTALL_ALDOWSC=true
    echo aldowsc repository updated.
fi

echo :: Checking zunkasrv repository...
cd $GS/zunkasrv
RES=`git pull`
if [[ ! $RES == Already* ]]; then
    INSTALL_ZUNKASRV=true
    echo zunkasrv repository updated.
fi

# Install aldowsc.
if [[ $INSTALL_ALDOWSC == true ]]; then
    echo :: Installing aldowsc...
    cd $GS/aldowsc
    go install
fi

# Install zunkasrv.
if [[ $INSTALL_ZUNKASRV == true ]]; then
    echo :: Installing zunkasrv...
    cd $GS/zunkasrv
    go install
    echo :: Setting zunka srv to be restarted...
    echo true > $ZUNKAPATH/restart-zunka-srv 
fi
