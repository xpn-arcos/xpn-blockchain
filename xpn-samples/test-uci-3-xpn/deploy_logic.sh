#!/bin/bash

#Authors
# Bogdan Gabriel Hortea
# Diego Camarmas Alonso

# Expand configuration 
export XPN_LOCALITY=0
export XPN_CONF=$HOME/xpn-blockchain/xpn-samples/test-uci-3-xpn/xpn/config.xml
#ssh -L *:3456:localhost:3456 -L *:3457:localhost:3457 tester005@10.119.12.168 &

#sleep 30

# Ask for number of patients
echo "Introduce number of patients"
#read numPatients
numPatients=150

# Validate number of patients
if ! [[ "$numPatients" =~ ^[0-9]+$ ]]; then
    echo "Error: Number of patients must be integer. "
    exit 1
fi

# Transaction ID
export TRANSACTION_ID="/tmp/transaction_id.var"
export LOCK_FILE="/tmp/transaction_id.lock"
echo 1 > "$TRANSACTION_ID"


# Create patients
for ((i=1; i<=numPatients; i++)); do

    transaction_id=0
    patient_id="$i"
    firstName=$(shuf -n 1 patient_data/first-names.txt)
    lastName=$(shuf -n 1 patient_data/last-names.txt)
    birthPlace=$(shuf -n 1 patient_data/countries.txt) 
    oxygenSaturation=none
    pulseRate=none
    temperature=none
    bloodPressureSystolic=none
    bloodPressureDiastolic=none
    oxygenSaturationDiagnosis=none
    pulseRateDiagnosis=none
    temperatureDiagnosis=none
    bloodPressureDiagnosis=none

    # Random birthdate
    year=$(( RANDOM % 80 + 1920 )) # year between 1920 y 2000
    month=$(( RANDOM % 12 + 1 ))    # month between 1 y 12

    # determine number of days in a month
    case $month in
        4|6|9|11) daysInMonth=30;;
        2) 
            if (( year % 4 == 0 && (year % 100 != 0 || year % 400 == 0) )); then
                daysInMonth=29 
            else
                daysInMonth=28
            fi
            ;;
        *) daysInMonth=31;;
    esac
    day=$(( RANDOM % daysInMonth + 1 ))


    birthDate=$(printf "%02d-%02d-%04d" $day $month $year)

    #Randon weight and height
    height=$(bc <<< "scale=2; 1.4 + 0.8 * $RANDOM / 32768")  # Between 1.4 and 2.2m
    weight=$(( RANDOM % 101 + 50 ))      





    # Create patient
    echo -e "GENERATING PATIENT WITH ID $patient_id ... \n"

    {
        flock -x 200
        transaction_id=$(cat "$TRANSACTION_ID")
        echo $((transaction_id+1)) > "$TRANSACTION_ID"
    } 200>"$LOCK_FILE"

    LD_PRELOAD=$HOME/bin/xpn/lib/xpn_bypass.so mkdir /tmp/expand/xpn/$patient_id
    file_hash=$(LD_PRELOAD=$HOME/bin/xpn/lib/xpn_bypass.so python3 write_data_xpn.py /tmp/expand/xpn/$patient_id/patient.txt "{\"id\":\"1\", \"patientId\":\"$patient_id\", \"firstName\":\"$firstName\", \"lastName\":\"$lastName\", \"birthDate\":\"$birthDate\", \"birthPlace\":\"$birthPlace\", \"weight\":\"$weight\", \"height\":\"$height\"}")
    
    #mkdir /tmp/expand/xpn/$patient_id
    #file_hash=$(python3 write_data_xpn.py /tmp/expand/xpn/$patient_id/patient.txt "{\"id\":\"$transaction_id\", \"patientId\":\"$patient_id\", \"firstName\":\"$firstName\", \"lastName\":\"$lastName\", \"birthDate\":\"$birthDate\", \"birthPlace\":\"$birthPlace\", \"weight\":\"$weight\", \"height\":\"$height\"}")

    peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c "{\"function\":\"CreateXpnTransaction\",\"Args\":[\"$transaction_id\", \"$file_hash\", \"/tmp/expand/xpn/$patient_id/patient.txt\"]}"

    sleep 5





    # Create contract
    echo -e "GENERATING CONTRACT FOR PATIENT WITH ID $patient_id ... \n "

    {
        flock -x 200
        transaction_id=$(cat "$TRANSACTION_ID")
        echo $((transaction_id+1)) > "$TRANSACTION_ID"
    } 200>"$LOCK_FILE"

    file_hash=$(LD_PRELOAD=$HOME/bin/xpn/lib/xpn_bypass.so python3 write_data_xpn.py /tmp/expand/xpn/$patient_id/contract.txt "{\"id\":\"1\", \"patientId\":\"$patient_id\", \"minOxygenSaturation\":\"95\", \"maxOxygenSaturation\":\"100\", \"minPulseRate\":\"60\", \"maxPulseRate\":\"100\", \"minTemperature\":\"35.5\", \"maxTemperature\":\"38\", \"minBloodPressureSystolic\":\"120\", \"maxBloodPressureSystolic\":\"180\", \"minBloodPressureDiastolic\":\"80\", \"maxBloodPressureDiastolic\":\"120\"}")
    
    #file_hash=$(python3 write_data_xpn.py /tmp/expand/xpn/$patient_id/contract.txt "{\"id\":\"$transaction_id\", \"patientId\":\"$patient_id\", \"minOxygenSaturation\":\"95\", \"maxOxygenSaturation\":\"100\", \"minPulseRate\":\"60\", \"maxPulseRate\":\"100\", \"minTemperature\":\"35.5\", \"maxTemperature\":\"38\", \"minBloodPressureSystolic\":\"120\", \"maxBloodPressureSystolic\":\"180\", \"minBloodPressureDiastolic\":\"80\", \"maxBloodPressureDiastolic\":\"120\"}")

    peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c "{\"function\":\"CreateXpnTransaction\",\"Args\":[\"$transaction_id\", \"$file_hash\", \"/tmp/expand/xpn/$patient_id/contract.txt\"]}"

    sleep 5





    # Create diagnosis
    echo -e "GENERATING DIAGNOSIS FOR PATIENT WITH ID $patient_id ... \n "

    {
        flock -x 200
        transaction_id=$(cat "$TRANSACTION_ID")
        echo $((transaction_id+1)) > "$TRANSACTION_ID"
    } 200>"$LOCK_FILE"

    file_hash=$(LD_PRELOAD=$HOME/bin/xpn/lib/xpn_bypass.so python3 write_data_xpn.py /tmp/expand/xpn/$patient_id/diagnosis.txt "{\"id\":\"1\", \"patientId\":\"$patient_id\", \"oxygenSaturationDiagnosis\":\"None\", \"pulseRateDiagnosis\":\"None\", \"temperatureDiagnosis\":\"None\", \"bloodPressureDiagnosis\":\"None\"}")
    
    #file_hash=$(python3 write_data_xpn.py /tmp/expand/xpn/$patient_id/diagnosis.txt "{\"id\":\"$transaction_id\", \"patientId\":\"$patient_id\", \"oxygenSaturationDiagnosis\":\"None\", \"pulseRateDiagnosis\":\"None\", \"temperatureDiagnosis\":\"None\", \"bloodPressureDiagnosis\":\"None\"}")

    peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c "{\"function\":\"CreateXpnTransaction\",\"Args\":[\"$transaction_id\", \"$file_hash\", \"/tmp/expand/xpn/$patient_id/diagnosis.txt\"]}"
    
    sleep 5


    # Generate real time data for a patient in a subrocess in background
    ./generate_data.sh "$patient_id" &





    #sleep 10
    # Update diagnosis data for a patient in a subrocess in background
    #./update_diagnosis.sh "$patient_id" &

done

wait