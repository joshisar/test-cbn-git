#!/bin/bash

# imports
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/utils.sh"

source $DIR/../network/organization/admin/admin.env
source $DIR/../.env

validateEnvVariables

approveForMyOrg() {
  set -x
  peer lifecycle chaincode approveformyorg -o $ORDERER_URL \
    --ordererTLSHostnameOverride $ORDERER_TLS_OVERRIDE_NAME \
    --tls --cafile "$ORDERER_TLS_CA_FILE" \
    --channelID $CHANNEL_NAME \
    --name ${CC_NAME} \
    --version ${CC_VERSION} \
    --package-id ${PACKAGE_ID} \
    --sequence ${CC_SEQUENCE} \
    ${INIT_REQUIRED} >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
  verifyResult $res "Chaincode definition approved on channel '$CHANNEL_NAME' failed"
  successln "Chaincode definition approved on channel '$CHANNEL_NAME'"
}

# Query whether the chaincode is installed
queryInstalled

## Approve the definition
approveForMyOrg

exit 0
