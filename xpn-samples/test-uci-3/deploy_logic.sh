#!/bin/bash

#Authors
# Bogdan Gabriel Hortea
#

# Ask for number of patients
echo "Introduce number of patients"
#read numPatients
numPatients=150

# Validate number of patients
if ! [[ "$numPatients" =~ ^[0-9]+$ ]]; then
    echo "Error: Number of patients must be integer. "
    exit 1
fi

# Create patients
for ((i=1; i<=numPatients; i++)); do

    id="$i"
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
    command="peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile \"${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem\" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles \"${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt\" --peerAddresses localhost:9051 --tlsRootCertFiles \"${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt\" -c '{\"function\":\"CreatePatient\",\"Args\":[\"$id\", \"$firstName\", \"$lastName\", \"$birthDate\", \"$birthPlace\", \"$weight\", \"$height\"]}'"
    echo -e "GENERATING PATIENT WITH ID $id ... \n"
    echo $command
    eval $command
    # POST patient
    #curl -X POST -H "Content-Type: application/json" -d "{\"id\": \"$id\", \"firstName\": \"$firstName\", \"lastName\": \"$lastName\", \"birthDate\": \"$birthDate\", \"birthPlace\": \"$birthPlace\", \"weight\": \"$weight\", \"height\": \"$height\", \"oxygenSaturation\": \"$oxygenSaturation\", \"pulseRate\": \"$pudiagnosislseRate\", \"temperature\": \"$temperature\", \"bloodPressureSystolic\": \"$bloodPressureSystolic\", \"bloodPressureDiastolic\": \"$bloodPressureDiastolic\",  \"oxygenSaturationDiagnosis\": \"$oxygenSaturationDiagnosis\",  \"pulseRateDiagnosis\": \"$pulseRateDiagnosis\",  \"temperatureDiagnosis\": \"$temperatureDiagnosis\",  \"bloodPressureDiagnosis\": \"$bloodPressureDiagnosis\"}" http://localhost:4000/patients
    sleep 5

    # Create contract
    contractCommand="peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile \"${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem\" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles \"${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt\" --peerAddresses localhost:9051 --tlsRootCertFiles \"${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt\" -c '{\"function\":\"CreateContract\",\"Args\":[\"$id\", \"95\", \"100\", \"60\", \"100\", \"35.5\", \"38.0\", \"120\", \"180\", \"80\", \"120\"]}'"
    echo -e "GENERATING CONTRACT FOR PATIENT WITH ID $id ... \n "
    echo $contractCommand
    eval $contractCommand
    sleep 5

    # Create diagnosis
    diagnosisCommand="peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile \"${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem\" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles \"${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt\" --peerAddresses localhost:9051 --tlsRootCertFiles \"${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt\" -c '{\"function\":\"CreateDiagnosis\",\"Args\":[\"$id\"]}'"
    echo -e "GENERATING DIAGNOSIS FOR PATIENT WITH ID $id ... \n "
    echo $diagnosisCommand
    eval $diagnosisCommand
    sleep 5


    # Generate real time data for a patient in a subrocess in background
    ./generate_data.sh "$id" &

    sleep 10
    # Update diagnosis data for a patient in a subrocess in background
    ./update_diagnosis.sh "$id" &

done

wait