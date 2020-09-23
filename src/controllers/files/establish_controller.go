package files

import (
	"github.com/genshen/ssh-web-console/src/models"
	"github.com/genshen/ssh-web-console/src/utils"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/ssh"
	"log"
	"math/rand"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
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
	ws, err := websocket.Accept(w, r, nil)
	if err != nil {
		http.Error(w, "Cannot setup WebSocket connection:", 400)
		log.Println("Error: Cannot setup WebSocket connection:", err)
		return
	}
	defer ws.Close(websocket.StatusNormalClosure, "closed")

	// add sftp client to list if success.
	user := session.Value.(models.UserInfo)
	sftpEntity, err := utils.NewSftpEntity(utils.SftpNode(utils.NewSSHNode(user.Host, user.Port)), user.Username, ssh.Password(user.Password))
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
	utils.Join(id.String(), sftpEntity) // note:key is not for user auth, but for identify different connections.
	defer utils.Leave(id.String())      // close sftp connection anf remove sftpEntity from list.

	wsjson.Write(r.Context(), ws, models.SftpWebSocketMessage{Type: models.SftpWebSocketMessageTypeID, Data: id.String()})

	// dispatch webSocket Messages.
	// process webSocket message one by one at present. todo improvement.
	for {
		_, _, err := ws.Read(r.Context())
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
