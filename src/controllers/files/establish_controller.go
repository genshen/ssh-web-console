package files

import (
	"net/http"
	"github.com/genshen/webConsole/src/utils"
	"github.com/genshen/webConsole/src/models"
	"log"
	"github.com/gorilla/websocket"
	"time"
	"github.com/oklog/ulid"
	"math/rand"
)

type SftpEstablish struct{}

func (e SftpEstablish) ShouldClearSessionAfterExec() bool {
	return false
}

// establish webSocket connection to browser to maintain connection with remote sftp server.
// If establish success, add sftp connection to a list.
// and then, handle all message from message (e.g.list files in one directory.).
func (e SftpEstablish) ServeAfterAuthenticated(w http.ResponseWriter, r *http.Request, claims *utils.Claims, session utils.Session) {
	// init webSocket connection
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		utils.Abort(w, "Not a webSocket handshake", 400)
		log.Println("Error: Not a websocket handshake", 400)
		return
	} else if err != nil {
		http.Error(w, "Cannot setup WebSocket connection:", 400)
		log.Println("Error: Cannot setup WebSocket connection:", err)
		return
	}
	defer ws.Close()

	// add sftp client to list if success.
	user := session.Value.(models.UserInfo)
	sftpEntity, err := NewSftpEntity(SftpNode{Host: user.Host, Port: user.Port}, user.Username, user.Password)
	if err != nil {
		http.Error(w, "Error while establishing sftp connection", 400)
		log.Println("Error: while establishing sftp connection", err)
		return
	}
	// generate unique id.
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	// add sftpEntity to list.
	Join(id.String(), sftpEntity) // note:key is not for user auth, but for identify different connections.
	defer Leave(id.String())      // close sftp connection anf remove sftpEntity from list.

	ws.WriteJSON(models.SftpWebSocketMessage{Type: models.SftpWebSocketMessageTypeID, Data: id.String()})

	// dispatch webSocket Messages.
	// process webSocket message one by one at present. todo improvement.
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			log.Println("Error: error reading webSocket message:", err)
			break
		}
		//if err = DispatchSftpMessage(msgType, p, sftpEntity.sftpClient); err != nil { // todo handle heartbeat message and so on.
		//	log.Println("Error: error write data to ssh server:", err)
		//	break
		//}
	}
}
