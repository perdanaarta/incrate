package main

import (
	"fmt"
	"incrate/services/api"
)

func main() {
	server := api.NewAPIsServer("127.0.0.1", 8080)

	if err := server.Run(); err != nil {
		fmt.Printf("Error occured: %s", err.Error())
	}
}
