package main

import (
	"fmt"
	"hlsapi/src/Domain/AppConfiguration"
	httpHandlers "hlsapi/src/IOChannel/http"
	"hlsapi/src/Subdomain"
	"net/http"
)

func main() {
	err := AppConfiguration.JsonConfigurationProvider{}.Initialize("appSettings.json")
	if err != nil {
		panic(err)
	}

	go Subdomain.Start()

	PORT := fmt.Sprintf(":%d", AppConfiguration.JsonConfigurationProvider{}.GetConfiguration().Server.Port)

	http.HandleFunc("/convert", httpHandlers.CreateM3U8)
	http.HandleFunc("/wipe", httpHandlers.Wipe)
	http.HandleFunc("/get", httpHandlers.Get)
	http.HandleFunc("/health", httpHandlers.RespondToAHealthCheck)

	err = http.ListenAndServe(PORT, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
