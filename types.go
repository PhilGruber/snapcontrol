package main

import (
	"bytes"
	"encoding/json"
)

type request struct {
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
}

type response struct {
	Id      int       `json:"id"`
	Jsonrpc string    `json:"jsonrpc"`
	Result  result    `json:"result"`
	Error   *rpcError `json:"error"`
}

type idOnly struct {
	Id string `json:"id"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type result struct {
	raw      json.RawMessage
	Client   *client `json:"client"`
	Server   *server `json:"server"`
	Group    *group  `json:"group"`
	StreamID string  `json:"stream_id"`
	Minor    int     `json:"minor"`
	Major    int     `json:"major"`
	Patch    int     `json:"patch"`
}

// UnmarshalJSON accepts both the object results used by status requests and
// the scalar "ok" result returned by stream control/property requests.
func (r *result) UnmarshalJSON(data []byte) error {
	r.raw = append(r.raw[:0], data...)
	if !r.isObject() {
		return nil
	}
	type resultAlias result
	return json.Unmarshal(data, (*resultAlias)(r))
}

func (r result) isObject() bool {
	return len(bytes.TrimSpace(r.raw)) > 0 && bytes.TrimSpace(r.raw)[0] == '{'
}

func (r result) isOK() bool {
	return bytes.Equal(bytes.TrimSpace(r.raw), []byte(`"ok"`))
}

type client struct {
	Config    clientConfig `json:"config"`
	Connected bool         `json:"connected"`
	Id        string       `json:"id"`
	Host      host         `json:"host"`
}

type host struct {
	Arch string `json:"arch"`
	Ip   string `json:"ip"`
	Mac  string `json:"mac"`
	Name string `json:"name"`
	Os   string `json:"os"`
}

type clientConfig struct {
	Instance int    `json:"instance"`
	Latency  int    `json:"latency"`
	Name     string `json:"name"`
	Volume   volume `json:"volume"`
}

type volumeRequest struct {
	Id     string `json:"id"`
	Volume volume `json:"volume"`
}

type nameRequest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type latencyRequest struct {
	Id      string `json:"id"`
	Latency int    `json:"latency"`
}

type volume struct {
	Muted   bool `json:"muted"`
	Percent int  `json:"percent"`
}

type muteRequest struct {
	Id   string `json:"id"`
	Mute bool   `json:"mute"`
}

type groupStreamRequest struct {
	Id       string `json:"id"`
	StreamID string `json:"stream_id"`
}

type groupClientsRequest struct {
	Id      string   `json:"id"`
	Clients []string `json:"clients"`
}

type streamControlRequest struct {
	Id      string         `json:"id"`
	Command string         `json:"command"`
	Params  map[string]any `json:"params,omitempty"`
}

type streamPropertyRequest struct {
	Id       string `json:"id"`
	Property string `json:"property"`
	Value    any    `json:"value"`
}

type streamURIRequest struct {
	StreamURI string `json:"streamUri"`
}

type server struct {
	Groups  []group       `json:"groups"`
	Server  serverDetails `json:"server"`
	Streams []stream      `json:"streams"`
}

type serverDetails struct {
	Host host `json:"host"`
}

type group struct {
	Clients  []client `json:"clients"`
	Id       string   `json:"id"`
	Muted    bool     `json:"muted"`
	Name     string   `json:"name"`
	StreamId string   `json:"stream_id"`
}

type stream struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Uri    uri    `json:"uri"`
}

type uri struct {
	Fragment string   `json:"fragment"`
	Host     string   `json:"host"`
	Path     string   `json:"path"`
	Raw      string   `json:"raw"`
	Scheme   string   `json:"scheme"`
	Query    uriQuery `json:"query"`
}

type uriQuery struct {
	ChunkMs      string `json:"chunk_ms"`
	Codec        string `json:"codec"`
	Name         string `json:"name"`
	SampleFormat string `json:"sampleformat"`
}
