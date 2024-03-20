package main

import (
	"fmt"
	"hlsapi/Daemons/StorageDaemon/src"
	"hlsapi/src/Domain/AppConfiguration"
	httpHandlers "hlsapi/src/IOChannel/http"
	"net/http"
)

func main() {
	AppConfiguration.Initialize("appSettings.json")

	go StorageDaemon.Start()

	PORT := fmt.Sprintf(":%d", AppConfiguration.ReadRoot().Server.Port)

	http.HandleFunc("/store", httpHandlers.StoreFileOnDisk)
	http.HandleFunc("/get", httpHandlers.GetFileFromDisk)

	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
