#!/bin/bash

#Authors
# Bogdan Gabriel Hortea

# Print the usage message
function printHelp() {
  echo "Usage: "
  echo "  deploy_use_case_docker_compose.sh [-l <deploy_logic.sh>]"
  echo "    -l <deploy_logic_N.sh> - "   
  echo " Ej: ./deploy_use_case_docker_compose.sh -l deploy_logic.sh"
}

if [ $# -eq 0 ]
  then
  printHelp
    exit 0
fi

while getopts "h?l:anv" opt; do
  case "$opt" in
  h | \?)
    printHelp
    exit 0
    ;;
  l)
    SCRIPT_LOGIC=$OPTARG
    ;;
esac
done

# STEP 1: CLEAN NETWORK --------------------------------------------------------------------------------------------------------
./network.sh down

# STEP 2: DEPLOY NETWORK --------------------------------------------------------------------------------------------------------
./network.sh up

# STEP 3: DEPLOY CHAINCODE --------------------------------------------------------------------------------------------------------
./network.sh createChannel
./network.sh deployCC -ccn basic -ccp ./chaincode-go -ccl go

# Exporting paths
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/

# Environment variables for Org1
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

# STEP 4: EXECUTE LOGIC SCRIPT ----------------------------------------------------------------------------------------------------
echo "STARTING LOGIC SCRPIT ..."
eval "./$SCRIPT_LOGIC" 

exit 0
