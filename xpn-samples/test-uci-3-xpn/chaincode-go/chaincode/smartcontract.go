package chaincode

import (
	"encoding/json"
	"fmt"
	"strings"
	"regexp"
	"strconv"
	"unicode"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// SmartContract provides functions for managing a Patient
type SmartContract struct {
	contractapi.Contract
}

type XpnTransaction struct {
	ID          string  `json:"ID"`
	Hash        string  `json:"Hash"`
	Path        string  `json:"Path"`
}

type Patient struct {
	ID          string  `json:"ID"`
	FirstName   string  `json:"FirstName"`
	MiddleName  string  `json:"MiddleName"`
	LastName    string  `json:"LastName"`
	BirthDate   string  `json:"BirthDate"`
	BirthPlace  string  `json:"BirthPlace"`
	Weight 		float64 `json:"Weight"`
	Height 		float64 `json:"Height"`
}

type RTData struct {
	ID               		string   `json:"ID"`
	Patient					string 	 `json:"Patient"`
	OxygenSaturation 		float64  `json:"OxygenSaturation"`
	PulseRate        		float64  `json:"PulseRate"`
	Temperature      		float64  `json:"Temperature"`
	BloodPressureSystolic   float64  `json:"BloodPressureSystolic"`
	BloodPressureDiastolic  float64  `json:"BloodPressureDiastolic"`

}

type Diagnosis struct {
	ID 							string `json:"ID"`
	Patient						string `json:"Patient"`
	OxygenSaturationDiagnosis	string `json:"OxygenSaturationDiagnosis"`
	PulseRateDiagnosis        	string `json:"PulseRateDiagnosis"`
	TemperatureDiagnosis      	string `json:"TemperatureDiagnosis"`
	BloodPressureDiagnosis  	string `json:"BloodPressureDiagnosis"`
}


type Contract struct {
	ID          				string   `json:"ID"`
	Patient						string 	 `json:"Patient"`
	MinOxygenSaturation			float64  `json:"MinOxygenSaturation"`
	MaxOxygenSaturation			float64  `json:"MaxOxygenSaturation"`
	MinPulseRate				float64  `json:"MinPulseRate"`
	MaxPulseRate				float64  `json:"MaxPulseRate"`
	MinTemperature				float64  `json:"MinTemperature"`
	MaxTemperature				float64  `json:"MaxTemperature"`
	MinBloodPressureSystolic	float64  `json:"MinBloodPressureSystolic"`
	MaxBloodPressureSystolic	float64  `json:"MaxBloodPressureSystolic"`
	MinBloodPressureDiastolic	float64  `json:"MinBloodPressureDiastolic"`
	MaxBloodPressureDiastolic	float64  `json:"MaxBloodPressureDiastolic"`	
}

// ---------------------------------------------------- XpnTransaction -------------------------------------------------------------- //
// Validations for the XpnTransaction parameters
func (s *SmartContract) validateXpnTransaction(xpntransaction XpnTransaction) error {
	// Check if ID is empty
	if strings.TrimSpace(xpntransaction.ID) == "" {
		return fmt.Errorf("ID must be non-empty")
	}

	//check ID is a number
	if _, err := strconv.Atoi(xpntransaction.ID); err != nil {
		return fmt.Errorf("ID must be a number")
	}

	// Check if Hash is empty
	if strings.TrimSpace(xpntransaction.Hash) == "" {
		return fmt.Errorf("Hash must be non-empty")
	}

	// Check if Path is empty
	if strings.TrimSpace(xpntransaction.Path) == "" {
		return fmt.Errorf("Path must be non-empty")
	}

	return nil
}


// CreateXpnTransaction creates a new xpntransaction with given details.
func (s *SmartContract) CreateXpnTransaction(ctx contractapi.TransactionContextInterface, id string, hash string, path string) error {
	
	exists, err := s.XpnTransactionExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("Cannot create xpntransaction. XpnTransaction with id %s already exists", id)
	}


	xpntransaction := XpnTransaction{
		ID:   id,
		Hash: hash,
		Path: path,
	}

	// validate the xpntransaction
	err = s.validateXpnTransaction(xpntransaction)
	if err != nil {
		return err
	}
	
	xpntransactionJSON, err := json.Marshal(xpntransaction)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, xpntransactionJSON)
}


// ReadXpnTransaction returns the xpntransaction stored in the world state with given id.
func (s *SmartContract) ReadXpnTransaction(ctx contractapi.TransactionContextInterface, id string) (*XpnTransaction, error) {
	xpntransactionJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read xpntransaction from world state: %v", err)
	}
	if xpntransactionJSON == nil {
		return nil, fmt.Errorf("Cannot read xpntransaction. XpnTransaction with id %s does not exist", id)
	}

	var xpntransaction XpnTransaction
	err = json.Unmarshal(xpntransactionJSON, &xpntransaction)
	if err != nil {
		return nil, err
	}

	return &xpntransaction, nil
}

// XpnTransactionExists returns true when xpntransaction with given ID exists in world state.
func (s *SmartContract) XpnTransactionExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	xpntransactionJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read xpntransaction from world state: %v", err)
	}

	return xpntransactionJSON != nil, nil
}










// ---------------------------------------------------- PATIENTS -------------------------------------------------------------- //
// Validations for the patient parameters
func (s *SmartContract) validatePatient(patient Patient) error {
	// Check if ID is empty
	if strings.TrimSpace(patient.ID) == "" {
		return fmt.Errorf("ID must be non-empty")
	}

	//check ID is a number
	if _, err := strconv.Atoi(patient.ID); err != nil {
		return fmt.Errorf("ID must be a number")
	}

	// Check if FirstName or LastName are empty
	if strings.TrimSpace(patient.FirstName) == "" || strings.TrimSpace(patient.LastName) == "" {
		return fmt.Errorf("both FirstName and LastName must be non-empty")
	}

	// Check if BirthDate is in the correct format
	match, _ := regexp.MatchString(`^([0-2][0-9]|(3)[0-1])(\-)(((0)[0-9])|((1)[0-2]))(\-)\d{4}$`, patient.BirthDate)
	if !match {
		return fmt.Errorf("BirthDate is not in a valid format, required: dd-mm-yyyy")
	}

	// Check if Weight and Height are positive numbers
	if patient.Weight <= 0 || patient.Height <= 0 {
		return fmt.Errorf("Weight and Height must be positive numbers")
	}

	return nil
}


// CreatePatient creates a new patient with given details.
func (s *SmartContract) CreatePatient(ctx contractapi.TransactionContextInterface, id string, firstName string, middleName string, 
	lastName string, birthDate string, birthPlace string, weight string, height string) error {
	
	weightFloat, err := strconv.ParseFloat(weight, 64)
    if err != nil {
        return fmt.Errorf("invalid weight: %s", weight)
    }

    heightFloat, err := strconv.ParseFloat(height, 64)
    if err != nil {
        return fmt.Errorf("invalid height: %s", height)
    }
	
	
	exists, err := s.PatientExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("Cannot create patient. Patient with id %s already exists", id)
	}


	patient := Patient{
		ID:         id,
		FirstName:  firstName,
		MiddleName: middleName,
		LastName:   lastName,
		BirthDate:  birthDate,
		BirthPlace: birthPlace,
		Weight:		weightFloat,
		Height:		heightFloat,
	}

	// validate the patient
	err = s.validatePatient(patient)
	if err != nil {
		return err
	}
	
	patientJSON, err := json.Marshal(patient)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, patientJSON)
}


// ReadPatient returns the patient stored in the world state with given id.
func (s *SmartContract) ReadPatient(ctx contractapi.TransactionContextInterface, id string) (*Patient, error) {
	patientJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read patient from world state: %v", err)
	}
	if patientJSON == nil {
		return nil, fmt.Errorf("Cannot read patient. Patient with id %s does not exist", id)
	}

	var patient Patient
	err = json.Unmarshal(patientJSON, &patient)
	if err != nil {
		return nil, err
	}

	return &patient, nil
}

// UpdatePatient updates an existing patient in the world state with provided parameters.
func (s *SmartContract) UpdatePatient(ctx contractapi.TransactionContextInterface, id string, firstName string, middleName string, 
	lastName string, birthDate string, birthPlace string, weight string, height string) error {
	
	weightFloat, err := strconv.ParseFloat(weight, 64)
    if err != nil {
        return fmt.Errorf("invalid weight: %s", weight)
    }

    heightFloat, err := strconv.ParseFloat(height, 64)
    if err != nil {
        return fmt.Errorf("invalid height: %s", height)
    }

	exists, err := s.PatientExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Cannot update patient. Patient with id %s does not exist", id)
	}

	// overwriting original patient with new patient
	patient := Patient{
		ID:         id,
		FirstName:  firstName,
		MiddleName: middleName,
		LastName:   lastName,
		BirthDate:  birthDate,
		BirthPlace: birthPlace,
		Weight:		weightFloat,
		Height:		heightFloat,
	}

	// validate the patient
	err = s.validatePatient(patient)
	if err != nil {
		return err
	}
			
	patientJSON, err := json.Marshal(patient)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, patientJSON)
}

// DeletePatient deletes a patient from the world state.
func (s *SmartContract) DeletePatient(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.PatientExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Cannot delete patient. Patient with id %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// PatientExists returns true when patient with given ID exists in world state.
func (s *SmartContract) PatientExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	patientJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read patient from world state: %v", err)
	}

	return patientJSON != nil, nil
}


// function to check if first character is a digit
func startsWithDigit(s string) bool {
    for _, r := range s {
        return unicode.IsDigit(r)
    }
    return false
}

// GetAllPatients returns all patients found in world state
func (s *SmartContract) GetAllPatients(ctx contractapi.TransactionContextInterface) ([]*Patient, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var patients []*Patient
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		// Check if the ID starts with a digit before unmarshalling and adding it to the list
		if startsWithDigit(queryResponse.Key) {
			var patient Patient
			err = json.Unmarshal(queryResponse.Value, &patient)
			if err != nil {
				return nil, err
			}
			patients = append(patients, &patient)
		}
	}

	return patients, nil
}

// ---------------------------------------------------- CONTRACTS -------------------------------------------------------------- //
// Validations for the contract parameters
func (s *SmartContract) validateContract(contract Contract) error {
	// Check if ID is empty
	if strings.TrimSpace(contract.ID) == "" {
		return fmt.Errorf("ID must be non-empty")
	}

	// Check if patient is empty
	if strings.TrimSpace(contract.Patient) == "" {
		return fmt.Errorf("patient must be non-empty")
	}

	//check patient is a number
	if _, err := strconv.Atoi(contract.Patient); err != nil {
		return fmt.Errorf("patient must be numeric")
	}

	// Check if all other fields are positive numbers
	if contract.MinOxygenSaturation <= 0 || contract.MaxOxygenSaturation <= 0 || contract.MinPulseRate <= 0 || 
	contract.MaxPulseRate <= 0 || contract.MinTemperature <= 0 || contract.MaxTemperature <= 0 || 
	contract.MinBloodPressureSystolic <= 0 || contract.MaxBloodPressureSystolic <= 0 || 
	contract.MinBloodPressureDiastolic <= 0 || contract.MaxBloodPressureDiastolic <= 0{
		return fmt.Errorf("contract's maximum and minimum values must be positive numbers")
	}

	return nil
}


// CreateContract creates a new contract with given details.
func (s *SmartContract) CreateContract(ctx contractapi.TransactionContextInterface, patient string, minOxygenSaturation string, 
	maxOxygenSaturation string, minPulseRate string, maxPulseRate string, minTemperature string, maxTemperature string, 
	minBloodPressureSystolic string, maxBloodPressureSystolic string, minBloodPressureDiastolic string, 
	maxBloodPressureDiastolic string) error {

	minOxygenSaturationFloat, err := strconv.ParseFloat(minOxygenSaturation, 64)
    if err != nil {
        return fmt.Errorf("invalid minOxygenSaturation: %s", minOxygenSaturation)
    }
    maxOxygenSaturationFloat, err := strconv.ParseFloat(maxOxygenSaturation, 64)
    if err != nil {
        return fmt.Errorf("invalid maxOxygenSaturation: %s", maxOxygenSaturation)
    }
	minPulseRateFloat, err := strconv.ParseFloat(minPulseRate, 64)
    if err != nil {
        return fmt.Errorf("invalid minPulseRate: %s", minPulseRate)
    }
    maxPulseRateFloat, err := strconv.ParseFloat(maxPulseRate, 64)
    if err != nil {
        return fmt.Errorf("invalid maxPulseRate: %s", maxPulseRate)
    }
	minTemperatureFloat, err := strconv.ParseFloat(minTemperature, 64)
    if err != nil {
        return fmt.Errorf("invalid minTemperature: %s", minTemperature)
    }
    maxTemperatureFloat, err := strconv.ParseFloat(maxTemperature, 64)
    if err != nil {
        return fmt.Errorf("invalid maxTemperature: %s", maxTemperature)
    }
	minBloodPressureSystolicFloat, err := strconv.ParseFloat(minBloodPressureSystolic, 64)
    if err != nil {
        return fmt.Errorf("invalid minBloodPressureSystolic: %s", minBloodPressureSystolic)
    }
    maxBloodPressureSystolicFloat, err := strconv.ParseFloat(maxBloodPressureSystolic, 64)
    if err != nil {
        return fmt.Errorf("invalid maxBloodPressureSystolic: %s", maxBloodPressureSystolic)
    }
	minBloodPressureDiastolicFloat, err := strconv.ParseFloat(minBloodPressureDiastolic, 64)
    if err != nil {
        return fmt.Errorf("invalid minBloodPressureDiastolic: %s", minBloodPressureDiastolic)
    }
    maxBloodPressureDiastolicFloat, err := strconv.ParseFloat(maxBloodPressureDiastolic, 64)
    if err != nil {
        return fmt.Errorf("invalid maxBloodPressureDiastolic: %s", maxBloodPressureDiastolic)
    }

	//create id based on patient
	id := "C" + patient

	// check if a contract is already created for a patient
	contractExists, err := s.ContractExists(ctx, id)
	if err != nil {
		return err
	}
	if contractExists {
		return fmt.Errorf("Cannot create contract. Contract with id %s already exists", id)
	}

	// Check if a patient with the given ID exists.
	patientExists, err := s.PatientExists(ctx, patient)
	if err != nil {
		return err
	}
	if !patientExists {
		return fmt.Errorf("Cannot create contract. Patient with id %s does not exist", patient)
	}

	contract := Contract{
		ID:         		  		id,	
		Patient:					patient,
		MinOxygenSaturation:  		minOxygenSaturationFloat,
		MaxOxygenSaturation:  		maxOxygenSaturationFloat,
		MinPulseRate: 		  		minPulseRateFloat,
		MaxPulseRate:		  		maxPulseRateFloat,
		MinTemperature:		  		minTemperatureFloat,
		MaxTemperature:		  		maxTemperatureFloat,
		MinBloodPressureSystolic:	minBloodPressureSystolicFloat,
		MaxBloodPressureSystolic:	maxBloodPressureSystolicFloat,
		MinBloodPressureDiastolic:	minBloodPressureDiastolicFloat,
		MaxBloodPressureDiastolic:	maxBloodPressureDiastolicFloat,
	}

	// validate the contract
	err = s.validateContract(contract)
	if err != nil {
		return err
	}

	contractJSON, err := json.Marshal(contract)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, contractJSON)
}


// ReadContract returns the contract stored in the world state with given id.
func (s *SmartContract) ReadContract(ctx contractapi.TransactionContextInterface, id string) (*Contract, error) {
	contractJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read contract from world state: %v", err)
	}
	if contractJSON == nil {
		return nil, fmt.Errorf("Cannot read contract. Contract with id %s does not exist", id)
	}

	var contract Contract
	err = json.Unmarshal(contractJSON, &contract)
	if err != nil {
		return nil, err
	}

	return &contract, nil
}

// UpdateContract updates an existing contract in the world state with provided parameters.
func (s *SmartContract) UpdateContract(ctx contractapi.TransactionContextInterface, patient string, minOxygenSaturation string, 
	maxOxygenSaturation string, minPulseRate string, maxPulseRate string, minTemperature string, maxTemperature string, 
	minBloodPressureSystolic string, maxBloodPressureSystolic string, minBloodPressureDiastolic string, 
	maxBloodPressureDiastolic string) error {

	minOxygenSaturationFloat, err := strconv.ParseFloat(minOxygenSaturation, 64)
    if err != nil {
        return fmt.Errorf("invalid minOxygenSaturation: %s", minOxygenSaturation)
    }
    maxOxygenSaturationFloat, err := strconv.ParseFloat(maxOxygenSaturation, 64)
    if err != nil {
        return fmt.Errorf("invalid maxOxygenSaturation: %s", maxOxygenSaturation)
    }
	minPulseRateFloat, err := strconv.ParseFloat(minPulseRate, 64)
    if err != nil {
        return fmt.Errorf("invalid minPulseRate: %s", minPulseRate)
    }
    maxPulseRateFloat, err := strconv.ParseFloat(maxPulseRate, 64)
    if err != nil {
        return fmt.Errorf("invalid maxPulseRate: %s", maxPulseRate)
    }
	minTemperatureFloat, err := strconv.ParseFloat(minTemperature, 64)
    if err != nil {
        return fmt.Errorf("invalid minTemperature: %s", minTemperature)
    }
    maxTemperatureFloat, err := strconv.ParseFloat(maxTemperature, 64)
    if err != nil {
        return fmt.Errorf("invalid maxTemperature: %s", maxTemperature)
    }
	minBloodPressureSystolicFloat, err := strconv.ParseFloat(minBloodPressureSystolic, 64)
    if err != nil {
        return fmt.Errorf("invalid minBloodPressureSystolic: %s", minBloodPressureSystolic)
    }
    maxBloodPressureSystolicFloat, err := strconv.ParseFloat(maxBloodPressureSystolic, 64)
    if err != nil {
        return fmt.Errorf("invalid maxBloodPressureSystolic: %s", maxBloodPressureSystolic)
    }
	minBloodPressureDiastolicFloat, err := strconv.ParseFloat(minBloodPressureDiastolic, 64)
    if err != nil {
        return fmt.Errorf("invalid minBloodPressureDiastolic: %s", minBloodPressureDiastolic)
    }
    maxBloodPressureDiastolicFloat, err := strconv.ParseFloat(maxBloodPressureDiastolic, 64)
    if err != nil {
        return fmt.Errorf("invalid maxBloodPressureDiastolic: %s", maxBloodPressureDiastolic)
    }

	
	id := "C" + patient

	// check if the contract exists
	contractExists, err := s.ContractExists(ctx, id)
	if err != nil {
		return err
	}
	if !contractExists {
		return fmt.Errorf("Cannot update contract. Contract with id %s does not exist", id)
	}

	// Check if a patient with the given ID exists.
	patientExists, err := s.PatientExists(ctx, patient)
	if err != nil {
		return err
	}
	if !patientExists {
		return fmt.Errorf("Cannot update contract. Patient with id %s does not exist", patient)
	}

	// overwriting original contract with new contract
	contract := Contract{
		ID:         		  		id,	
		Patient:					patient,
		MinOxygenSaturation:  		minOxygenSaturationFloat,
		MaxOxygenSaturation:  		maxOxygenSaturationFloat,
		MinPulseRate: 		  		minPulseRateFloat,
		MaxPulseRate:		  		maxPulseRateFloat,
		MinTemperature:		  		minTemperatureFloat,
		MaxTemperature:		  		maxTemperatureFloat,
		MinBloodPressureSystolic:	minBloodPressureSystolicFloat,
		MaxBloodPressureSystolic:	maxBloodPressureSystolicFloat,
		MinBloodPressureDiastolic:	minBloodPressureDiastolicFloat,
		MaxBloodPressureDiastolic:	maxBloodPressureDiastolicFloat,
	}

	// validate the contract
	err = s.validateContract(contract)
	if err != nil {
		return err
	}

	contractJSON, err := json.Marshal(contract)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, contractJSON)
}

// DeleteContract deletes a contract from the world state.
func (s *SmartContract) DeleteContract(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.ContractExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Cannot delete contract. Contract with id %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}


// ContractExists returns true when contract with given ID exists in world state.
func (s *SmartContract) ContractExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	contractJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read contract from world state: %v", err)
	}

	return contractJSON != nil, nil
}


// GetAllContracts returns all contracts found in world state
func (s *SmartContract) GetAllContracts(ctx contractapi.TransactionContextInterface) ([]*Contract, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var contracts []*Contract
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		if strings.HasPrefix(queryResponse.Key, "C") {
			var contract Contract
			err = json.Unmarshal(queryResponse.Value, &contract)
			if err != nil {
				return nil, err
			}
			contracts = append(contracts, &contract)
		}
	}

	return contracts, nil
}


// ------------------------------------------------ REAL TIME DATA --------------------------------------------------------- //
// Validations for the RTData parameters
func (s *SmartContract) validateRTData(rtdata RTData) error {
	// Check if ID is empty
	if strings.TrimSpace(rtdata.ID) == "" {
		return fmt.Errorf("ID must be non-empty")
	}

	// Check if patient is empty
	if strings.TrimSpace(rtdata.Patient) == "" {
		return fmt.Errorf("patient must be non-empty")
	}

	//check patient is a number
	if _, err := strconv.Atoi(rtdata.Patient); err != nil {
		return fmt.Errorf("patient must be numeric")
	}

	// Check if all other fields are positive numbers 
	if rtdata.OxygenSaturation <= 0 || rtdata.PulseRate <= 0 || rtdata.Temperature <= 0 || 
	rtdata.BloodPressureSystolic <= 0 || rtdata.BloodPressureDiastolic <= 0 {
		return fmt.Errorf("real time measurements must be positive numbers")
	}

	return nil
}


// CreateRTData creates new real time measurements with given details.
func (s *SmartContract) CreateRTData(ctx contractapi.TransactionContextInterface, patient string, oxygenSaturation string,
	pulseRate string, temperature string, bloodPressureSystolic string, bloodPressureDiastolic string) error {

	oxygenSaturationFloat, err := strconv.ParseFloat(oxygenSaturation, 64)
    if err != nil {
        return fmt.Errorf("invalid oxygenSaturation: %s", oxygenSaturation)
    }
    pulseRateFloat, err := strconv.ParseFloat(pulseRate, 64)
    if err != nil {
        return fmt.Errorf("invalid pulseRate: %s", pulseRate)
    }
	temperatureFloat, err := strconv.ParseFloat(temperature, 64)
    if err != nil {
        return fmt.Errorf("invalid temperature: %s", temperature)
    }
    bloodPressureSystolicFloat, err := strconv.ParseFloat(bloodPressureSystolic, 64)
    if err != nil {
        return fmt.Errorf("invalid bloodPressureSystolic: %s", bloodPressureSystolic)
    }
	bloodPressureDiastolicFloat, err := strconv.ParseFloat(bloodPressureDiastolic, 64)
    if err != nil {
        return fmt.Errorf("invalid bloodPressureDiastolic: %s", bloodPressureDiastolic)
    }

	id := "RTD" + patient

	rtDataExists, err := s.RTDataExists(ctx, id)
	if err != nil {
		return err
	}
	if rtDataExists {
		return fmt.Errorf("Cannot create real time measurements. Measurements with id %s already exists", id)
	}

	// Check if a patient with the given ID exists.
	patientExists, err := s.PatientExists(ctx, patient)
	if err != nil {
		return err
	}
	if !patientExists {
		return fmt.Errorf("Cannot create measurements. Patient with id %s does not exist", patient)
	}

	rtData:= RTData{
		ID:         	  		 id,
		Patient:			     patient,
		OxygenSaturation: 		 oxygenSaturationFloat,
		PulseRate: 		  		 pulseRateFloat,
		Temperature:   	  		 temperatureFloat,
		BloodPressureSystolic:   bloodPressureSystolicFloat,
		BloodPressureDiastolic:  bloodPressureDiastolicFloat,
	}

	// validate the RTData
	err = s.validateRTData(rtData)
	if err != nil {
		return err
	}

	rtDataJSON, err := json.Marshal(rtData)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, rtDataJSON)
}


// ReadRTData returns real time measurements for a patient stored in the world state with given id.
func (s *SmartContract) ReadRTData(ctx contractapi.TransactionContextInterface, id string) (*RTData, error) {
	rtDataJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read real time data from world state: %v", err)
	}
	if rtDataJSON == nil {
		return nil, fmt.Errorf("Cannot read measurements. Measurements with id %s do not exist", id)
	}

	var rtData RTData
	err = json.Unmarshal(rtDataJSON, &rtData)
	if err != nil {
		return nil, err
	}

	return &rtData, nil
}


// UpdateRTData updates existing real time measurements for a patient in the world state with provided parameters.
func (s *SmartContract) UpdateRTData(ctx contractapi.TransactionContextInterface, patient string, oxygenSaturation string,
	pulseRate string, temperature string, bloodPressureSystolic string, bloodPressureDiastolic string) error {

	oxygenSaturationFloat, err := strconv.ParseFloat(oxygenSaturation, 64)
    if err != nil {
        return fmt.Errorf("invalid oxygenSaturation: %s", oxygenSaturation)
    }
    pulseRateFloat, err := strconv.ParseFloat(pulseRate, 64)
    if err != nil {
        return fmt.Errorf("invalid pulseRate: %s", pulseRate)
    }
	temperatureFloat, err := strconv.ParseFloat(temperature, 64)
    if err != nil {
        return fmt.Errorf("invalid temperature: %s", temperature)
    }
    bloodPressureSystolicFloat, err := strconv.ParseFloat(bloodPressureSystolic, 64)
    if err != nil {
        return fmt.Errorf("invalid bloodPressureSystolic: %s", bloodPressureSystolic)
    }
	bloodPressureDiastolicFloat, err := strconv.ParseFloat(bloodPressureDiastolic, 64)
    if err != nil {
        return fmt.Errorf("invalid bloodPressureDiastolic: %s", bloodPressureDiastolic)
    }

	id := "RTD" + patient

	rtDataexists, err := s.RTDataExists(ctx, id)
	if err != nil {
		return err
	}
	if !rtDataexists {
		return fmt.Errorf("Cannot update measurements. Measurements with id %s do not exist", id)
	}

	// Check if a patient with the given ID exists.
	patientExists, err := s.PatientExists(ctx, patient)
	if err != nil {
		return err
	}
	if !patientExists {
		return fmt.Errorf("Cannot update measurements. Patient with id %s does not exist", patient)
	}

	// overwriting original measurements with new measurements
	rtData:= RTData{
		ID:         	  		 id,
		Patient:				 patient,
		OxygenSaturation: 		 oxygenSaturationFloat,
		PulseRate: 		  		 pulseRateFloat,
		Temperature:   	  		 temperatureFloat,
		BloodPressureSystolic:   bloodPressureSystolicFloat,
		BloodPressureDiastolic:  bloodPressureDiastolicFloat,
	}

	// validate the RTData
	err = s.validateRTData(rtData)
	if err != nil {
		return err
	}

	rtDataJSON, err := json.Marshal(rtData)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, rtDataJSON)

}


// RTDataExists returns true when real time measurements with given ID exists in world state.
func (s *SmartContract) RTDataExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	rtDataJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read real time data from world state: %v", err)
	}

	return rtDataJSON != nil, nil
}


// DeleteRTData deletes a contract from the world state.
func (s *SmartContract) DeleteRTData(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.RTDataExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Cannot delete measurements. Measurements with id %s do not exist", id)
	}

	return ctx.GetStub().DelState(id)
}


// GetAllRTData returns all contracts found in world state
func (s *SmartContract) GetAllRTData(ctx contractapi.TransactionContextInterface) ([]*RTData, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var measurements []*RTData
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		if strings.HasPrefix(queryResponse.Key, "R") {
			var rtData RTData
			err = json.Unmarshal(queryResponse.Value, &rtData)
			if err != nil {
				return nil, err
			}
			measurements = append(measurements, &rtData)
		}
	}

	return measurements, nil
}


// ------------------------------------------------ DIAGNOSIS --------------------------------------------------------- //
func (s *SmartContract) validateDiagnosis(diagnosis Diagnosis) error {
	// Check if ID is empty
	if strings.TrimSpace(diagnosis.ID) == "" {
		return fmt.Errorf("ID must be non-empty")
	}

	// Check if patient is empty
	if strings.TrimSpace(diagnosis.Patient) == "" {
		return fmt.Errorf("patient must be non-empty")
	}

	//check patient is a number
	if _, err := strconv.Atoi(diagnosis.Patient); err != nil {
		return fmt.Errorf("patient must be numeric")
	}

	return nil
}


// CreateDiagnosis creates an empty diagnosis for a patient with given id
func (s *SmartContract) CreateDiagnosis(ctx contractapi.TransactionContextInterface, patient string) error {
	var oxygenSaturationDiagnosis string
	oxygenSaturationDiagnosis = "None"
	var pulseRateDiagnosis string
	pulseRateDiagnosis = "None"
	var temperatureDiagnosis string
	temperatureDiagnosis = "None"
	var bloodPressureDiagnosis string
	bloodPressureDiagnosis = "None"
	id := "D" + patient

	
	diagnosisExists, err := s.DiagnosisExists(ctx, id)
	if err != nil {
		return err
	}
	if diagnosisExists {
		return fmt.Errorf("Cannot create diagnosis. Diagnosis with id %s already exists", id)
	}

	// Check if a patient with the given ID exists.
	patientExists, err := s.PatientExists(ctx, patient)
	if err != nil {
		return err
	}
	if !patientExists {
		return fmt.Errorf("Cannot create diagnosis. Patient with id %s does not exist", patient)
	}
	
	diagnosis:= Diagnosis{
		ID:        					id,
		Patient:					patient,
		OxygenSaturationDiagnosis:  oxygenSaturationDiagnosis,	
		PulseRateDiagnosis:     	pulseRateDiagnosis,
		TemperatureDiagnosis:       temperatureDiagnosis,
		BloodPressureDiagnosis:  	bloodPressureDiagnosis,
		
	}

	// validate diagnosis
	err = s.validateDiagnosis(diagnosis)
	if err != nil {
		return err
	}

	diagnosisJSON, err := json.Marshal(diagnosis)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, diagnosisJSON)

}


// ReadDiagnosis returns the diagnosis for a patient stored in the world state with given id.
func (s *SmartContract) ReadDiagnosis(ctx contractapi.TransactionContextInterface, id string) (*Diagnosis, error) {
	diagnosisJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read diagnosis from world state: %v", err)
	}
	if diagnosisJSON == nil {
		return nil, fmt.Errorf("Cannot read diagnosis. Diagnosis with id %s does not exist", id)
	}

	var diagnosis Diagnosis
	err = json.Unmarshal(diagnosisJSON, &diagnosis)
	if err != nil {
		return nil, err
	}

	return &diagnosis, nil
}


// UpdateDiagnosis checks real time data for a specific patient and issue the diagnosis and 
// an alert if any measurement is outside the contract limits.
func (s *SmartContract) UpdateDiagnosis(ctx contractapi.TransactionContextInterface, patient string) error {
	contractID := "C" + patient
	contract, err := s.ReadContract(ctx, contractID)
	if err != nil {
		return fmt.Errorf("Could not read contract: %s", err.Error())
	}

	rtdID := "RTD" + patient
	rtData, err := s.ReadRTData(ctx, rtdID)
	if err != nil {
		return fmt.Errorf("Could not read RTData: %s", err.Error())
	}

	id := "D" + patient

	diagnosis := Diagnosis{ID: id, Patient: patient}

	if rtData.OxygenSaturation < contract.MinOxygenSaturation {
		diagnosis.OxygenSaturationDiagnosis = "Low oxygen saturation"
	} else if rtData.OxygenSaturation > contract.MaxOxygenSaturation {
		diagnosis.OxygenSaturationDiagnosis = "High oxygen saturation"
	} else {
		diagnosis.OxygenSaturationDiagnosis = "Oxygen Saturation in correct range"
	}

	if rtData.PulseRate < contract.MinPulseRate {
		diagnosis.PulseRateDiagnosis = "Alert. Bradycardia"
	} else if rtData.PulseRate > contract.MaxPulseRate {
		diagnosis.PulseRateDiagnosis = "Alert. Tachycardia"
	} else {
		diagnosis.PulseRateDiagnosis = "Pulse rate in correct range"
	}

	if rtData.Temperature < contract.MinTemperature {
		diagnosis.TemperatureDiagnosis = "Alert. Low body temperature"
	} else if rtData.Temperature > contract.MaxTemperature {
		diagnosis.TemperatureDiagnosis = "Alert. Fever"
	} else {
		diagnosis.TemperatureDiagnosis = "Temperature in correct range"
	}

	if rtData.BloodPressureSystolic < contract.MinBloodPressureSystolic && 
	rtData.BloodPressureDiastolic < contract.MinBloodPressureDiastolic {
		diagnosis.BloodPressureDiagnosis = "Normal blood pressure"
	} else if rtData.BloodPressureSystolic >= contract.MinBloodPressureSystolic && rtData.BloodPressureSystolic < 130 && 
	rtData.BloodPressureDiastolic < contract.MinBloodPressureDiastolic {
		diagnosis.BloodPressureDiagnosis = "Elevated blood pressure"
	} else if (rtData.BloodPressureSystolic >= 130 && rtData.BloodPressureSystolic < 140) || 
	(rtData.BloodPressureDiastolic >= contract.MinBloodPressureDiastolic && rtData.BloodPressureDiastolic < 89) {
		diagnosis.BloodPressureDiagnosis = "Hypertension Stage 1"
	} else if rtData.BloodPressureSystolic >= 140 || rtData.BloodPressureDiastolic >= 90 {
		diagnosis.BloodPressureDiagnosis = "Hypertension Stage 2"
	} else if rtData.BloodPressureSystolic > contract.MaxBloodPressureSystolic || 
	rtData.BloodPressureDiastolic > contract.MaxBloodPressureSystolic {
		diagnosis.BloodPressureDiagnosis = "Hypertensive crisis - Immediate medical attention required"
	}
	

	diagnosisExists, err := s.DiagnosisExists(ctx, id)
	if err != nil {
		return err
	}
	if !diagnosisExists {
		return fmt.Errorf("Cannot update diagnosis. Diagnosis with id %s does not exist", id)
	}

	// Check if a patient with the given ID exists.
	patientExists, err := s.PatientExists(ctx, patient)
	if err != nil {
		return err
	}
	if !patientExists {
		return fmt.Errorf("Cannot update diagnosis. Patient with id %s does not exist", patient)
	}

	// validate diagnosis
	err = s.validateDiagnosis(diagnosis)
	if err != nil {
		return err
	}

	diagnosisJSON, err := json.Marshal(diagnosis)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, diagnosisJSON)

}


// DiagnosisExists returns true when diagnosis with given ID exists in world state.
func (s *SmartContract) DiagnosisExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	diagnosisJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read diagnosis from world state: %v", err)
	}

	return diagnosisJSON != nil, nil
}


// DeleteDiagnosis deletes a contract from the world state.
func (s *SmartContract) DeleteDiagnosis(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.DiagnosisExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Cannot delete diagnosis. Diagnosis with id %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}


// GetAllDiagnosis returns all diagnosis found in world state
func (s *SmartContract) GetAllDiagnosis(ctx contractapi.TransactionContextInterface) ([]*Diagnosis, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var allDiagnosis []*Diagnosis
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		// Check if the ID starts with "D" before unmarshalling and adding it to the list
		if strings.HasPrefix(queryResponse.Key, "D") {
			var diagnosis Diagnosis
			err = json.Unmarshal(queryResponse.Value, &diagnosis)
			if err != nil {
				return nil, err
			}
			allDiagnosis = append(allDiagnosis, &diagnosis)
		}
	}

	return allDiagnosis, nil
}

