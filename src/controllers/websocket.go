package controllers

import (
	"github.com/genshen/ssh-web-console/src/models"
	"github.com/genshen/ssh-web-console/src/utils"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"time"
)

//const SSH_EGG = `genshen<genshenchu@gmail.com> https://github.com/genshen/sshWebConsole"`

type SSHWebSocketHandle struct {
	upgrader websocket.Upgrader
}

func NewSSHWSHandle() *SSHWebSocketHandle {
	var handle SSHWebSocketHandle
	handle.upgrader.ReadBufferSize = 1024
	handle.upgrader.WriteBufferSize = 1024
	return &handle
}

// clear session after ssh closed.
func (c *SSHWebSocketHandle) ShouldClearSessionAfterExec() bool {
	return true
}

// handle webSocket connection.
// first,we establish a ssh connection to ssh server when a webSocket comes;
// then we deliver ssh data via ssh connection between browser and ssh server.
// That is, read webSocket data from browser (e.g. 'ls' command) and send data to ssh server via ssh connection;
// the other hand, read returned ssh data from ssh server and write back to browser via webSocket API.
func (c *SSHWebSocketHandle) ServeAfterAuthenticated(w http.ResponseWriter, r *http.Request, claims *utils.Claims, session utils.Session) {
	// init webSocket connection
	ws, err := c.upgrader.Upgrade(w, r, nil)
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
	sshEntity := utils.SSHShellSession{
		Node: utils.Node{
			Host: claims.Host,
			Port: claims.Port,
		},
	}
	// set io for ssh session
	var wsBuff WebSocketBufferWriter
	sshEntity.WriterPipe = &wsBuff

	var sshConn utils.SSHConnInterface = &sshEntity // set interface
	userInfo := session.Value.(models.UserInfo)
	err = sshConn.Connect(userInfo.Username, userInfo.Password)
	if err != nil {
		utils.Abort(w, "Cannot setup ssh connection:", 500)
		log.Println("Error: Cannot setup ssh connection:", err)
		return
	}
	defer sshConn.Close()

	// config ssh
	cols := utils.GetQueryInt32(r, "cols", 120)
	rows := utils.GetQueryInt32(r, "rows", 32)
	if err = sshConn.Config(cols, rows); err != nil {
		log.Println("Error: configure ssh error:", err)
		return
	}

	// an egg:
	//if err := sshEntity.Session.Setenv("SSH_EGG", SSH_EGG); err != nil {
	//	log.Println(err)
	//}
	// after configure, the WebSocket is ok.
	defer wsBuff.Flush(websocket.TextMessage, ws)

	done := make(chan bool, 3)
	setDone := func() { done <- true }

	// most messages are ssh output, not webSocket input
	writeMessageToSSHServer := func(wc io.WriteCloser) { // read messages from webSocket
		defer setDone()
		for {
			msgType, p, err := ws.ReadMessage()
			// if WebSocket is closed by some reason, then this func will return,
			// and 'done' channel will be set, the outer func will reach to the end.
			// then ssh session will be closed in defer.
			if err != nil {
				log.Println("Error: error reading webSocket message:", err)
				return
			}
			if err = DispatchMessage(sshEntity.Session, msgType, p, wc); err != nil {
				log.Println("Error: error write data to ssh server:", err)
				return
			}
		}
	}

	stopper := make(chan bool) // timer stopper
	// check webSocketWriterBuffer(if not empty,then write back to webSocket) every 120 ms.
	writeBufferToWebSocket := func() {
		defer setDone()
		tick := time.NewTicker(time.Millisecond * time.Duration(utils.Config.SSH.BufferCheckerCycleTime))
		//for range time.Tick(120 * time.Millisecond){}
		defer tick.Stop()
		for {
			select {
			case <-tick.C:
				if err := wsBuff.Flush(websocket.TextMessage, ws); err != nil {
					log.Println("Error: error sending data via webSocket:", err)
					return
				}
			case <-stopper:
				return
			}
		}
	}

	go writeMessageToSSHServer(sshEntity.StdinPipe)
	go writeBufferToWebSocket()
	go func() {
		defer setDone()
		if err := sshEntity.Session.Wait(); err != nil {
			log.Println("ssh exist from server", err)
		}
		// if ssh is closed (wait returns), then 'done', web socket will be closed.
		// by the way, buffered data will be flushed before closing WebSocket.
	}()

	<-done
	stopper <- true // stop tick timer(if tick is finished by due to the bad WebSocket, this line will just only set channel(no bad effect). )
	log.Println("Info: websocket finished!")
}
