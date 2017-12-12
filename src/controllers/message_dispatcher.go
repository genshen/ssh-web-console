package controllers

import (
	"io"
	"encoding/json"
	"encoding/base64"
	"github.com/genshen/webConsole/src/models"
)

func DispatchMessage(messageType int, message []byte, wc io.WriteCloser) error {
	socketMsg := models.WebSocketMessage{}
	if err := json.Unmarshal(message, &socketMsg); err != nil {
		return nil // todo ignore unmarshal error
	}

	switch socketMsg.Type {
	case models.WebSocketMessageTypeHeartbeat:
		return nil
	case models.WebSocketMessageTypeTerminal:
		if decodeBytes, err := base64.StdEncoding.DecodeString(socketMsg.DataBase64); err == nil { // todo ignore error
			wc.Write(decodeBytes)
		}
	}
	return nil
}
