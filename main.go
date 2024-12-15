package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	if len(os.Args) > 1 && os.Args[1] == "help" {
		printHelp()
		return
	}

	if len(os.Args) < 3 {
		fmt.Println("Usage: snapcontrol <client|group|server|stream> <command> [args]")
		return
	}

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
			printOrError(fmt.Sprintf("Client %s: %s. Volume: %d%%\n", cl.Id, cl.Config.Name, cl.Config.Volume.Percent), err)
		case "volume":
			volume, _ := strconv.Atoi(os.Args[4])
			if volume < 0 || volume > 100 {
				fmt.Println("Volume must be between 0 and 100")
			}
			err := client.ClientSetVolume(clientId, volume)
			printOrError(fmt.Sprintf("Volume set to %d%%", volume), err)
		case "name":
			name := os.Args[4]
			err := client.ClientSetName(clientId, name)
			printOrError(fmt.Sprintf("Name set to %s\n", name), err)
		case "latency":
			latency, _ := strconv.Atoi(os.Args[4])
			err := client.SetClientLatency(clientId, latency)
			printOrError(fmt.Sprintf("Latency set to %d\n", latency), err)
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
				fmt.Println()
				for _, stream := range svr.Streams {
					fmt.Printf("Stream %s: %s (%s)\n", stream.Id, stream.Status, stream.Uri.Scheme)
				}
			}
		case "version":
			version := client.ServerGetRPCVersion()
			fmt.Println("Server version: ", version)
		case "deleteclient":
			err := client.ServerDeleteClient(os.Args[3])
			printOrError("Client deleted", err)
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

func printOrError(msg string, err error) {
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(msg)
	}
}

func printHelp() {
	fmt.Println("Usage: snapcontrol <client|group|server|stream> <command> [args]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("client status <id>")
	fmt.Println("\tShow status of a specific client")
	fmt.Println()
	fmt.Println("client volume <id> <volume>")
	fmt.Println("\tSet volume of a specific client")
	fmt.Println()
	fmt.Println("client name <id> <name>")
	fmt.Println("\tChange name of a specific client")
	fmt.Println()
	fmt.Println("client latency <id> <latency>")
	fmt.Println("\tSet latency of a specific client")
	fmt.Println()
	fmt.Println("group status <id>")
	fmt.Println("group mute <id> <true|false>")
	fmt.Println("group clients <id>")
	fmt.Println("group name <id> <name>")
	fmt.Println()
	fmt.Println("server status")
	fmt.Println("\tShow all groups and clients")
	fmt.Println("server version")
	fmt.Println("\tShow RPC version of server")
	fmt.Println("server deleteclient <id>")
	fmt.Println()
	fmt.Println("stream addstream <id>")
	fmt.Println("stream removestream <id>")
	fmt.Println("stream control <id> <play|pause|stop>")
}
