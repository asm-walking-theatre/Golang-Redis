package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/Jeffail/gabs"
	"github.com/gomodule/redigo/redis"
)

type Address struct {
	Latitude  float32 `json:"customer_latitude"`
	Longitude float32 `json:"customer_longitude"`
	Pincode   int     `json:"customer_pincode"`
}

type loc struct {
	Lat float32 `json: lat`
	Lon float32 `json: lon`
}

var Request_Count int = 1
var PincodeVersionMap = make(map[int]int)

func ResponseFromRedis(lat, long *float32) int {

	conn, err := redis.Dial("tcp", "localhost:6379")
	checkError(err, "Error in connecting to redis")
	defer conn.Close()
	result := "-1"
	reply, err := conn.Do("GEORADIUS", "maps", *long, *lat, 1, "km", "ASC")
	checkError(err, "Error in GEORADIUS response")

	data := reply.([]interface{})

	for key, value := range data {
		data[key] = string(value.([]byte))
		if key == 0 {
			result = data[key].(string)
			result = result[0:6]
		}
	}

	result_int, err := strconv.Atoi(result)
	checkError(err, "Error in string conversion to int")

	return result_int
}

func AddToRedis(lat, long *float32, pincode int) error {
	conn, err := redis.Dial("tcp", "localhost:6379")
	checkError(err, "Error in connecting to Redis")
	defer conn.Close()

	var pincode_v int = pincode*100 + PincodeVersionMap[pincode]
	PincodeVersionMap[pincode]++

	reply, err := conn.Do("GEOADD", "maps", *long, *lat, pincode_v)

	var CheckAdded int64 = 1

	if reply == CheckAdded {
		fmt.Printf("Successfully added lat:%v, long:%v for pincode: %v \n", *lat, *long, pincode)
	} else {
		fmt.Printf("new type of error")
	}

	return err
}

func checkError(err error, response string) {
	if err != nil {
		panic(err)
		fmt.Println(response)
	}
}

func RequestToGoogleMapsAPI(lat, long *float32) int {

	var myurl string = "https://nominatim.openstreetmap.org/reverse?format=json&lat=" + fmt.Sprintf("%v", *lat) + "&lon=" + fmt.Sprintf("%v", *long)
	response, err := http.Get(myurl)
	checkError(err, "Error in GET request from openstreet maps API")
	defer response.Body.Close()

	contentFromAPI, err := ioutil.ReadAll(response.Body)
	checkError(err, "Error in reading body of GET request")

	jsonParsed, err := gabs.ParseJSON(contentFromAPI)
	checkError(err, "Error in Parsing JSON to get pincode")

	pincode := jsonParsed.Path("address.postcode").String()
	pincode = pincode[1:7]

	pincode_int, err := strconv.Atoi(pincode)
	checkError(err, "Error in string to int conversion")

	return pincode_int
}

func CustomerHandler(w http.ResponseWriter, r *http.Request) {

	location := loc{}

	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error reading the body", err)
	}

	err = json.Unmarshal(jsn, &location)
	if err != nil {
		log.Fatal("Decoding error: ", err)
	}

	log.Printf("Received: %v\n", location)

	var customerAddress Address

	customerAddress.Latitude = location.Lat
	customerAddress.Longitude = location.Lon

	fmt.Println("Request Count : ", Request_Count)

	if Request_Count == 1 {

		FlushallInRedis()
		customerAddress.Pincode = RequestToGoogleMapsAPI(&location.Lat, &location.Lon)
		AddToRedis(&location.Lat, &location.Lon, customerAddress.Pincode)
		Request_Count++

	} else {
		response := ResponseFromRedis(&location.Lat, &location.Lon)
		Request_Count++

		if response == -1 {
			customerAddress.Pincode = RequestToGoogleMapsAPI(&location.Lat, &location.Lon)
			AddToRedis(&location.Lat, &location.Lon, customerAddress.Pincode)

		} else {
			customerAddress.Pincode = response
		}

	}

	customerAddressJson, err := json.Marshal(customerAddress)
	if err != nil {
		fmt.Fprintf(w, "Error: converting to JSOn %s", err)
	}
	fmt.Println("Success")
	w.Header().Set("Content-Type", "application/json")
	w.Write(customerAddressJson)

}

func FlushallInRedis() {
	conn, err := redis.Dial("tcp", "localhost:6379")
	checkError(err, "Error in flushall at redis")
	defer conn.Close()
	var i int = 5
	reply, err := conn.Do("FLUSHALL")
	if i == 3 {
		fmt.Println("Cant leave reply from redis", reply)
	}
	checkError(err, "Error in flushall at redis")
}

func main() {
	http.HandleFunc("/", CustomerHandler)
	http.ListenAndServe(":8080", nil)

}
