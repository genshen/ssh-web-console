package models

const (
	SftpWebSocketMessageTypeHeartbeat = "heartbeat"
	SftpWebSocketMessageTypeID        = "cid"
)

type SftpWebSocketMessage struct {
	Type string        `json:"type"`
	Data interface{} `json:"data"` // json.RawMessage
}
