package main

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

type clientId struct {
	Id string `json:"id"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type result struct {
	Client *client `json:"client"`
	Server *server `json:"server"`
}

type client struct {
	Config    clientConfig `json:"config"`
	Connected bool         `json:"connected"`
	Id        string       `json:"id"`
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

type volume struct {
	Muted   bool `json:"muted"`
	Percent int  `json:"percent"`
}

type server struct {
	Groups  []group `json:"groups"`
	Server  any     `json:"server"`
	Streams []any   `json:"streams"`
}

type group struct {
	Clients  []client `json:"clients"`
	Id       string   `json:"id"`
	Muted    bool     `json:"muted"`
	Name     string   `json:"name"`
	StreamId string   `json:"stream_id"`
}
