package main

import (
	"flag"
	"fmt"
	"os"
	"porter/api"
	"porter/config"
	"porter/hw"
)

func main() {
	configFile := flag.String("config", "config.json", "configuration file path")
	flag.Parse()

	appCfg, err := config.Load(*configFile)
	if err != nil {
		fmt.Printf("Error reading configuration file: %s", err.Error())
		os.Exit(1)
	}

	if err = hw.Init(); err != nil {
		fmt.Printf("Error initializing GPIO: %s", err.Error())
		os.Exit(1)
	}
	defer func() {
		err := hw.Close()
		if err != nil {
			panic(err)
		}
	}()

	apiServer := api.NewServer(appCfg)

	panic(apiServer.Serve())
}
