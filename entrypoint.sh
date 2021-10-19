#!/bin/bash
set -e

### Usage
### ----------------------------------------
### Set MARKETPLACE_CLI_VERSION env var or pass arg
### Print ls - export VERBOSE=true
### ./entrypoint.sh "$MARKETPLACE_CLI_VERSION"
### ----------------------------------------



_ROOT_DIR="${PWD}"
_WORKDIR="${_ROOT_DIR}/vmware-mkpcli"
#_MARKETPLACE_CLI_VERSION=${"0.2.4":-$MARKETPLACE_CLI_VERSION} # Use env or arg

_MARKETPLACE_CLI_VERSION="v0.4.9"

_MKP_CLI_VERSION=${_MARKETPLACE_CLI_VERSION}
_DOWNLOAD_URL=""
_DOWNLOAD_FILENAME="mkpcli-linux"



msg_log(){
    msg=$1
    echo -e ">> [LOG]: ${msg}"
}


set_workdir(){
    env
    echo "--------------"
    echo $MARKETPLACE_CLI_VERSION
    echo "${MARKETPLACE_CLI_VERSION}"
    echo "--------------"   
    mkdir -p "${_WORKDIR}"
    cd "${_WORKDIR}"
}

set_download_url(){
    msg_log "Setting _DOWNLOAD_URL"

    msg_log "${MKP_CLI_VERSION}"
    
     _DOWNLOAD_URL="https://github.com/vmware-labs/marketplace-cli/releases/download/${_MKP_CLI_VERSION}/mkpcli-linux"
     
    msg_log "_DOWNLOAD_URL = ${_DOWNLOAD_URL}"
}    

download_mkp_cli(){
    msg_log "Downloading ..."
    wget "$_DOWNLOAD_URL" &&
    [[ $_VERBOSE = "true" ]] && ls -lah "$_DOWNLOAD_FILENAME"
    chmod +x mkpcli-linux
}

install_mkp_cli(){
    msg_log "Installing ..."
    mv mkpcli-linux /usr/bin/mkpcli
  
}

test_mkp_cli(){
    msg_log "Printing MKP CLI installed version"
    mkpcli version
}

#Main
set_workdir
set_download_url
download_mkp_cli
install_mkp_cli
test_mkp_cli