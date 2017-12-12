package controllers

import (
	"io"
	"log"
	"bufio"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/genshen/webConsole/src/utils"
	"github.com/genshen/webConsole/src/models"
	"time"
	"bytes"
)

type SSHWebSocketHandle struct {
}

//to handle webSocket connection
func (c SSHWebSocketHandle) ServeAfterAuthenticated(w http.ResponseWriter, r *http.Request, claims *utils.Claims, session *utils.Session) {
	// init websocket connection
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		utils.Abort(w, "Not a websocket handshake", 400)
		log.Println("Error: Not a websocket handshake", 400)
		return
	} else if err != nil {
		http.Error(w, "Cannot setup WebSocket connection:", 400)
		log.Println("Error: Cannot setup WebSocket connection:", err)
		return
	}
	defer ws.Close()

	//setup ssh connection
	sshEntity := utils.SSH{
		Node: utils.Node{
			Host: claims.Host,
			Port: claims.Port,
		},
	}
	userInfo := session.Value.(models.UserInfo)
	_, err = sshEntity.Connect(userInfo.Username, userInfo.Password)
	if err != nil {
		utils.Abort(w, "Cannot setup ssh connection:", 500)
		log.Println("Error: Cannot setup ssh connection:", err)
		return
	}
	defer sshEntity.Close()

	cols := utils.GetQueryInt32(r, "cols", 120)
	rows := utils.GetQueryInt32(r, "rows", 32)

	//set ssh IO mode and ssh shell
	sshIOMode := utils.Config.SSH.IOMode
	if sshIOMode == utils.SSH_IO_MODE_CHANNEL {
		_, err = sshEntity.ConfigShellChannel(cols, rows)
	} else {
		_, err = sshEntity.ConfigShellSession(int(cols), int(rows))
	}
	if err != nil {
		log.Println("Error: configure ssh session error:", err)
		return
	}

	done := make(chan bool, 4)
	runeChan := make(chan rune)
	setDone := func() { done <- true }

	// most messages are ssh output,not webSocket input,so we add a buffer in function readMessageFromSSHServer.
	writeMessageToSSHServer := func(wc io.WriteCloser) { // read messages from webSocket
		defer setDone()
		for {
			msgType, p, err := ws.ReadMessage()
			if err != nil {
				log.Println("Error: error reading webSocket message:", err)
				return
			}
			if err = DispatchMessage(msgType, p, wc); err != nil {
				log.Println("Error: error sending data to ssh server:", err)
				return
			}
		}
	}

	// read turn from ssh server, and store it to byte buffer.
	readMessageFromSSHServer := func(reader io.Reader) {
		defer setDone()
		// read rune.
		br := bufio.NewReader(reader)
		for {
			r, size, err := br.ReadRune()
			if err != nil {
				log.Println("Error: error reading data from ssh server:", err)
				return
			}
			if size > 0 { // store rune to buffer. (?) may have bug: char '\', if not use buffer.
				runeChan <- r
			}
		}
	}

	writeBufferToWebSocket := func() {
		defer setDone()
		tick := time.NewTicker(time.Millisecond * time.Duration(utils.Config.SSH.BufferCheckerCycleTime)) // check buffer(if not empty,then write back to webSocket) every 120 ms.
		//for range time.Tick(120 * time.Millisecond){}
		defer tick.Stop()
		// r := make(chan rune)
		var buffer bytes.Buffer
		for {
			select {
			case <-tick.C:
				if buffer.Len() != 0 {
					err = ws.WriteMessage(websocket.TextMessage, []byte(buffer.Bytes()))
					if err != nil { //todo error
						log.Println("Error: error sending data via webSocket:", err)
						return
					}
				}
				buffer.Reset()
			case r := <-runeChan:
				buffer.WriteRune(r)
			}
		}
	}

	if sshIOMode == utils.SSH_IO_MODE_CHANNEL {
		go writeMessageToSSHServer(sshEntity.Channel)
		go readMessageFromSSHServer(sshEntity.Channel)
		go writeBufferToWebSocket()
	} else {
		go writeMessageToSSHServer(sshEntity.IO.StdIn)
		go readMessageFromSSHServer(sshEntity.IO.StdOut)
		go readMessageFromSSHServer(sshEntity.IO.StdErr)
		go writeBufferToWebSocket()
	}
	<-done
	log.Println("Info: websocket finished!")
}
