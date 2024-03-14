package main

import (
	"fmt"
	"hlsapi/src/Configuration"
	"hlsapi/src/FileEndpoints"
	"net/http"
)

func main() {
	Configuration.Init(fmt.Sprintf("%s", "appSettings.json"))
	configuration := Configuration.ReadConfiguration()

	PORT := fmt.Sprintf(":%d", configuration.Port)

	http.HandleFunc("/store", FileEndpoints.StoreFileOnDisk)
	http.HandleFunc("/get", FileEndpoints.GetFileFromDisk)

	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
