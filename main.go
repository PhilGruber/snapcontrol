package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func main() {

	if len(os.Args) > 1 && os.Args[1] == "help" {
		printHelp()
		return
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: snapcontrol <client|group|server> <command> [args]")
		fmt.Println("Use 'snapcontrol help' for more information")
		return
	}

	subsystem := strings.ToLower(os.Args[1])

	hostname := "127.0.0.1"
	if os.Getenv("SNAPCONTROL_HOSTNAME") != "" {
		hostname = os.Getenv("SNAPCONTROL_HOSTNAME")
	}

	client := newRpcClient(hostname, 1705, false)

	switch subsystem {
	case "client":
		if len(os.Args) < 4 {
			fmt.Println("Usage: snapcontrol client <status|volume|name|latency> <clientId> [<value>]")
			return
		}
		command := strings.ToLower(os.Args[2])
		clientId := os.Args[3]
		switch command {
		case "status":
			cl, err := client.ClientGetStatus(clientId)
			if err != nil {
				printOrError("", err)
				return
			}
			fmt.Printf("Client %s: %s. Volume: %d%%\n", cl.Id, cl.Config.Name, cl.Config.Volume.Percent)
		case "volume":
			if len(os.Args) < 5 {
				fmt.Println("Usage: snapcontrol client volume <clientId> <volume>")
				return
			}
			volumeStr := os.Args[4]
			if volumeStr == "" {
				fmt.Println("Invalid volume value: empty value")
				return
			}
			relative := volumeStr[0] == '+' || volumeStr[0] == '-'
			volume, err := strconv.Atoi(volumeStr)
			if err != nil {
				fmt.Println("Invalid volume value:", err)
				return
			}
			if !relative && (volume < 0 || volume > 100) {
				fmt.Println("Volume must be between 0 and 100")
				return
			}
			if relative {
				cl, err := client.ClientGetStatus(clientId)
				if err != nil {
					fmt.Println("Error getting client status:", err)
					return
				}
				volume += cl.Config.Volume.Percent
				if volume < 0 {
					volume = 0
				} else if volume > 100 {
					volume = 100
				}
			}
			err = client.ClientSetVolume(clientId, volume)
			printOrError(fmt.Sprintf("Volume set to %d%%", volume), err)
		case "name":
			if len(os.Args) < 5 {
				fmt.Println("Usage: snapcontrol client name <clientId> <name>")
				return
			}
			name := os.Args[4]
			err := client.ClientSetName(clientId, name)
			printOrError(fmt.Sprintf("Name set to %s\n", name), err)
		case "latency":
			if len(os.Args) < 5 {
				fmt.Println("Usage: snapcontrol client latency <clientId> <latency>")
				return
			}
			latency, err := strconv.Atoi(os.Args[4])
			if err != nil {
				fmt.Println("Invalid latency value:", err)
				return
			}
			err = client.SetClientLatency(clientId, latency)
			printOrError(fmt.Sprintf("Latency set to %d\n", latency), err)
		default:
			fmt.Println("Usage: snapcontrol client <status|volume|name|latency> <clientId> [<value>]")
		}
	case "group":
		if len(os.Args) < 4 {
			fmt.Println("Usage: snapcontrol group <status|mute|stream|clients|name> <groupId>")
			return
		}
		command := strings.ToLower(os.Args[2])
		groupId := os.Args[3]
		switch command {
		case "status":
			group, err := client.GroupGetStatus(groupId)
			if err != nil {
				printOrError("", err)
				return
			}
			fmt.Printf("\t%-36s %-16s %-16s %-9s %-12s\n", "Id", "Name", "Host", "Volume", "Latency")
			for _, client := range group.Clients {
				fmt.Printf("\t%-36s %-16s %-16s %5d%% %8dms\n", client.Id, client.Config.Name, client.Host.Name, client.Config.Volume.Percent, client.Config.Latency)
			}
		case "mute":
			if len(os.Args) < 5 {
				fmt.Println("Usage: snapcontrol group mute <groupId> <true|false>")
				return
			}
			mute, err := strconv.ParseBool(os.Args[4])
			if err != nil {
				fmt.Println("Mute must be true or false")
				return
			}
			err = client.SetGroupMute(groupId, mute)
			printOrError("Mute set", err)
		case "clients":
			if len(os.Args) < 5 {
				fmt.Println("Usage: snapcontrol group clients <groupId> <clientId|name>...")
				return
			}
			_, err := client.SetGroupClients(groupId, os.Args[4:])
			printOrError("Group clients set", err)
		case "stream":
			if len(os.Args) < 5 {
				fmt.Println("Usage: snapcontrol group stream <groupId> <streamId>")
				return
			}
			err := client.SetGroupStream(groupId, os.Args[4])
			printOrError("Group stream set", err)
		case "name":
			if len(os.Args) < 5 {
				fmt.Println("Usage: snapcontrol group name <groupId> <name>")
				return
			}
			err := client.SetGroupName(groupId, os.Args[4])
			printOrError("Name set", err)
		default:
			fmt.Println("Usage: snapcontrol group <status|mute|stream|clients|name> <groupId>")
		}
	case "server":
		if len(os.Args) < 3 {
			fmt.Println("Usage: snapcontrol server <status|version|deleteclient> [args]")
			return
		}
		command := strings.ToLower(os.Args[2])
		switch command {
		case "status":
			svr, err := client.ServerGetStatus()
			if err != nil {
				printOrError("", err)
				return
			}
			for _, group := range svr.Groups {
				name := group.Name
				if name == "" {
					name = group.Id
				}
				fmt.Printf("Group: %s\n", name)
				for _, client := range group.Clients {
					fmt.Printf("\tClient %-36s %-16s %-16s %5d%% %8dms\n", client.Id, client.Config.Name, client.Host.Name, client.Config.Volume.Percent, client.Config.Latency)
				}
				fmt.Println()
			}
			for _, stream := range svr.Streams {
				fmt.Printf("Stream %s: %s (%s)\n", stream.Id, stream.Status, stream.Uri.Scheme)
			}
		case "version":
			version, err := client.ServerGetRPCVersion()
			if err != nil {
				printOrError("", err)
				return
			}
			fmt.Println("Server version: ", version)
		case "deleteclient":
			if len(os.Args) < 4 {
				fmt.Println("Usage: snapcontrol server deleteclient <clientId>")
				return
			}
			err := client.ServerDeleteClient(os.Args[3])
			printOrError("Client deleted", err)
		default:
			fmt.Println("Usage: snapcontrol server <status|version|deleteclient> [args]")
		}
	case "stream":
		if len(os.Args) < 3 {
			fmt.Println("Usage: snapcontrol stream <add|remove|control|property> [args]")
			return
		}
		command := strings.ToLower(os.Args[2])
		switch command {
		case "add", "addstream":
			if len(os.Args) < 4 {
				fmt.Println("Usage: snapcontrol stream add <streamUri>")
				return
			}
			streamURI, err := url.Parse(os.Args[3])
			if err != nil || streamURI.Scheme == "" {
				fmt.Println("Invalid stream URI: a URI with a scheme is required")
				return
			}
			id, err := client.StreamAdd(os.Args[3])
			printOrError(fmt.Sprintf("Stream added: %s", id), err)
		case "remove", "removestream":
			if len(os.Args) < 4 {
				fmt.Println("Usage: snapcontrol stream remove <streamId>")
				return
			}
			id, err := client.StreamRemove(os.Args[3])
			printOrError(fmt.Sprintf("Stream removed: %s", id), err)
		case "control":
			if len(os.Args) < 5 {
				fmt.Println("Usage: snapcontrol stream control <streamId> <command> [key=value ...]")
				return
			}
			params, err := parseParameters(os.Args[5:])
			if err != nil {
				fmt.Println("Invalid control parameter:", err)
				return
			}
			err = client.StreamControl(os.Args[3], os.Args[4], params)
			printOrError("Stream command sent", err)
		case "property", "setproperty":
			if len(os.Args) < 6 {
				fmt.Println("Usage: snapcontrol stream property <streamId> <property> <jsonValue>")
				return
			}
			value, err := parseRequiredJSONValue(os.Args[5])
			if err != nil {
				fmt.Println("Invalid property value:", err)
				return
			}
			err = client.StreamSetProperty(os.Args[3], os.Args[4], value)
			printOrError("Stream property set", err)
		default:
			fmt.Println("Usage: snapcontrol stream <add|remove|control|property> [args]")
		}
	case "version":
		fmt.Println("snapcontrol version " + AppVersion)
	default:
		fmt.Println("Usage: snapcontrol <client|group|server|stream> <command> [args]")
		fmt.Println("Use 'snapcontrol help' for more information")
	}
}

func parseParameters(args []string) (map[string]any, error) {
	params := make(map[string]any, len(args))
	for _, arg := range args {
		key, value, found := strings.Cut(arg, "=")
		if !found || key == "" {
			return nil, fmt.Errorf("%q must use key=value form", arg)
		}
		params[key] = parseJSONValue(value)
	}
	return params, nil
}

func parseJSONValue(value string) any {
	var parsed any
	if json.Unmarshal([]byte(value), &parsed) == nil {
		return parsed
	}
	return value
}

func parseRequiredJSONValue(value string) (any, error) {
	var parsed any
	if err := json.Unmarshal([]byte(value), &parsed); err != nil {
		return nil, fmt.Errorf("must be valid JSON (quote string values): %w", err)
	}
	return parsed, nil
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
	fmt.Println("group stream <id> <streamId>")
	fmt.Println("\tAssign a stream to a group")
	fmt.Println("group clients <id> <clientId|name>...")
	fmt.Println("\tReplace the clients assigned to a group")
	fmt.Println("group name <id> <name>")
	fmt.Println()
	fmt.Println("server status")
	fmt.Println("\tShow all groups and clients")
	fmt.Println("server version")
	fmt.Println("\tShow RPC version of server")
	fmt.Println("server deleteclient <id>")
	fmt.Println()
	fmt.Println("stream add <streamUri>")
	fmt.Println("stream remove <streamId>")
	fmt.Println("stream control <streamId> <command> [key=value ...]")
	fmt.Println("stream property <streamId> <property> <jsonValue>")
}
