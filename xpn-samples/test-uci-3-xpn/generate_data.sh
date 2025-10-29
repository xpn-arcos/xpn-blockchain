#!/bin/bash
#Authors
# Bogdan Gabriel Hortea
# Diego Camarmas Alonso

patient_id="$1"
i=1
#endTime=$(($(date +%s) + 120)) #120 seconds generating data. Adjust to desired value
endTime=$(($(date +%s) + 43200)) #120 seconds generating data. Adjust to desired value

while [ $(date +%s) -lt $endTime ]; do

    transaction_id=0
    oxygenSaturation=$(shuf -i 90-100 -n 1) #between 90 an 100
    pulseRate=$(shuf -i 40-120 -n 1) #between 40 and 120
    temperature=$(LC_ALL=C awk "BEGIN {print 34 + 6 * $(shuf -i 0-100 -n 1) / 100}") #between 34 and 40
    bloodPressureSystolic=$(shuf -i 100-200 -n 1) #between 100 and 200
    bloodPressureDiastolic=$(shuf -i 60-140 -n 1) #between 60 and 140

    echo -e "GENERATING VITAL SIGNS FOR PATIENT WITH ID $patient_id ... \n "

    {
        flock -x 200
        transaction_id=$(cat "$TRANSACTION_ID")
        echo $((transaction_id+1)) > "$TRANSACTION_ID"
    } 200>"$LOCK_FILE"

    time {
        file_hash=$(LD_PRELOAD=$HOME/bin/xpn/lib/xpn_bypass.so python3 write_data_xpn.py /tmp/expand/xpn/$patient_id/vital_signs\_$i.txt "{\"id\":\"1\", \"patientId\":\"$patient_id\", \"oxygenSaturation\":\"$oxygenSaturation\", \"pulseRate\":\"$pulseRate\", \"temperature\":\"$temperature\", \"bloodPressureSystolic\":\"$bloodPressureSystolic\", \"bloodPressureDiastolic\":\"$bloodPressureDiastolic\"}")
        
        #file_hash=$(python3 write_data_xpn.py /tmp/expand/xpn/$patient_id/vital_signs\_$i.txt "{\"id\":\"$transaction_id\", \"patientId\":\"$patient_id\", \"oxygenSaturation\":\"$oxygenSaturation\", \"pulseRate\":\"$pulseRate\", \"temperature\":\"$temperature\", \"bloodPressureSystolic\":\"$bloodPressureSystolic\", \"bloodPressureDiastolic\":\"$bloodPressureDiastolic\"}")

        peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c "{\"function\":\"CreateXpnTransaction\",\"Args\":[\"$transaction_id\", \"$file_hash\", \"/tmp/expand/xpn/$patient_id/vital_signs_$i.txt\"]}"
    }

    i=$((i+1))

    #each 10 seconds, generate new data
    sleep 10
    
done
