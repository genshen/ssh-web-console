package models

const (
	SSHWebSocketMessageTypeTerminal  = "terminal"
	SSHWebSocketMessageTypeHeartbeat = "heartbeat"
)

type SSHWebSocketMessage struct {
	Type       string `json:"type"`
	DataBase64 string `json:"data"` // json.RawMessage
}
