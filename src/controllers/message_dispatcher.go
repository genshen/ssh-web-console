package controllers

import (
	"encoding/base64"
	"encoding/json"
	"github.com/genshen/ssh-web-console/src/models"
	"golang.org/x/crypto/ssh"
	"io"
	"nhooyr.io/websocket"
)

func DispatchMessage(sshSession *ssh.Session, messageType websocket.MessageType, wsData []byte, wc io.WriteCloser) error {
	var socketData json.RawMessage
	socketStream := models.SSHWebSocketMessage{
		Data: &socketData,
	}

	if err := json.Unmarshal(wsData, &socketStream); err != nil {
		return nil // skip error
	}

	switch socketStream.Type {
	case models.SSHWebSocketMessageTypeHeartbeat:
		return nil
	case models.SSHWebSocketMessageTypeResize:
		var resize models.WindowResize
		if err := json.Unmarshal(socketData, &resize); err != nil {
			return nil // skip error
		}
		sshSession.WindowChange(resize.Rows, resize.Cols)
	case models.SSHWebSocketMessageTypeTerminal:
		var message models.TerminalMessage
		if err := json.Unmarshal(socketData, &message); err != nil {
			return nil
		}
		if decodeBytes, err := base64.StdEncoding.DecodeString(message.DataBase64); err != nil { // todo ignore error
			return nil // skip error
		} else {
			if _, err := wc.Write(decodeBytes); err != nil {
				return err
			}
		}
	}
	return nil
}
