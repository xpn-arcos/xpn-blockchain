#!/bin/bash
#Authors
# Bogdan Gabriel Hortea
#

id="$1"
i=1
endTime=$(($(date +%s) + 120)) #120 seconds generating data. Adjust to desired value

while [ $(date +%s) -lt $endTime ]; do

    oxygenSaturation=$(shuf -i 90-100 -n 1) #between 90 an 100
    pulseRate=$(shuf -i 40-120 -n 1) #between 40 and 120
    temperature=$(LC_ALL=C awk "BEGIN {print 34 + 6 * $(shuf -i 0-100 -n 1) / 100}") #between 34 and 40
    bloodPressureSystolic=$(shuf -i 100-200 -n 1) #between 100 and 200
    bloodPressureDiastolic=$(shuf -i 60-140 -n 1) #between 60 and 140

    if (( $(bc <<< "$i == 1") )); then
        rtDataCommand="peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile \"${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem\" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles \"${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt\" --peerAddresses localhost:9051 --tlsRootCertFiles \"${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt\" -c '{\"function\":\"CreateRTData\",\"Args\":[\"$id\", \"$oxygenSaturation\", \"$pulseRate\", \"$temperature\", \"$bloodPressureSystolic\", \"$bloodPressureDiastolic\"]}'"
    else
        rtDataCommand="peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile \"${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem\" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles \"${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt\" --peerAddresses localhost:9051 --tlsRootCertFiles \"${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt\" -c '{\"function\":\"UpdateRTData\",\"Args\":[\"$id\", \"$oxygenSaturation\", \"$pulseRate\", \"$temperature\", \"$bloodPressureSystolic\", \"$bloodPressureDiastolic\"]}'"
    fi

    echo $rtDataCommand
    time eval $rtDataCommand
    #curl -X POST -H "Content-Type: application/json" -d "{\"id\": \"$id\", \"oxygenSaturation\": \"$oxygenSaturation\", \"pulseRate\": \"$pulseRate\", \"temperature\": \"$temperature\", \"bloodPressureSystolic\": \"$bloodPressureSystolic\", \"bloodPressureDiastolic\": \"$bloodPressureDiastolic\"}" http://localhost:4000/patients/$id/real-time-data

    i=$((i+1))

    #each 10 seconds, generate new data
    sleep 10
    
done
