package controllers

import (
	"bytes"
	"context"
	"nhooyr.io/websocket"
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
func (w *WebSocketBufferWriter) Flush(ctx context.Context, messageType websocket.MessageType, ws *websocket.Conn) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.buffer.Len() != 0 {
		err := ws.Write(ctx, messageType, w.buffer.Bytes())
		if err != nil {
			return err
		}
		w.buffer.Reset()
	}
	return nil
}
