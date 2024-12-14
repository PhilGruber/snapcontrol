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
			cl, err := client.ClientGetStatus(clientId)
			if err != nil {
				fmt.Println(err)
				return
			}
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
	case "group":
		switch command {
		case "status":
			/* TODO: Implement */
		case "mute":
			/* TODO: Implement */
		case "clients":
			/* TODO: Implement */
		case "name":
			/* TODO: Implement */
		}
	case "server":
		switch command {
		case "status":
			svr := client.ServerGetStatus()
			for _, group := range svr.Groups {
				name := group.Name
				if name == "" {
					name = group.Id
				}
				fmt.Printf("Group %s\n", name)
				for _, client := range group.Clients {
					fmt.Printf("\tClient %s: %s (%s)\n", client.Id, client.Config.Name, client.Host.Name)
				}
			}
		case "version":
			version := client.ServerGetRPCVersion()
			fmt.Println("Server version: ", version)
		case "deleteclient":
			err := client.ServerDeleteClient(os.Args[3])
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Client deleted")
			}
		}
	case "stream":
		switch command {
		case "addstream":
			/* TODO: Implement */
		case "removestream":
			/* TODO: Implement */
		case "control":
			/* TODO: Implement */
		}
	}
}
