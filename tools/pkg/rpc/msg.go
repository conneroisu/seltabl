package rpc

type BaseMessage struct {
	Testing bool   `json:"testing"`
	Method  string `json:"method"`
}