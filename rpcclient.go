package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	version         = "2.0"
	maxResponseSize = 4 << 20
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

	if response.Result.Client == nil {
		return nil, errors.New("invalid response: missing client status")
	}
	c.log(fmt.Sprintf("Client %s status: %s\n", id, response.Result.Client.Config.Name))
	return response.Result.Client, nil
}

func (c *rpcClient) ClientSetVolume(id string, vol int) error {
	client, err := c.ClientGetStatus(id)
	if err != nil {
		return err
	}
	request := request{
		Id:      1,
		Jsonrpc: version,
		Method:  "Client.SetVolume",
		Params: volumeRequest{
			Id: client.Id,
			Volume: volume{
				Muted:   client.Config.Volume.Muted,
				Percent: vol,
			},
		},
	}
	response, err := c.sendRequest(request)
	if err != nil {
		return err
	}
	return requireObjectResult(response.Result, "Client.SetVolume")
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
	response, err := c.sendRequest(request)
	if err != nil {
		return err
	}
	return requireObjectResult(response.Result, "Client.SetName")
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
	response, err := c.sendRequest(request)
	if err != nil {
		return err
	}
	return requireObjectResult(response.Result, "Client.SetLatency")
}

func (c *rpcClient) ServerGetStatus() (*server, error) {
	request := request{
		Id:      1,
		Jsonrpc: version,
		Method:  "Server.GetStatus",
		Params:  nil,
	}
	response, err := c.sendRequest(request)
	if err != nil {
		return nil, err
	}
	if response.Result.Server == nil {
		return nil, errors.New("invalid response: missing server status")
	}
	return response.Result.Server, nil
}

func (c *rpcClient) ServerGetRPCVersion() (string, error) {
	request := request{
		Id:      1,
		Jsonrpc: version,
		Method:  "Server.GetRPCVersion",
	}
	response, err := c.sendRequest(request)
	if err != nil {
		return "", err
	}

	if err := requireObjectResult(response.Result, "Server.GetRPCVersion"); err != nil {
		return "", err
	}
	return fmt.Sprintf("%d.%d.%d", response.Result.Major, response.Result.Minor, response.Result.Patch), nil
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
	response, err := c.sendRequest(request)
	if err != nil {
		return err
	}
	return requireObjectResult(response.Result, "Server.DeleteClient")
}

func (c *rpcClient) GroupGetStatus(id string) (*group, error) {
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
		return nil, err
	}
	if response.Result.Group == nil {
		return nil, errors.New("invalid response: missing group status")
	}
	return response.Result.Group, nil
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
	response, err := c.sendRequest(request)
	if err != nil {
		return err
	}
	return requireObjectResult(response.Result, "Group.SetMute")
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
	response, err := c.sendRequest(request)
	if err != nil {
		return err
	}
	return requireObjectResult(response.Result, "Group.SetName")
}

func (c *rpcClient) SetGroupStream(id string, streamID string) error {
	response, err := c.sendRequest(request{
		Id:      1,
		Jsonrpc: version,
		Method:  "Group.SetStream",
		Params: groupStreamRequest{
			Id:       id,
			StreamID: streamID,
		},
	})
	if err != nil {
		return err
	}
	return requireObjectResult(response.Result, "Group.SetStream")
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
	if response.Result.Server == nil {
		return nil, errors.New("invalid Group.SetClients response: missing server status")
	}
	return response.Result.Server, nil
}

func (c *rpcClient) StreamControl(id string, command string, params map[string]any) error {
	response, err := c.sendRequest(request{
		Id:      1,
		Jsonrpc: version,
		Method:  "Stream.Control",
		Params:  streamControlRequest{Id: id, Command: command, Params: params},
	})
	if err != nil {
		return err
	}
	if !response.Result.isOK() {
		return errors.New("invalid Stream.Control response: expected result \"ok\"")
	}
	return nil
}

func (c *rpcClient) StreamSetProperty(id string, property string, value any) error {
	response, err := c.sendRequest(request{
		Id:      1,
		Jsonrpc: version,
		Method:  "Stream.SetProperty",
		Params:  streamPropertyRequest{Id: id, Property: property, Value: value},
	})
	if err != nil {
		return err
	}
	if !response.Result.isOK() {
		return errors.New("invalid Stream.SetProperty response: expected result \"ok\"")
	}
	return nil
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
	if response.Result.StreamID == "" {
		return "", errors.New("invalid Stream.AddStream response: missing stream_id")
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
	if response.Result.StreamID == "" {
		return "", errors.New("invalid Stream.RemoveStream response: missing stream_id")
	}
	return response.Result.StreamID, nil
}

func (c *rpcClient) sendRequest(request request) (*response, error) {
	c.log(fmt.Sprintf("Connecting to %s:%d\n", c.url, c.port))
	address := fmt.Sprintf("%s:%d", c.url, c.port)
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("connect to Snapcast at %s: %w", address, err)
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return nil, fmt.Errorf("set request deadline: %w", err)
	}

	data, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}
	data = append(data, '\n')
	c.log(fmt.Sprintf("Sending request: %s\n", string(data)))
	_, err = conn.Write(data)
	if err != nil {
		return nil, fmt.Errorf("send request to Snapcast: %w", err)
	}

	var response response
	decoder := json.NewDecoder(io.LimitReader(conn, maxResponseSize))
	err = decoder.Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("decode Snapcast response: %w", err)
	}
	if response.Jsonrpc != version {
		return nil, fmt.Errorf("invalid Snapcast response: jsonrpc = %q", response.Jsonrpc)
	}
	if response.Id != request.Id {
		return nil, fmt.Errorf("invalid Snapcast response: id = %d, want %d", response.Id, request.Id)
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

func requireObjectResult(result result, method string) error {
	if !result.isObject() {
		return fmt.Errorf("invalid %s response: expected an object result", method)
	}
	return nil
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

	svr, err := c.ServerGetStatus()
	if err == nil && svr.Groups != nil {
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
