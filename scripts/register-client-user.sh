#!/bin/bash

# imports
DIR="${BASH_SOURCE%/*}"
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
. "$DIR/utils.sh"

USERNAME="$1"
PASSWORD="$2"

NETWORK_DIR=/home/coder/${slug}/network

if [ ! -f "${NETWORK_DIR}/ca/admin/msp/keystore/key.pem" ]; then
  infoln "Enrolling CA Admin..."
  # CA Admin user folder does not exist, enroll CA Admin
  set -x

  fabric-ca-client enroll -d -u https://${CA_NAME}:${CA_PASSWORD}@${CA_NAME}.settlemint.com \
    --csr.hosts ${PEER_K8S_SERVICE_ADDRESS} \
    --csr.hosts ${PEER_ADDRESS} \
    --csr.hosts ${PEER_PUBLIC_ADDRESS} \
    --tls.certfiles "${NETWORK_DIR}/ca/tls-ca-cert.pem" \
    --mspdir "${NETWORK_DIR}/ca/admin/msp"
  RES=$?
  { set +x; } 2>/dev/null

  if [ $RES -eq 0 ]; then
    sleep 1

    infoln "Renaming CA Admin private key..."
    # Rename CA Admin private key file
    mv $(find ${NETWORK_DIR}/ca/admin/msp/keystore/ -maxdepth 1 -name "*_sk") \
      ${NETWORK_DIR}/ca/admin/msp/keystore/key.pem

    infoln "Renaming CA Admin CA certificate..."
    # Rename CA Admin CA certificate
    mv $(find ${NETWORK_DIR}/ca/admin/msp/cacerts/ -maxdepth 1 -name "*-settlemint-com.pem") \
      ${NETWORK_DIR}/ca/admin/msp/cacerts/ca-cert.pem
  fi
fi

infoln "Registering User ${USERNAME} as client..."
set -x
fabric-ca-client register --caname ${CA_NAME} -d \
  --id.name ${USERNAME} \
  --id.secret $PASSWORD \
  --id.type client \
  --csr.hosts ${PEER_K8S_SERVICE_ADDRESS} \
  --csr.hosts ${PEER_ADDRESS} \
  --csr.hosts ${PEER_PUBLIC_ADDRESS} \
  --tls.certfiles ${NETWORK_DIR}/ca/tls-ca-cert.pem \
  --mspdir ${NETWORK_DIR}/ca/admin/msp
RES=$?
{ set +x; } 2>/dev/null

sleep 1

infoln "Enrolling User ${USERNAME} as client..."
set -x
fabric-ca-client enroll -d -u https://${USERNAME}:${PASSWORD}@${CA_NAME}.settlemint.com \
  --csr.hosts ${PEER_K8S_SERVICE_ADDRESS} \
  --csr.hosts ${PEER_ADDRESS} \
  --csr.hosts ${PEER_PUBLIC_ADDRESS} \
  --tls.certfiles ${NETWORK_DIR}/ca/tls-ca-cert.pem \
  --mspdir ${NETWORK_DIR}/organization/${USERNAME}/msp
RES=$?
{ set +x; } 2>/dev/null

if [ $RES -eq 0 ]; then
  sleep 1

  infoln "Renaming ${USERNAME} private key..."
  mv $(find ${NETWORK_DIR}/organization/${USERNAME}/msp/keystore/ -maxdepth 1 -name "*_sk") \
    ${NETWORK_DIR}/organization/${USERNAME}/msp/keystore/key.pem

  infoln "Renaming ${USERNAME} CA certificate..."
  mv $(find ${NETWORK_DIR}/organization/${USERNAME}/msp/cacerts/ -maxdepth 1 -name "*-settlemint-com.pem") \
    ${NETWORK_DIR}/organization/${USERNAME}/msp/cacerts/ca-cert.pem

  cp ${NETWORK_DIR}/organization/peer/msp/config.yaml ${NETWORK_DIR}/organization/${USERNAME}/msp/config.yaml
  cp ${NETWORK_DIR}/organization/peer/peer.env ${NETWORK_DIR}/organization/${USERNAME}/${USERNAME}.env

  sed -i "s|organization/peer/msp$|organization/$USERNAME/msp|" ${NETWORK_DIR}/organization/${USERNAME}/${USERNAME}.env
fi
