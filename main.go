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
		fmt.Printf("Error reading configuration file: %s\n", err.Error())
		os.Exit(1)
	}

	if err = hw.Init(); err != nil {
		fmt.Printf("Error initializing GPIO: %s\n", err.Error())
		os.Exit(1)
	}
	defer func() {
		err := hw.Close()
		if err != nil {
			fmt.Printf("Error un-initializing GPIO: %s\n", err.Error())
			os.Exit(1)
		}
	}()

	apiServer := api.NewServer(appCfg)
	if err = apiServer.Serve(); err != nil {
		fmt.Printf("API server died with error: %s\n", err.Error())
		os.Exit(1)
	}

}
