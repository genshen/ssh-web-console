package controllers

import (
	"bytes"
	"github.com/gorilla/websocket"
	"sync"
)

// copy data from WebSocket to ssh server
// and copy data from ssh server to WebSocket

// write data to WebSocket
// the data comes from ssh server.
type WebSocketBufferWriter struct {
	buffer bytes.Buffer
	mu     sync.Mutex
}

// implement Write interface to write bytes from ssh server into bytes.Buffer.
func (w *WebSocketBufferWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Write(p)
}

// flush all data in this buff into WebSocket.
func (w *WebSocketBufferWriter) Flush(messageType int, ws *websocket.Conn) error {
	if w.buffer.Len() != 0 {
		err := ws.WriteMessage(messageType, []byte(w.buffer.Bytes()))
		if err != nil {
			return err
		}
		w.buffer.Reset()
	}
	return nil
}
