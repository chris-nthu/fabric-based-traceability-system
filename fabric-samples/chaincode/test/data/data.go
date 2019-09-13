package main

import (
	"fmt"
	"math/rand"
	"time"
	"strconv"
	_"reflect"
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
	/*
	products := []Product {
		Product{GPS_Location: "(41.40338, 2.17403)", Temperature: "26.3", Humidity: "70.1"},
		Product{GPS_Location: "(53.22449, 3.79878)", Temperature: "27.8", Humidity: "32.3"},
		Product{GPS_Location: "(49.12345, 2.66534)", Temperature: "28.5", Humidity: "50.0"},
	}

	fmt.Print(products[0])
	fmt.Println()
	fmt.Print(products[1])
	fmt.Println()
	fmt.Print(products[2])
	fmt.Println()
	*/


	var products [10]Product
	var longitude_string string
	var latitude_string string
	var location_string string
	var temperature_string string
	var humidity_string string

	i := 0
	for i < len(products) {
		longitude_string = fmt.Sprintf("%.6f", randomProduceNumber("longitude"))
		latitude_string = fmt.Sprintf("%.6f", randomProduceNumber("latitude"))
		temperature_string = fmt.Sprintf("%.1f", randomProduceNumber("temperature"))
		humidity_string = fmt.Sprintf("%.1f", randomProduceNumber("humidity"))
		location_string = "(" + longitude_string + ", " + latitude_string + ")"

		products[i] = Product{GPS_Location: location_string, Temperature: temperature_string, Humidity: humidity_string}
		//fmt.Printf("%d\n", i)

		i = i + 1
	}

	i = 0
	for i <len(products) {
		fmt.Println(products[i])

		i = i + 1
	}

	fmt.Println()
	fmt.Println()

	products2 := []Product {
		Product{GPS_Location: "(41.40338, 2.17403)", Temperature: "26.3", Humidity: "70.1"},
		Product{GPS_Location: "(53.22449, 3.79878)", Temperature: "27.8", Humidity: "32.3"},
		Product{GPS_Location: "(49.12345, 2.66534)", Temperature: "28.5", Humidity: "50.0"},
	}

	i = 0
	for i <len(products2) {
		fmt.Println(products2[i])

		i = i + 1
	}
}
 