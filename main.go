package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"os"
	"time"
)

type Config struct {
	Broker   string `json:"broker"`
	DeviceId string `json:"deviceId"`
}

func readConfig(filename string) Config {
	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var config Config
	err = json.Unmarshal(bytes, &config)

	if err != nil {
		panic(err)
	}

	return config
}

func main() {
	config := readConfig("config.json")
	fmt.Println(config)
	mqttClient := mqtt.NewClient(mqtt.NewClientOptions().AddBroker(config.Broker))

	mqttClient.Connect()

	fmt.Print("Connecting")
	for !mqttClient.IsConnected() {
		fmt.Print(".")
		time.Sleep(1 * time.Second)
	}
	fmt.Println("\nConnected")

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter command (e - exit, s - send message): ")
	for scanner.Scan() {
		switch scanner.Text() {

		case "e":
			mqttClient.Disconnect(0)
			return

		case "s":
			fmt.Print("Sensor uri: ")
			scanner.Scan()
			sensorUri := scanner.Text()

			fmt.Print("Value: ")
			scanner.Scan()
			value := scanner.Text()

			topic := fmt.Sprintf("xeo/%s/sensor/%s", config.DeviceId, sensorUri)
			mqttClient.Publish(topic, 1, false, value)

			fmt.Println("Message sent")
			break
		}

		fmt.Println("Enter command (e - exit, s - send message): ")
	}

}
