package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

const (
	version = "2.0"
)

type rpcClient struct {
	url       string
	port      int
	verbose   bool
	clientIds map[string]string
}

func newRpcClient(url string, port int, verbose bool) *rpcClient {
	return &rpcClient{
		url:       url,
		port:      port,
		verbose:   verbose,
		clientIds: make(map[string]string),
	}
}

func (c *rpcClient) ClientGetStatus(id string) (*client, error) {
	request := request{
		Id:      1,
		Jsonrpc: version,
		Method:  "Client.GetStatus",
		Params: idOnly{
			Id: c.resolveClientId(id),
		},
	}
	response, err := c.sendRequest(request)
	if err != nil {
		return nil, err
	}

	c.log(fmt.Sprintf("Client %s status: %s\n", id, response.Result.Client.Config.Name))
	return response.Result.Client, nil
}

func (c *rpcClient) ClientSetVolume(id string, vol int) error {
	request := request{
		Id:      1,
		Jsonrpc: version,
		Method:  "Client.SetVolume",
		Params: volumeRequest{
			Id: c.resolveClientId(id),
			Volume: volume{
				Percent: vol,
			},
		},
	}
	_, err := c.sendRequest(request)
	return err
}

func (c *rpcClient) ClientSetName(id string, name string) error {
	request := request{
		Id:      1,
		Jsonrpc: version,
		Method:  "Client.SetName",
		Params: nameRequest{
			Id:   c.resolveClientId(id),
			Name: name,
		},
	}
	_, err := c.sendRequest(request)
	return err
}

func (c *rpcClient) SetClientLatency(id string, latency int) error {
	request := request{
		Id:      1,
		Jsonrpc: version,
		Method:  "Client.SetLatency",
		Params: latencyRequest{
			Id:      c.resolveClientId(id),
			Latency: latency,
		},
	}
	_, err := c.sendRequest(request)
	return err
}

func (c *rpcClient) ServerGetStatus() *server {
	request := request{
		Id:      1,
		Jsonrpc: version,
		Method:  "Server.GetStatus",
		Params:  nil,
	}
	response, err := c.sendRequest(request)
	if err != nil {
		fmt.Printf("%s: %v\n", err.Error(), response)
		return nil
	}

	return response.Result.Server
}

func (c *rpcClient) ServerGetRPCVersion() string {
	request := request{
		Id:      1,
		Jsonrpc: version,
		Method:  "Server.GetRPCVersion",
	}
	response, err := c.sendRequest(request)
	fmt.Printf("Response: %v\n", response)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return fmt.Sprintf("%d.%d.%d", response.Result.Major, response.Result.Minor, response.Result.Patch)
}

func (c *rpcClient) ServerDeleteClient(id string) error {
	request := request{
		Id:      1,
		Jsonrpc: version,
		Method:  "Server.DeleteClient",
		Params: idOnly{
			Id: c.resolveClientId(id),
		},
	}
	_, err := c.sendRequest(request)
	return err
}

func (c *rpcClient) GroupGetStatus(id string) *group {
	request := request{
		Id:      1,
		Jsonrpc: version,
		Method:  "Group.GetStatus",
		Params: idOnly{
			Id: id,
		},
	}
	response, err := c.sendRequest(request)
	if err != nil {
		return nil
	}

	return response.Result.Group
}

func (c *rpcClient) SetGroupMute(id string, mute bool) error {
	request := request{
		Id:      1,
		Jsonrpc: version,
		Method:  "Group.SetMute",
		Params: muteRequest{
			Id:   id,
			Mute: mute,
		},
	}
	_, err := c.sendRequest(request)
	return err
}

func (c *rpcClient) SetGroupName(id string, name string) error {
	request := request{
		Id:      1,
		Jsonrpc: version,
		Method:  "Group.SetName",
		Params: nameRequest{
			Id:   id,
			Name: name,
		},
	}
	_, err := c.sendRequest(request)
	return err
}

func (c *rpcClient) SetGroupStream(id string, streamID string) error {
	_, err := c.sendRequest(request{
		Id:      1,
		Jsonrpc: version,
		Method:  "Group.SetStream",
		Params: groupStreamRequest{
			Id:       id,
			StreamID: streamID,
		},
	})
	return err
}

func (c *rpcClient) SetGroupClients(id string, clientIDs []string) (*server, error) {
	resolvedIDs := make([]string, len(clientIDs))
	for i, clientID := range clientIDs {
		resolvedIDs[i] = c.resolveClientId(clientID)
	}
	response, err := c.sendRequest(request{
		Id:      1,
		Jsonrpc: version,
		Method:  "Group.SetClients",
		Params: groupClientsRequest{
			Id:      id,
			Clients: resolvedIDs,
		},
	})
	if err != nil {
		return nil, err
	}
	return response.Result.Server, nil
}

func (c *rpcClient) StreamControl(id string, command string, params map[string]any) error {
	_, err := c.sendRequest(request{
		Id:      1,
		Jsonrpc: version,
		Method:  "Stream.Control",
		Params:  streamControlRequest{Id: id, Command: command, Params: params},
	})
	return err
}

func (c *rpcClient) StreamSetProperty(id string, property string, value any) error {
	_, err := c.sendRequest(request{
		Id:      1,
		Jsonrpc: version,
		Method:  "Stream.SetProperty",
		Params:  streamPropertyRequest{Id: id, Property: property, Value: value},
	})
	return err
}

func (c *rpcClient) StreamAdd(streamURI string) (string, error) {
	response, err := c.sendRequest(request{
		Id:      1,
		Jsonrpc: version,
		Method:  "Stream.AddStream",
		Params:  streamURIRequest{StreamURI: streamURI},
	})
	if err != nil {
		return "", err
	}
	return response.Result.StreamID, nil
}

func (c *rpcClient) StreamRemove(id string) (string, error) {
	response, err := c.sendRequest(request{
		Id:      1,
		Jsonrpc: version,
		Method:  "Stream.RemoveStream",
		Params:  idOnly{Id: id},
	})
	if err != nil {
		return "", err
	}
	return response.Result.StreamID, nil
}

func (c *rpcClient) sendRequest(request request) (*response, error) {
	c.log(fmt.Sprintf("Connecting to %s:%d\n", c.url, c.port))
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.url, c.port))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	data, _ := json.Marshal(request)
	data = append(data, '\n')
	c.log(fmt.Sprintf("Sending request: %s\n", string(data)))
	_, err = conn.Write(data)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 10240)
	length, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	buf = buf[:length]

	c.log(fmt.Sprintf("Response: %s\n", string(buf)))

	var response response
	err = json.Unmarshal(buf, &response)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if response.Error != nil {
		if response.Error.Data == nil {
			return nil, errors.New(response.Error.Message)
		}
		return nil, fmt.Errorf("%s: %v", response.Error.Message, response.Error.Data)
	}

	c.log(fmt.Sprintf("Result: %+v\n", response.Result))
	return &response, nil
}

func (c *rpcClient) log(s string) {
	if c.verbose {
		fmt.Println(s)
	}
}

func (c *rpcClient) resolveClientId(name string) string {

	if c.clientIds[name] != "" {
		return c.clientIds[name]
	}

	svr := c.ServerGetStatus()
	if svr != nil && svr.Groups != nil {
		for _, grp := range svr.Groups {
			for _, client := range grp.Clients {
				if client.Config.Name == name || client.Host.Name == name || client.Id == name {
					c.clientIds[name] = client.Id
					return client.Id
				}
			}
		}
	}

	return name

}
