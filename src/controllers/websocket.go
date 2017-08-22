package controllers

import (
	"github.com/gorilla/websocket"
	"net/http"
	"github.com/astaxie/beego"
	"github.com/genshen/webConsole/src/utils"
	"bufio"
	"io"
	"github.com/genshen/webConsole/src/models"
)

type WebSocketController struct {
	BaseController
}

//to handle webSocket connection
func (this *WebSocketController) SSHWebSocketHandle() {
	this.EnableRender = false
	ws, err := websocket.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}

	defer ws.Close()

	v := this.GetSession("userinfo")
	if v == nil {
		beego.Error("Cannot get Session data:", err)
		return
	}
	user := v.(models.UserInfo)
	//setup ssh connection
	sshEntity := utils.SSH{
		Node: utils.Node{
			Host: user.Host,
			Port: user.Port,
		},
	}
	_, err = sshEntity.Connect(user.Username, user.Password)
	if err != nil {
		beego.Error("Cannot setup ssh connection:", err)
		return
	}

	cols, err := this.GetUint32("cols", 120)
	if err != nil {
		beego.Error("get params cols error:", err)
		return
	}
	rows, err := this.GetUint32("rows", 32)
	if err != nil {
		beego.Error("get params cols error:", err)
		return
	}

	//set ssh IO mode and ssh shell
	sshIOMode := beego.AppConfig.DefaultInt(utils.KEY_SSH_IO_MODE, utils.SSH_IO_MODE_CHANNEL)
	if sshIOMode == utils.SSH_IO_MODE_CHANNEL {
		_, err = sshEntity.ConfigShellChannel(cols, rows)
	} else {
		_, err = sshEntity.ConfigShellSession(int(cols), int(rows))
	}
	if err != nil {
		beego.Error("configure ssh session error:", err)
		return
	}

	defer sshEntity.Close()

	done := make(chan bool, 3)
	setDone := func() { done <- true }

	writeMessageToSSHServer := func(wc io.WriteCloser) { //read messages from webSocket
		defer setDone()
		for {
			_, p, err := ws.ReadMessage()
			if err != nil {
				beego.Error("error reading webSocket message:", err)
				return
			}
			_, err = wc.Write(p)
			if err != nil {
				beego.Error("error sending data to ssh server:", err)
				return
			}
		}
	}

	readMessageFromSSHServer := func(reader io.Reader) {
		br := bufio.NewReader(reader)
		//buf := []byte{}
		go func() {
			defer setDone()
			for {
				r, size, err := br.ReadRune()
				if err != nil {
					beego.Error("error reading data from ssh server:", err)
					return
				}
				if size > 0 {
					//if string(r) == "\\" { //todo bug: char '\'
					//	continue
					//}
					err = ws.WriteMessage(websocket.TextMessage, []byte(string(r)))
					if err != nil { //todo error
						beego.Error("error sending data via webSocket:", err)
						return
					}
				}
			}
		}()
	}

	if sshIOMode == utils.SSH_IO_MODE_CHANNEL {
		go writeMessageToSSHServer(sshEntity.Channel);
		go readMessageFromSSHServer(sshEntity.Channel)
	} else {
		go writeMessageToSSHServer(sshEntity.IO.StdIn);
		go readMessageFromSSHServer(sshEntity.IO.StdOut)
		go readMessageFromSSHServer(sshEntity.IO.StdErr)
	}
	<-done
	beego.Info("websocket finished!")
}
