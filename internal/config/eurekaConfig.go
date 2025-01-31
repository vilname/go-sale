package config

import (
	"fmt"
	"github.com/ArthurHlt/go-eureka-client/eureka"
	"os"
	"strconv"
)

func EurekaConfig() {
	eurekaName := os.Getenv("EUREKA_NAME")
	serverIp := os.Getenv("SERVER_IP")
	domain := os.Getenv("DOMAIN_NAME")
	port, _ := strconv.ParseInt(os.Getenv("PORT"), 10, 64)

	client := eureka.NewClient([]string{
		fmt.Sprintf("https://%s:8761/eureka", domain), //From a spring boot based eureka server
		// add others servers here
	})
	instance := eureka.NewInstanceInfo(eurekaName, "telemetry-sale", serverIp, int(port), 30, false) //Create a new instance to register
	instance.Metadata = &eureka.MetaData{
		Map: make(map[string]string),
	}
	//instance.Metadata.Map["foo"] = "bar"                          //add metadata for example
	client.RegisterInstance("telemetry-sale", instance)   // Retrieves all applications from eureka server(s)
	client.GetApplication(instance.App)                   // retrieve the application "test"
	client.GetInstance(instance.App, instance.HostName)   // retrieve the instance from "test.com" inside "test"" app
	client.SendHeartbeat(instance.App, instance.HostName) // say to eureka that your app is alive (here you must send heartbeat before 30 sec)
}
