package rpc

// BaseMessage is the base message for a rpc message
type BaseMessage struct {
	Testing bool   `json:"testing"`
	Method  string `json:"method"`
}
