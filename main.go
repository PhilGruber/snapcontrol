package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	subsystem := strings.ToLower(os.Args[1])
	command := strings.ToLower(os.Args[2])

	hostname := "127.0.0.1"
	if os.Getenv("SNAPCONTROL_HOSTNAME") != "" {
		hostname = os.Getenv("SNAPCONTROL_HOSTNAME")
	}

	client := newRpcClient(hostname, 1705, false)

	switch subsystem {
	case "client":
		clientId := os.Args[3]
		switch command {
		case "status":
			cl := client.ClientGetStatus(clientId)
			fmt.Printf("Client %s: %s. Volume: %d%%\n", cl.Id, cl.Config.Name, cl.Config.Volume.Percent)
		case "volume":
			volume, _ := strconv.Atoi(os.Args[4])
			if volume < 0 || volume > 100 {
				fmt.Println("Volume must be between 0 and 100")
			}
			err := client.ClientSetVolume(clientId, volume)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Volume set to %d%%\n", volume)
			}
		case "name":
			name := os.Args[4]
			err := client.ClientSetName(clientId, name)
			if err == nil {
				fmt.Printf("Name set to %s\n", name)
			} else {
				fmt.Println(err)
			}
		case "latency":
			/* TODO: Implement */
		}
	case "server":
		switch command {
		case "status":
			svr := client.ServerGetStatus()
			for _, group := range svr.Groups {
				fmt.Printf("Group %s: %s\n", group.Id, group.Name)
				for _, client := range group.Clients {
					fmt.Printf("  Client %s: %s\n", client.Id, client.Config.Name)
				}
			}
		}
	}
}
