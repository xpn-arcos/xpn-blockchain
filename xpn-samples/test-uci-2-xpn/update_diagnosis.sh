#!/bin/bash
#Authors
# Bogdan Gabriel Hortea
#

id="$1"
endTime=$(($(date +%s) + 120)) #120s, same value as in generate_data

while [ $(date +%s) -lt $endTime ]; do

    diagnosisCommand="peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile \"${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem\" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles \"${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt\" --peerAddresses localhost:9051 --tlsRootCertFiles \"${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt\" -c '{\"function\":\"UpdateDiagnosis\",\"Args\":[\"$id\"]}'"
    echo $diagnosisCommand
    eval $diagnosisCommand

    #sleep to allow the data to be udpated before getting and posting it
    sleep 1

    #get the data
    diagnosis=$(peer chaincode query -C mychannel -n basic -c "{\"function\":\"ReadDiagnosis\",\"Args\":[\"D$id\"]}")
    oxygenSaturationDiagnosis=$(echo $diagnosis | jq -r '.OxygenSaturationDiagnosis')
    pulseRateDiagnosis=$(echo $diagnosis | jq -r '.PulseRateDiagnosis')
    temperatureDiagnosis=$(echo $diagnosis | jq -r '.TemperatureDiagnosis')
    bloodPressureDiagnosis=$(echo $diagnosis | jq -r '.BloodPressureDiagnosis')

    #post it
    #curl -X POST -H "Content-Type: application/json" -d "{\"id\": \"$id\", \"oxygenSaturationDiagnosis\": \"$oxygenSaturationDiagnosis\", \"pulseRateDiagnosis\": \"$pulseRateDiagnosis\", \"temperatureDiagnosis\": \"$temperatureDiagnosis\", \"bloodPressureDiagnosis\": \"$bloodPressureDiagnosis\"}" http://localhost:4000/patients/$id/diagnosis

    #update each 5 seconds
    sleep 5


done
