#!/bin/bash

C_RESET='\033[0m'
C_RED='\033[0;31m'
C_GREEN='\033[0;32m'
C_BLUE='\033[0;34m'
C_YELLOW='\033[1;33m'

# println echos string
function println() {
  echo -e "$1"
}

# errorln echos in red color
function errorln() {
  println "${C_RED}${1}${C_RESET}"
}

# successln echos in green color
function successln() {
  println "${C_GREEN}${1}${C_RESET}"
}

# infoln echos in blue color
function infoln() {
  println "${C_BLUE}${1}${C_RESET}"
}

# warnln echos in yellow color
function warnln() {
  println "${C_YELLOW}${1}${C_RESET}"
}

# fatalln echos in red color and exits with fail status
function fatalln() {
  errorln "$1"
  exit 1
}

verifyResult() {
  if [ $1 -ne 0 ]; then
    fatalln "$2"
  fi
}

validateEnvVariables() {
  println "executing from ${DIR} with the following"
  println "- FABRIC_CFG_PATH: ${C_GREEN}${FABRIC_CFG_PATH}${C_RESET}"
  println "- CHANNEL_NAME: ${C_GREEN}${CHANNEL_NAME}${C_RESET}"
  println "- CC_NAME: ${C_GREEN}${CC_NAME}${C_RESET}"
  println "- CC_SRC_PATH: ${C_GREEN}${CC_SRC_PATH}${C_RESET}"
  println "- CC_SRC_LANGUAGE: ${C_GREEN}${CC_SRC_LANGUAGE}${C_RESET}"
  println "- CC_RUNTIME_LANGUAGE: ${C_GREEN}${CC_RUNTIME_LANGUAGE}${C_RESET}"
  println "- CC_VERSION: ${C_GREEN}${CC_VERSION}${C_RESET}"
  println "- CC_SEQUENCE: ${C_GREEN}${CC_SEQUENCE}${C_RESET}"
  println "- CC_INIT_FCN: ${C_GREEN}${CC_INIT_FCN}${C_RESET}"
  println "- DELAY: ${C_GREEN}${DELAY}${C_RESET}"
  println "- MAX_RETRY: ${C_GREEN}${MAX_RETRY}${C_RESET}"
  println "- VERBOSE: ${C_GREEN}${VERBOSE}${C_RESET}"

  # User has not provided a name
  if [ -z "$CC_NAME" ] || [ "$CC_NAME" = "NA" ]; then
    fatalln "No chaincode name was provided. Please provide it in the environment variable CC_NAME"

  # User has not provided a path
  elif [ -z "$CC_SRC_PATH" ] || [ "$CC_SRC_PATH" = "NA" ]; then
    fatalln "No chaincode path was provided. Please provide it in the environment variable CC_SRC_PATH"
  fi

  if [ "$CC_RUNTIME_LANGUAGE" != "golang" ]; then

    fatalln "The chaincode language ${CC_RUNTIME_LANGUAGE} is not supported by this script. Supported chaincode language is go"
    exit 1
  fi

  INIT_REQUIRED="--init-required"
  # check if the init fcn should be called
  if [ "$CC_INIT_FCN" = "NA" ]; then
    INIT_REQUIRED=""
  fi
}

queryInstalled() {
  set -x
  peer lifecycle chaincode queryinstalled >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
  PACKAGE_ID=$(sed -n "/${CC_NAME}_${CC_VERSION}/{s/^Package ID: //; s/, Label:.*$//; p;}" log.txt)
  verifyResult $res "Query installed has failed"
  successln "Query installed successfully on channel"
}
