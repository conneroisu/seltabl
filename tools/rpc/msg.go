package rpc

// BaseMessage is the base message type
type BaseMessage struct {
	Testing bool   `json:"testing"`
	Method  string `json:"method"`
}
