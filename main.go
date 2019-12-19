package main

import (
	"flag"
	"fmt"
	"os"
	"porter/api"
	"porter/config"
	"porter/hw"
	"time"
)

func main() {
	configFile := flag.String("config", "config.json", "configuration file path")
	flag.Parse()

	appCfg, err := config.Load(*configFile)
	if err != nil {
		fmt.Printf("Error reading configuration file: %s", *configFile, err.Error())
		os.Exit(1)
	}

	if err = hw.Init(); err != nil {
		fmt.Printf("Error initializing GPIO: %s", err.Error())
		os.Exit(1)
	}
	defer hw.Close()

	apiServer := api.NewServer(appCfg)
	go panic(apiServer.Serve())

	for {
		time.Sleep(60 * time.Second)
	}
}
