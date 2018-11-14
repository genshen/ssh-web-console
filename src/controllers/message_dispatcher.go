package controllers

import (
	"encoding/base64"
	"encoding/json"
	"github.com/genshen/webConsole/src/models"
	"io"
)

func DispatchMessage(messageType int, message []byte, wc io.WriteCloser) error {
	socketMsg := models.SSHWebSocketMessage{}
	if err := json.Unmarshal(message, &socketMsg); err != nil {
		return err
	}

	switch socketMsg.Type {
	case models.SSHWebSocketMessageTypeHeartbeat:
		return nil
	case models.SSHWebSocketMessageTypeTerminal:
		if decodeBytes, err := base64.StdEncoding.DecodeString(socketMsg.DataBase64); err != nil { // todo ignore error
			return err
		} else {
			if _, err := wc.Write(decodeBytes); err != nil {
				return err
			}
		}
	}
	return nil
}
