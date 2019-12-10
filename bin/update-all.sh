#!/usr/bin/env bash

[[ -z "$GS" ]] && printf "error: GS - Go Source enviorment variable not defined.\n" >&2 && exit 1 
[[ -z "$ZUNKAPATH" ]] && printf "error: ZUNKAPATH enviorment variable not defined.\n" >&2 && exit 1 

pull_roll () {
    [[ -z $1 ]] && echo "error: pull_roll function called without path argument." && exit 1 

    cd $1

    printf "\n:: Synchronizing %s ...\n" $1
    REV_OLD=`git rev-parse HEAD`
    git pull
    REV_NEW=`git rev-parse HEAD`
    SOME_FILES_CHANGED=`git diff $REV_OLD --name-only`
    GOPKG_LOCK_CHANGED=`git diff $REV_OLD --name-only | grep "Gopkg\.lock$"`
    SECRET_FILES_CHANGED=`git diff $REV_OLD --name-only | grep "\.secret$"`
    # Confirm if local repository not have modifications.
    if [[ $REV_NEW == $REV_OLD && ! -z $SOME_FILES_CHANGED ]];then
        printf "Local repository have modifications.\n"
        return
    fi
    if [[ ! -z $GOPKG_LOCK_CHANGED ]]; then
        echo :: Running dep ensure -vendor-only ...
        dep ensure -vendor-only
    fi
    if [[ ! -z $SECRET_FILES_CHANGED ]]; then
        echo :: Revealing secret files...
        git secret reveal
    fi
    if [[ ! -z $SOME_FILES_CHANGED ]]; then 
        echo " rev-new:" `git rev-parse HEAD`
        echo " rev-old:" $REV_OLD
        return 1
    fi
}

# # bluetang 
# pull_roll $GS/bluetang
# if [[ $? == 1 ]]; then
    # INSTALL_ZUNKASRV=true
# fi

# # currency
# pull_roll $GS/currency
# if [[ $? == 1 ]]; then
    # INSTALL_ZUNKASRV=true
    # INSTALL_ALDOWSC=true
# fi

# # aldoutil
# pull_roll $GS/aldoutil
# if [[ $? == 1 ]]; then
    # INSTALL_ZUNKASRV=true
    # INSTALL_ALDOWSC=true
# fi

# aldowsc
pull_roll $GS/aldowsc
if [[ $? == 1 ]]; then
    INSTALL_ALDOWSC=true
fi
# aldowsc.service
if [[ ! -z `git diff --name-only $REV_OLD | grep "install-aldowsc-service\.sh$"` ]]; then
    printf "\n:: Installing alwdowsc.service...\n"
    ./bin/install-aldowsc-service.sh
fi

# zoomwsc
pull_roll $GS/zoomwsc
if [[ $? == 1 ]]; then
    INSTALL_ZOOMWSC=true
fi

# zunkasrv
pull_roll $GS/zunkasrv
if [[ $? == 1 ]]; then
    INSTALL_ZUNKASRV=true
fi

########################################################
# Install
########################################################
# Install aldowsc.
if [[ $INSTALL_ALDOWSC == true ]]; then
    printf "\n:: Installing aldowsc...\n"
    cd $GS/aldowsc
    go install
fi

# Install zoomwsc.
if [[ $INSTALL_ZOOMWSC == true ]]; then
    printf "\n:: Installing zoomwsc...\n"
    cd $GS/zoomwsc
    go install
fi

# Install zunkasrv.
if [[ $INSTALL_ZUNKASRV == true ]]; then
    printf "\n:: Installing zunkasrv...\n"
    cd $GS/zunkasrv
    go install
    printf "\n:: Signaling to restart zunka_srv...\n"
    echo true > $ZUNKAPATH/restart-zunka-srv 
fi
