package models

const (
	WebSocketMessageTypeTerminal   = "terminal"
	WebSocketMessageTypeHeartbeat = "heartbeat"
)

type WebSocketMessage struct {
	Type       string `json:"type"`
	DataBase64 string `json:"data"` // json.RawMessage
}
