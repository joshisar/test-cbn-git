#!/bin/bash

# imports
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/utils.sh"

source $DIR/../network/organization/admin/admin.env
source $DIR/../.env

validateEnvVariables

eval "$(shdotenv -e $DIR/../.secrets/env)"

commitChaincodeDefinition() {
  set -x
  peer lifecycle chaincode commit -o $ORDERER_URL \
    --ordererTLSHostnameOverride $ORDERER_TLS_OVERRIDE_NAME \
    --tls --cafile "$ORDERER_TLS_CA_FILE" \
    --channelID $CHANNEL_NAME \
    --name $CC_NAME \
    $PEER_NODES_ADDRESSES_AND_TLS_CONFIG \
    --version $CC_VERSION \
    --sequence $CC_SEQUENCE \
    $INIT_REQUIRED >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
  verifyResult $res "Chaincode definition commit failed on channel '$CHANNEL_NAME' failed"
  successln "Chaincode definition committed on channel '$CHANNEL_NAME'"
}

queryCommitted() {
  EXPECTED_RESULT="Version: ${CC_VERSION}, Sequence: ${CC_SEQUENCE}, Endorsement Plugin: escc, Validation Plugin: vscc"
  infoln "Querying chaincode definition on channel '$CHANNEL_NAME'..."
  local rc=1
  local COUNTER=1
  # continue to poll
  # we either get a successful response, or reach MAX RETRY
  while [ $rc -ne 0 -a $COUNTER -lt $MAX_RETRY ]; do
    sleep $DELAY
    infoln "Attempting to Query committed status, Retry after $DELAY seconds."
    set -x
    peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name ${CC_NAME} >&log.txt
    res=$?
    { set +x; } 2>/dev/null
    test $res -eq 0 && VALUE=$(cat log.txt | grep -o '^Version: '$CC_VERSION', Sequence: [0-9]*, Endorsement Plugin: escc, Validation Plugin: vscc')
    test "$VALUE" = "$EXPECTED_RESULT" && let rc=0
    COUNTER=$(expr $COUNTER + 1)
  done
  cat log.txt
  if test $rc -eq 0; then
    successln "Query chaincode definition successful on channel '$CHANNEL_NAME'"
  else
    fatalln "After $MAX_RETRY attempts, Query chaincode definition result is INVALID!"
  fi
}

## Now that we know for sure both orgs have approved, commit the definition
commitChaincodeDefinition

## Query on peer to see that the definition committed successfully
queryCommitted

exit 0
