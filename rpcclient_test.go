package main

import (
	"bufio"
	"encoding/json"
	"net"
	"reflect"
	"strings"
	"testing"
)

func TestAdditionalAPICalls(t *testing.T) {
	tests := []struct {
		name   string
		call   func(*rpcClient) error
		method string
		params map[string]any
		result string
	}{
		{"set group stream", func(c *rpcClient) error { return c.SetGroupStream("group-1", "stream-1") }, "Group.SetStream", map[string]any{"id": "group-1", "stream_id": "stream-1"}, `{}`},
		{"set group clients", func(c *rpcClient) error {
			c.clientIds["kitchen"] = "client-1"
			_, err := c.SetGroupClients("group-1", []string{"kitchen"})
			return err
		}, "Group.SetClients", map[string]any{"id": "group-1", "clients": []any{"client-1"}}, `{"server":{"groups":[]}}`},
		{"control stream", func(c *rpcClient) error { return c.StreamControl("stream-1", "seek", map[string]any{"offset": 60}) }, "Stream.Control", map[string]any{"id": "stream-1", "command": "seek", "params": map[string]any{"offset": float64(60)}}, `"ok"`},
		{"set stream property", func(c *rpcClient) error { return c.StreamSetProperty("stream-1", "shuffle", true) }, "Stream.SetProperty", map[string]any{"id": "stream-1", "property": "shuffle", "value": true}, `"ok"`},
		{"add stream", func(c *rpcClient) error { _, err := c.StreamAdd("pipe:///tmp/snapfifo"); return err }, "Stream.AddStream", map[string]any{"streamUri": "pipe:///tmp/snapfifo"}, `{"stream_id":"stream-1"}`},
		{"remove stream", func(c *rpcClient) error { _, err := c.StreamRemove("stream-1"); return err }, "Stream.RemoveStream", map[string]any{"id": "stream-1"}, `{"stream_id":"stream-1"}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			defer listener.Close()
			received := make(chan request, 1)
			go serveOnce(listener, test.result, received)

			client := newRpcClient("127.0.0.1", listener.Addr().(*net.TCPAddr).Port, false)
			if err := test.call(client); err != nil {
				t.Fatal(err)
			}
			request := <-received
			if request.Method != test.method {
				t.Fatalf("method = %q, want %q", request.Method, test.method)
			}
			params, ok := request.Params.(map[string]any)
			if !ok || !reflect.DeepEqual(params, test.params) {
				t.Fatalf("params = %#v, want %#v", request.Params, test.params)
			}
		})
	}
}

func serveOnce(listener net.Listener, result string, received chan<- request) {
	conn, err := listener.Accept()
	if err != nil {
		return
	}
	defer conn.Close()
	var request request
	if json.NewDecoder(bufio.NewReader(conn)).Decode(&request) != nil {
		return
	}
	received <- request
	_, _ = conn.Write([]byte(`{"id":1,"jsonrpc":"2.0","result":` + result + "}\n"))
}

func TestInvalidInputs(t *testing.T) {
	if _, err := parseParameters([]string{"not-a-pair"}); err == nil {
		t.Fatal("parseParameters accepted a parameter without key=value")
	}
	if _, err := parseRequiredJSONValue("not-json"); err == nil {
		t.Fatal("parseRequiredJSONValue accepted invalid JSON")
	}
	value, err := parseRequiredJSONValue(`{"volume":55}`)
	if err != nil || !reflect.DeepEqual(value, map[string]any{"volume": float64(55)}) {
		t.Fatalf("valid JSON value = %#v, %v", value, err)
	}
}

func TestConnectionErrorIsReturned(t *testing.T) {
	client := newRpcClient("127.0.0.1", -1, false)
	_, err := client.sendRequest(request{Id: 1, Jsonrpc: version, Method: "Server.GetStatus"})
	if err == nil {
		t.Fatal("sendRequest succeeded with an invalid address")
	}
	if !strings.Contains(err.Error(), "connect to Snapcast") {
		t.Fatalf("error = %q, want connection context", err)
	}
}
