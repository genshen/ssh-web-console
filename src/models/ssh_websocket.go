package models

const (
	SSHWebSocketMessageTypeTerminal  = "terminal"
	SSHWebSocketMessageTypeHeartbeat = "heartbeat"
	SSHWebSocketMessageTypeResize    = "resize"
)

type SSHWebSocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"` // json.RawMessage
}

// normal terminal message
type TerminalMessage struct {
	DataBase64 string `json:"base64"`
}

// terminal window resize
type WindowResize struct {
	Cols int `json:"cols"`
	Rows int `json:"rows"`
}
