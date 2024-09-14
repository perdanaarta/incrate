package main

import (
	"flag"
	"fmt"
	"incrate/config"
	"incrate/services/api"
	"incrate/services/log"
)

type Flag struct {
	FileCfg string
}

func loadLogWriter() {
	writer := log.NewLogWriter()
	{
		writer.AddConsoleWriter()
		// writer.AddFileWriter("app.log")
		writer.SetDefault()
	}
}

func loadFlags() *Flag {
	var cfg Flag

	// Define flags and bind them to the struct fields
	flag.StringVar(&cfg.FileCfg, "config", "", "Config file path")

	flag.Parse()

	return &cfg
}

func main() {
	loadLogWriter()
	fg := loadFlags()

	config.ConfigFile = fg.FileCfg
	conf := config.New()

	server := api.NewAPIsServer(conf.Server.Host, int(conf.Server.Port))

	if err := server.Run(); err != nil {
		fmt.Printf("Error occured: %s", err.Error())
	}
}
