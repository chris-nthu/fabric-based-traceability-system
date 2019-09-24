package main

import (
	"fmt"
	"math/rand"
	"time"
	"strconv"
	_"reflect"
	"os"
	"bufio"
)

type Product struct {
	GPS_Location	string	`json:"location"`
	Temperature	string	`json:"temperature"`
	Humidity	string	`json:"humidity"`
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
			num1 := float64(rand.Intn(100))
			num2, err := strconv.ParseFloat(fmt.Sprintf("%.1f", rand.Float32()), 64)
			checkError(err)

			return num1 + num2

		default:
			return 0
	}
}

func checkError(err error) {
	if err != nil {
		// panic(err)
		fmt.Printf("checkError: %s\n", err)
	}
}

func main() {
	var products [1000]Product
	var key_string string
	var longitude_string string
	var latitude_string string
	var temperature_string string
	var humidity_string string

	outputf, err := os.OpenFile("file.csv", os.O_CREATE|os.O_WRONLY, 0664)
	checkError(err)

	defer func() {
		if err := outputf.Close(); err != nil {
			panic(err)
		}
	}()

	outputWriter := bufio.NewWriter(outputf)

	i := 0
	for i < len(products) {
		key_string = "No" + strconv.Itoa(i)
		longitude_string = fmt.Sprintf("%.6f", randomProduceNumber("longitude"))
		latitude_string = fmt.Sprintf("%.6f", randomProduceNumber("latitude"))
		temperature_string = fmt.Sprintf("%.1f", randomProduceNumber("temperature"))
		humidity_string = fmt.Sprintf("%.1f", randomProduceNumber("humidity"))

		outStr := key_string + "," + longitude_string + "," + latitude_string + "," + temperature_string + "," + humidity_string + "\n"

		outputWriter.WriteString(outStr)

		//fmt.Println(longitude_string)
		//fmt.Println(latitude_string)
		//fmt.Println(temperature_string)
		//fmt.Println(humidity_string)

		i = i + 1
	}

	outputWriter.Flush()
}
 