#!/bin/bash

# imports
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/utils.sh"

source $DIR/../network/organization/admin/admin.env
source $DIR/../.env

validateEnvVariables

eval "$(shdotenv -e $DIR/../.secrets/env)"

chaincodeInvokeInit() {
  set -x
  fcn_call='{"function":"'${CC_INIT_FCN}'","Args":[]}'
  infoln "invoke fcn call:${fcn_call}"
  peer chaincode invoke -o $ORDERER_URL \
    --ordererTLSHostnameOverride $ORDERER_TLS_OVERRIDE_NAME \
    --tls --cafile "$ORDERER_TLS_CA_FILE" \
    --channelID $CHANNEL_NAME \
    --name $CC_NAME \
    $PEER_NODES_ADDRESSES_AND_TLS_CONFIG \
    --isInit -c ${fcn_call} >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
  verifyResult $res "Invoke execution failed "
  successln "Invoke transaction successful on channel '$CHANNEL_NAME'"
}

## Invoke the chaincode - this does require that the chaincode have the 'initLedger'
## method defined
if [ "$CC_INIT_FCN" = "NA" ]; then
  infoln "Chaincode initialization is not required"
else
  chaincodeInvokeInit
fi

exit 0
