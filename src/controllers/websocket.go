package controllers

import (
	"context"
	"fmt"
	"github.com/genshen/ssh-web-console/src/models"
	"github.com/genshen/ssh-web-console/src/utils"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"time"
)

//const SSH_EGG = `genshen<genshenchu@gmail.com> https://github.com/genshen/sshWebConsole"`

type SSHWebSocketHandle struct {
	bufferFlushCycle int
}

func NewSSHWSHandle(bfc int) *SSHWebSocketHandle {
	var handle SSHWebSocketHandle
	handle.bufferFlushCycle = bfc
	return &handle
}

// clear session after ssh closed.
func (c *SSHWebSocketHandle) ShouldClearSessionAfterExec() bool {
	return true
}

// handle webSocket connection.
func (c *SSHWebSocketHandle) ServeAfterAuthenticated(w http.ResponseWriter, r *http.Request, claims *utils.Claims, session utils.Session) {
	// init webSocket connection
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		http.Error(w, "Cannot setup WebSocket connection:", 400)
		log.Println("Error: Cannot setup WebSocket connection:", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "closed")

	userInfo := session.Value.(models.UserInfo)
	cols := utils.GetQueryInt32(r, "cols", 120)
	rows := utils.GetQueryInt32(r, "rows", 32)
	sshAuth := ssh.Password(userInfo.Password)
	if err := c.SSHShellOverWS(r.Context(), conn, claims.Host, claims.Port, userInfo.Username, sshAuth, cols, rows); err != nil {
		log.Println("Error,", err)
		utils.Abort(w, err.Error(), 500)
	}
}

// ssh shell over websocket
// first,we establish a ssh connection to ssh server when a webSocket comes;
// then we deliver ssh data via ssh connection between browser and ssh server.
// That is, read webSocket data from browser (e.g. 'ls' command) and send data to ssh server via ssh connection;
// the other hand, read returned ssh data from ssh server and write back to browser via webSocket API.
func (c *SSHWebSocketHandle) SSHShellOverWS(ctx context.Context, ws *websocket.Conn, host string, port int, username string, auth ssh.AuthMethod, cols, rows uint32) error {
	//setup ssh connection
	sshEntity := utils.SSHShellSession{
		Node: utils.Node{
			Host: host,
			Port: port,
		},
	}
	// set io for ssh session
	var wsBuff WebSocketBufferWriter
	sshEntity.WriterPipe = &wsBuff

	var sshConn utils.SSHConnInterface = &sshEntity // set interface
	err := sshConn.Connect(username, auth)
	if err != nil {
		return fmt.Errorf("cannot setup ssh connection %w", err)
	}
	defer sshConn.Close()

	// config ssh
	sshSession, err := sshConn.Config(cols, rows)
	if err != nil {
		return fmt.Errorf("configure ssh error: %w", err)
	}

	// an egg:
	//if err := sshSession.Setenv("SSH_EGG", SSH_EGG); err != nil {
	//	log.Println(err)
	//}
	// after configure, the WebSocket is ok.
	defer wsBuff.Flush(ctx, websocket.MessageText, ws)

	done := make(chan bool, 3)
	setDone := func() { done <- true }

	// most messages are ssh output, not webSocket input
	writeMessageToSSHServer := func(wc io.WriteCloser) { // read messages from webSocket
		defer setDone()
		for {
			msgType, p, err := ws.Read(ctx)
			// if WebSocket is closed by some reason, then this func will return,
			// and 'done' channel will be set, the outer func will reach to the end.
			// then ssh session will be closed in defer.
			if err != nil {
				log.Println("Error: error reading webSocket message:", err)
				return
			}
			if err = DispatchMessage(sshSession, msgType, p, wc); err != nil {
				log.Println("Error: error write data to ssh server:", err)
				return
			}
		}
	}

	stopper := make(chan bool) // timer stopper
	// check webSocketWriterBuffer(if not empty,then write back to webSocket) every 120 ms.
	writeBufferToWebSocket := func() {
		defer setDone()
		tick := time.NewTicker(time.Millisecond * time.Duration(c.bufferFlushCycle))
		//for range time.Tick(120 * time.Millisecond){}
		defer tick.Stop()
		for {
			select {
			case <-tick.C:
				if err := wsBuff.Flush(ctx, websocket.MessageText, ws); err != nil {
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
		if err := sshSession.Wait(); err != nil {
			log.Println("ssh exist from server", err)
		}
		// if ssh is closed (wait returns), then 'done', web socket will be closed.
		// by the way, buffered data will be flushed before closing WebSocket.
	}()

	<-done
	stopper <- true // stop tick timer(if tick is finished by due to the bad WebSocket, this line will just only set channel(no bad effect). )
	log.Println("Info: websocket finished!")
	return nil
}
