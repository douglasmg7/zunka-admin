#!/usr/bin/env bash

[[ -z "$GS" ]] && printf "error: GS - Go Source enviorment variable not defined.\n" >&2 && exit 1 
[[ -z "$ZUNKAPATH" ]] && printf "error: ZUNKAPATH enviorment variable not defined.\n" >&2 && exit 1 

########################################################
# bluetang 
########################################################
echo :: Fetching bluetang repository...
cd $GS/bluetang
git fetch
SOME_FILES_CHANGED=`git diff --name-only master...origin/master`
SECRET_FILES_CHANGED=`git diff --name-only master...origin/master | grep "*.secret"`
if [[ ! -z $SOME_FILES_CHANGED ]]; then 
    INSTALL_ZUNKASRV=true
    git merge
    echo Merging bluetang repository...
fi
if [[ ! -z $SECRET_FILES_CHANGED ]]; then
    echo :: Revealing secret files...
    git secret reveal
fi

########################################################
# currency
########################################################
echo :: Fetching currency repository...
cd $GS/currency
git fetch
SOME_FILES_CHANGED=`git diff --name-only master...origin/master`
SECRET_FILES_CHANGED=`git diff --name-only master...origin/master | grep "*.secret"`
if [[ ! -z $SOME_FILES_CHANGED ]]; then
    INSTALL_ZUNKASRV=true
    INSTALL_ALDOWSC=true
    git merge
    echo Merging currency repository...
fi
if [[ ! -z $SECRET_FILES_CHANGED ]]; then
    echo :: Revealing secret files...
    git secret reveal
fi

########################################################
# aldoutil
########################################################
echo :: Fetching aldoutil repository...
cd $GS/aldoutil
git fetch
SOME_FILES_CHANGED=`git diff --name-only master...origin/master`
SECRET_FILES_CHANGED=`git diff --name-only master...origin/master | grep "*.secret"`
if [[ ! -z $SOME_FILES_CHANGED ]]; then
    INSTALL_ZUNKASRV=true
    INSTALL_ALDOWSC=true
    git merge
    echo Merging aldoutil repository...
fi
if [[ ! -z $SECRET_FILES_CHANGED ]]; then
    echo :: Revealing secret files...
    git secret reveal
fi

########################################################
# aldowsc
########################################################
echo :: Fetching aldowsc repository...
cd $GS/aldowsc
git fetch
SOME_FILES_CHANGED=`git diff --name-only master...origin/master`
SECRET_FILES_CHANGED=`git diff --name-only master...origin/master | grep "*.secret"`
if [[ ! -z $SOME_FILES_CHANGED ]]; then
    INSTALL_ALDOWSC=true
    git merge
    echo Merging aldowsc repository...
fi
if [[ ! -z $SECRET_FILES_CHANGED ]]; then
    echo :: Revealing secret files...
    git secret reveal
fi

########################################################
# zoomwsc
########################################################
echo :: Fetching zoomwsc repository...
cd $GS/zoomwsc
git fetch
SOME_FILES_CHANGED=`git diff --name-only master...origin/master`
SECRET_FILES_CHANGED=`git diff --name-only master...origin/master | grep "*.secret"`
if [[ ! -z $SOME_FILES_CHANGED ]]; then
    INSTALL_ZOOMWSC=true
    git merge
    echo Merging zoomwsc repository...
fi
if [[ ! -z $SECRET_FILES_CHANGED ]]; then
    echo :: Revealing secret files...
    git secret reveal
fi

########################################################
# zunkasrv
########################################################
echo :: Fetching zunkasrv repository...
cd $GS/zunkasrv
git fetch
SOME_FILES_CHANGED=`git diff --name-only master...origin/master`
SECRET_FILES_CHANGED=`git diff --name-only master...origin/master | grep "*.secret"`
if [[ ! -z $SOME_FILES_CHANGED ]]; then
    INSTALL_ZUNKASRV=true
    git merge
    echo Merging zunkasrv repository...
fi
if [[ ! -z $SECRET_FILES_CHANGED ]]; then
    echo :: Revealing secret files...
    git secret reveal
fi

########################################################
# Install
########################################################
# Install aldowsc.
if [[ $INSTALL_ALDOWSC == true ]]; then
    echo :: Installing aldowsc...
    cd $GS/aldowsc
    go install
fi

# Install zoomwsc.
if [[ $INSTALL_ZOOMWSC == true ]]; then
    echo :: Installing zoomwsc...
    cd $GS/zoomwsc
    go install
fi

# Install zunkasrv.
if [[ $INSTALL_ZUNKASRV == true ]]; then
    echo :: Installing zunkasrv...
    cd $GS/zunkasrv
    go install
    echo :: Signaling to restart zunka_srv...
    echo true > $ZUNKAPATH/restart-zunka-srv 
fi
