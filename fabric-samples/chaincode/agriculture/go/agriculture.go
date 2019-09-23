/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"fmt"
	"bytes"
	"encoding/json"
	"strconv"
	_ "math/rand"
	_ "time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {

}

/*
 * Define the Agricultural Products Traceability System structure, with 3 properties.
 * Structure tags are used by encoding/json library.
 */
type Product struct {
	GPS_Location	string	`json:"location"`
	Temperature	string	`json:"temperature"`
	Humidity	string	`json:"humidity"`
}

type Product2 struct {
	GPS_Location	int	`json:"location"`
	Temperature	int	`json:"temperature"`
	Humidity	int	`json:"humidity"`
}

/*
 * The Init method is called when the Smart Contract "agriculture" is instantiated by the blockchain network
 * Best practice is to have any Ledger initalization in separate function --see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	// Retrieve the requested Smart Contract function and atguments
	function, args := APIstub.GetFunctionAndParameters()

	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryProduct" {
		return s.queryProduct(APIstub, args)
	} else if function == "queryAllProducts" {
		return s.queryAllProducts(APIstub)
	} else if function == "createProduct" {
		return s.createProduct(APIstub, args)
	} else if function == "changeProductLocation" {
		return s.changeProductLocation(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} /*else if function == "testSubmit" {
		return s.testSubmit(APIstub)
	}*/

	return shim.Error("Invalid Smart Contract function name.")
}

// Used to query the specific product information from blockchain ledger
func (s *SmartContract) queryProduct(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	productAsBytes, _ := APIstub.GetState(args[0])

	return shim.Success(productAsBytes)
}

// Used to query all the product information from blockchain ledger
func (s *SmartContract) queryAllProducts(APIstub shim.ChaincodeStubInterface) sc.Response {
	startKey := "No0"
	endKey := "No5000"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err!= nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(",\"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
	}

	buffer.WriteString("]")

	fmt.Printf("-queryAllProducts:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// Used to add a new transaction (product information) to blockchain ledger
func (s *SmartContract) createProduct(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	var product = Product{GPS_Location: args[1], Temperature: args[2], Humidity: args[3]}

	productAsBytes, _ := json.Marshal(product)
	APIstub.PutState(args[0], productAsBytes)

	return shim.Success(nil)
}

// Used to change the existing product information from blockchain ledger
func (s *SmartContract) changeProductLocation(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrct number of arguments. Expecting 2")
	}

	productAsBytes, _ := APIstub.GetState(args[0])
	product := Product{}

	json.Unmarshal(productAsBytes, &product)
	product.GPS_Location = args[1]

	productAsBytes, _ = json.Marshal(product)
	APIstub.PutState(args[0], productAsBytes)

	return shim.Success(nil)
}

// Used to initialize the blockchain ledger
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	products := []Product {
		Product{GPS_Location: "(41.40338, 2.17403)", Temperature: "26.3", Humidity: "70.1"},
		Product{GPS_Location: "(53.22449, 3.79878)", Temperature: "27.8", Humidity: "32.3"},
		Product{GPS_Location: "(49.12345, 2.66534)", Temperature: "28.5", Humidity: "50.0"},
	}

	i := 0
	for i < len(products) {
		fmt.Println("i is ", i)
		productAsBytes, _ := json.Marshal(products[i])
		APIstub.PutState("No" + strconv.Itoa(i), productAsBytes)
		fmt.Println("Added", products[i])
		i = i + 1
	}

	return shim.Success(nil)
}

/*
// Used to submit lots of transaction
func (s *SmartContract) testSubmit(APIstub shim.ChaincodeStubInterface) sc.Response {

	var products [9000]Product2
	//var longitude_string string
	//var latitude_string string
	//var location_string string
	//var temperature_string string
	//var humidity_string string
	var temperature_int int
	var location_int int
	var humidity_int int

	i := 0
	for i < len(products) {
		//longitude_string = fmt.Sprintf("%.6f", randomProduceNumber("longitude"))
		//latitude_string = fmt.Sprintf("%.6f", randomProduceNumber("latitude"))
		//temperature_string = fmt.Sprintf("%.1f", randomProduceNumber("temperature"))
		//humidity_string = fmt.Sprintf("%.1f", randomProduceNumber("humidity"))
		//location_string = "(" + longitude_string + ", " + latitude_string + ")"
		temperature_int = rand.Intn(80) - 40
		location_int = rand.Intn(360) - 180
		humidity_int = rand.Intn(100)
		//temperature_string = fmt.Sprintf("%.1f", randomProduceNumber("temperature"))
		//temperature_float, err := strconv.ParseFloat(temperature_string, 64)
		//checkError(err)

		products[i] = Product2{GPS_Location: location_int, Temperature: temperature_int, Humidity: humidity_int}
		//fmt.Printf("%d\n", i)

		fmt.Println("i is ", i)
		productAsBytes, _ := json.Marshal(products[i])
		APIstub.PutState("No" + strconv.Itoa(i), productAsBytes)
		fmt.Println("Added", products[i])

		i = i + 1
	}

	return shim.Success(nil)
}

func randomProduceNumber(str string) float64 {
	rand.Seed(time.Now().UnixNano())

	switch str {
		case "longitude":
			num1 := float64(rand.Intn(360) - 180)
			num2, err := strconv.ParseFloat(fmt.Sprintf("%.6f", rand.Float32()), 64)
			checkError(err)

			return num1 + num2
			
		case "latitude":
			num1 := float64(rand.Intn(180) - 90)
			num2, err := strconv.ParseFloat(fmt.Sprintf("%.6f", rand.Float32()), 64)
			checkError(err)

			return num1 + num2

		case "temperature":
			num1 := float64(rand.Intn(70) - 35)
			num2, err := strconv.ParseFloat(fmt.Sprintf("%.1f", rand.Float32()), 64)
			checkError(err)

			return num1 + num2

		case "humidity":
			num1 := float64(rand.Intn(70) - 35)
			num2, err := strconv.ParseFloat(fmt.Sprintf("%.1f", rand.Float32()), 64)
			checkError(err)

			return num1 + num2

		default:
			return 0
	}
}
*/

func checkError(err error) {
	if err != nil {
		// panic(err)
		fmt.Printf("checkError: %s\n", err)
	}
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {
	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	checkError(err)
}

