package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
)

const (
	version = "2.0"
)

type rpcClient struct {
	url  string
	port int
}

func newRpcClient(url string, port int) *rpcClient {
	return &rpcClient{
		url:  url,
		port: port,
	}
}

func (c *rpcClient) ClientGetStatus(id string) *client {
	request := request{
		Id:      1,
		Jsonrpc: version,
		Method:  "Client.GetStatus",
		Params: clientId{
			Id: id,
		},
	}
	response, err := c.sendRequest(request)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Printf("Client %s status: %s\n", id, response.Result.Client.Config.Name)
	return response.Result.Client
}

func (c *rpcClient) ClientSetVolume(id string, vol int) {
	muted := vol == 0
	request := request{
		Id:      1,
		Jsonrpc: version,
		Method:  "Client.SetVolume",
		Params: volumeRequest{
			Id: id,
			Volume: volume{
				Muted:   muted,
				Percent: vol,
			},
		},
	}
	_, err := c.sendRequest(request)
	if err != nil {
		log.Println(err)
	}
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
		log.Println(err)
		return nil
	}

	return response.Result.Server
}

func (c *rpcClient) sendRequest(request request) (*response, error) {
	log.Printf("Connecting to %s:%d\n", c.url, c.port)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.url, c.port))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	data, _ := json.Marshal(request)
	data = append(data, '\n')
	log.Printf("Sending request: %s\n", string(data))
	_, err = conn.Write(data)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 10240)
	length, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}

	buf = buf[:length]

	fmt.Printf("Response: %s\n", string(buf))

	var response response
	err = json.Unmarshal(buf, &response)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	fmt.Printf("Error: %v\n", response.Error)

	if response.Error != nil {
		return nil, errors.New(response.Error.Message + ": " + response.Error.Data.(string))
	}
	log.Printf("Result: %s\n", response.Result)
	return &response, nil

}
