package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	subsystem := strings.ToLower(os.Args[1])
	command := strings.ToLower(os.Args[2])

	hostname := "127.0.0.1"
	if os.Getenv("SNAPCONTROL_HOSTNAME") != "" {
		hostname = os.Getenv("SNAPCONTROL_HOSTNAME")
	}

	client := newRpcClient(hostname, 1705)

	switch subsystem {
	case "client":
		clientId := os.Args[3]
		switch command {
		case "status":
			cl := client.ClientGetStatus(clientId)
			fmt.Printf("Client %s: %s. Volume: %d%%\n", cl.Id, cl.Config.Name, cl.Config.Volume.Percent)
		}
	case "server":
		switch command {
		case "status":
			svr := client.ServerGetStatus()
			for _, group := range svr.Groups {
				log.Printf("Group %s: %s\n", group.Id, group.Name)
				for _, client := range group.Clients {
					log.Printf("  Client %s: %s\n", client.Id, client.Config.Name)
				}
			}
		}
	}
}
