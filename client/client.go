package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	//"github.com/gomodule/redigo/redis"
	//"strconv"
)

type Address struct {
	Latitude  float32 `json:"customer_latitude"`
	Longitude float32 `json:"customer_longitude"`
	Pincode   int     `json:"customer_pincode"`
}

type loc struct {
	Lat float32 `json: "lat"`
	Lon float32 `json: "lon"`
}

func client(lat, lon float32) {

	location_send := loc{Lat: lat, Lon: lon}
	locJson, err := json.Marshal(location_send)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8080", bytes.NewBuffer(locJson))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response: ", string(body))
	resp.Body.Close()
}

func main() {
	client(22.34601, 87.23197)
	client(19.7174, 77.1494)
	client(21.146633, 79.088860)
	client(19.1383, 77.3210)
	client(19.07609, 72.877426)
	client(19.717401, 77.149401)
	client(22.34601, 87.23197)
	client(19.7174, 77.1494)
	client(21.146633, 79.088860)
	client(19.1383, 77.3210)
	client(19.07609, 72.877426)
	client(19.717401, 77.149401)
}
