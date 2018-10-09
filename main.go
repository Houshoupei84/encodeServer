/*
Package main provides the encodeServer executable. See README.md for usage information.
*/
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ifIMust/encodeServer/server"
)

const (
	defaultPort = "8080"
)

func main() {
	port := defaultPort
	if len(os.Args) > 1 {
		port = os.Args[1]
		intPort, err := strconv.Atoi(port)
		if err != nil {
			fmt.Println("Invalid (non-numeric) port specified.")
			return
		}
		if intPort < 1024 || intPort > 49150 {
			fmt.Printf("Specified port (%d) out of supported range (1024-49150)\n", intPort)
			return
		}

	}
	server := server.NewServer(port)
	server.Run()
	<-server.ShutdownComplete
}
