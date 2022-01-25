package main

import (
	Library "github.com/shd00700/library0"
	"github.com/stianeikeland/go-rpio"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/enitt-dev/go-utils/convert"
	"github.com/tarm/serial"
	"log"
	"os"
	"sync"
	"time"
	"fmt"
)

var Json []byte
var uri = "tcp://broker.hivemq.com:1883"
var topic1 = "safety_hsj/start"
var topic2 = "safety_hsj/sensor"
var topic3 = "safety_hsj/device"
var client = Library.Connect("pub", uri)

const (
	cool = rpio.Pin(6)
	heat = rpio.Pin(19)
	heat2 = rpio.Pin(13)
)

var sensorTypeMap = map[int]string{
	0x10: "co2",
	0x11: "tvoc",
	0x20: "humi",
	0x21: "humi2",
	0x30: "temp2",
	0x31: "temp",
}

type sensorData struct {
	SensorType string  `json:"sensor"`
	Value      float32 `json:"value"`
}

type Statedata struct {
	Heat string  `json:"heating_fan"`
	Cool string  `json:"cooling_fan"`
	Lamp string  `json:"heating_lamp"`
}

var MW = []string{"power"}

var data = []string{"data"}

func MQTTPublish(client mqtt.Client, topic string,json interface{}) {


		//time.Sleep(time.Second)
		Library.MQTTPublish(client, topic,json)
		println("data Publish")
}
func MQTTPublish2(client mqtt.Client, topic string,json interface{}) {


		//time.Sleep(time.Second)
		Library.MQTTPublish(client, topic,json)
		//println("data Publish")

}

func MQTTPublish3(client mqtt.Client, topic string,json interface{}) {


		//time.Sleep(time.Second)
		Library.MQTTPublish(client, topic,json)
		//println("data Publish")

}

func MQTTSubscribe(uri2 string, topic2 string) {

	for {
		Library.Listen(uri2, topic2, nil)
		//println("data Publish")
		time.Sleep(time.Second)
	}
}

//var humisensor = map[string]interface{}{}
//var GyroXsensor = map[string]interface{}{}
//var GyroYsensor = map[string]interface{}{}
//var GyroZsensor = map[string]interface{}{}
//var	result = []map[string]interface{}{}

func JsonMaker(SensorType string,value float32,leng []byte) []byte {

	sensor3 := map[string]interface{}{}
	result := []map[string]interface{}{}

	for i := 0; i < len(leng); i++ {
		sensor3[string(SensorType)] = value
	}

	result = append(result,sensor3 )
	jsonmaker, err := json.Marshal(result)
	//println("result:", result)

	if err != nil {
		panic(err)
	}
	return jsonmaker
}
func JsonMaker2() []byte {

	powerset := map[string]interface{}{}
	result := []map[string]interface{}{}
	powerset[string(MW[0])] = "on"
	result = append(result,powerset )
	jsonmaker2, err := json.Marshal(result)

	Library.MQTTPublish(client, topic1,jsonmaker2)

	if err != nil {
		panic(err)
	}
	return jsonmaker2
}

var SensorType string
var bb float32

func DataPaser(wg *sync.WaitGroup)interface{} {

	config := &serial.Config{
		Name: "/dev/ttyUSB2",
		Baud: 115200,
	}

	stream, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(stream)
	}

	err = stream.Flush()
	if err != nil {
		log.Fatal(err)
	}

	defer stream.Close()

	var (
		validatorQ = make(chan []byte)
	)

	go recv(validatorQ, stream)

	for {
		time.Sleep(time.Second*2)
		select {
		case recvBytes := <-validatorQ:

			log.Println()

			for i := 0; i < len(recvBytes); i = i + 10 {
				bytes := recvBytes[i : i+10]
				//sd := sensorData{
				//	SensorType: sensorTypeMap[int(bytes[4])],
				//	Value:      convert.BytesToFloat32(bytes[6:]),
				//}


				SensorType = sensorTypeMap[int(bytes[4])]
				bb = convert.BytesToFloat32(bytes[6:])
				Json = JsonMaker(SensorType, bb,recvBytes)

				go MQTTPublish(client, topic2,Json)
				time.Sleep(time.Second)

				//writeToFile(Json)
			}
		}
	}
	return Json
}

var state = map[string]int{}
func sys(wg *sync.WaitGroup){
	if err := rpio.Open(); err!= nil{
		fmt.Println(err)
		os.Exit(1)
	}
	cool.Output()
	heat.Output()
	heat2.Output()
	for {

		cool.High()
		heat.Low()
		heat.Low()
		time.Sleep(time.Second*12)
		cool.Low()
		heat.High()
		heat.High()
		time.Sleep(time.Second*10)

	}


	//for {
	//	if SensorType == "temp" {
	//		if bb >= 36 {
	//			cool.High()
	//			heat.Low()
	//			sd2 := Statedata{
	//				Heat: "off",
	//				Cool: "on",
	//				Lamp: "on",
	//			}
	//			ll, _ := json.Marshal(sd2)
	//			Json3 := string(append(ll, '\n'))
	//			MQTTPublish3(client, topic3, Json3)
	//
	//			println("비상\n\n")
	//			time.Sleep(time.Second)
	//
	//		} else if bb <= 35 && bb >= 25{
	//			cool.Low()
	//			heat.Low()
	//			sd2 := Statedata{
	//				Heat: "off",
	//				Cool: "off",
	//				Lamp: "on",
	//			}
	//			ll, _ := json.Marshal(sd2)
	//			Json3 := string(append(ll, '\n'))
	//			//Library.MQTTPublish(client, topic,Json3)
	//			MQTTPublish3(client, topic3, Json3)
	//
	//			println("대상대상\n\n")
	//			time.Sleep(time.Second)
	//
	//		} else {
	//			cool.Low()
	//			heat.High()
	//			sd2 := Statedata{
	//				Heat: "on",
	//				Cool: "off",
	//				Lamp: "on",
	//			}
	//			ll, _ := json.Marshal(sd2)
	//			Json3 := string(append(ll, '\n'))
	//			//Library.MQTTPublish(client, topic,Json3)
	//			MQTTPublish3(client, topic3, Json3)
	//
	//			println("엘즈\n\n")
	//			time.Sleep(time.Second)
	//
	//		}
			//	else if bb < 24{
			//		cool.Low()
			//		heat.High()
			//		sd2 := Statedata{
			//			Heat:   "on",
			//			Cool: "off",
			//			Lamp:   "on",
			//		}
			//		ll,_ := json.Marshal(sd2)
			//		Json3 := string(append(ll, '\n'))
			//		//Library.MQTTPublish(client, topic,Json3)
			//		MQTTPublish3(client, topic3,Json3)
			//
			//		println("cool OFF\n\n")
			//		//time.Sleep(time.Second)
			//	} else if bb == 25{
			//		cool.Low()
			//		heat.Low()
			//		sd2 := Statedata{
			//			Heat:   "off",
			//			Cool: "off",
			//			Lamp:   "on",
			//		}
			//		ll,_ := json.Marshal(sd2)
			//		Json3 := string(append(ll, '\n'))
			//		//Library.MQTTPublish(client, topic,Json3)
			//		MQTTPublish3(client, topic3,Json3)
			//
			//		println("cool OFF\n\n")
			//		//time.Sleep(time.Second)
			//	}
			//	time.Sleep(time.Second)
			//}
		}



func main() {
	JsonMaker2()
	var wg sync.WaitGroup

	wg.Add(1)
	go DataPaser(&wg)

	//wg.Add(2)
	//go sys(&wg)

	wg.Wait()
}


func writeToFile(msg string) {
	const fileName = "/home/pi/go/Test01/sensorData6.log"

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(msg); err != nil {
		panic(err)
	}

}

