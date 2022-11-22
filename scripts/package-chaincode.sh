#!/bin/bash

# imports
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/utils.sh"

source $DIR/../network/organization/admin/admin.env
source $DIR/../.env

validateEnvVariables

packageChaincode() {
  infoln "Packaging chaincode..."
  set -x
  peer lifecycle chaincode package ${CC_SRC_PATH}/${CC_NAME}.tar.gz \
    --path ${CC_SRC_PATH} \
    --lang ${CC_RUNTIME_LANGUAGE} \
    --label ${CC_NAME}_${CC_VERSION} >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
  verifyResult $res "Chaincode packaging has failed"
  successln "Chaincode is packaged"
}

compileSourceCode() {
  # do some language specific preparation to the chaincode before packaging
  if [ "$CC_RUNTIME_LANGUAGE" = "golang" ]; then

    infoln "Vendoring Go dependencies at $CC_SRC_PATH"
    pushd $CC_SRC_PATH
    go mod vendor
    popd
    successln "Finished vendoring Go dependencies"
  fi
}

# Compile the source code
compileSourceCode

# Package the chaincode
packageChaincode

exit 0
