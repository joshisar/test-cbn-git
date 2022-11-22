#!/bin/bash

# imports
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/utils.sh"

source $DIR/../network/organization/admin/admin.env
source $DIR/../.env

validateEnvVariables

installChaincode() {
  infoln "Installing chaincode..."
  set -x
  peer lifecycle chaincode install ${CC_SRC_PATH}/${CC_NAME}.tar.gz >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
  verifyResult $res "Chaincode installation has failed"
  successln "Chaincode is installed"
}

# Install chaincode
installChaincode

# Query whether the chaincode is installed
queryInstalled

exit 0
