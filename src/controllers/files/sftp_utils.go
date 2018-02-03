package files

import (
	"sync"
	"github.com/pkg/sftp"
	"github.com/genshen/webConsole/src/utils"
	"log"
)

type SftpNode utils.Node // struct alias.

type SftpEntity struct {
	sshEntity  *utils.SSH   // from utils/ssh_utils
	sftpClient *sftp.Client // sftp session created by sshEntity.Client..
}

// close sftp session and ssh client
func (con *SftpEntity) Close() {
	defer con.sshEntity.Close()

	err := con.sftpClient.Close()
	if err != nil { // todo for debug.
		log.Println(err)
	}
}

var (
	mutex       = new(sync.RWMutex)
	subscribers = make(map[string]SftpEntity)
)

func NewSftpEntity(user SftpNode, username, password string) (SftpEntity, error) {
	sshEntity := utils.SSH{
		Node: utils.Node{
			Host: user.Host,
			Port: user.Port,
		},
	}
	// init ssh connection.
	err := sshEntity.Connect(username, password)
	if err != nil {
		return SftpEntity{}, err
	}

	// make a new sftp client
	client, err := sftp.NewClient(sshEntity.Client)
	if err != nil {
		return SftpEntity{}, err
	}
	return SftpEntity{sshEntity: &sshEntity, sftpClient: client}, nil
}

// add a sftp client to subscribers list.
func Join(key string, sftpEntity SftpEntity) {
	mutex.Lock()
	//subscribers.PushBack(client)
	if c, ok := subscribers[key]; ok {
		c.Close() // if client have exists, close the client.
	}
	subscribers[key] = sftpEntity // store sftpEntity.
	mutex.Unlock()
}

// make a copy of SftpEntity matched with given key.
// return sftpEntity and exist flag (bool).
func Fork(key string) (SftpEntity, bool) {
	mutex.Lock()
	defer mutex.Unlock()
	//subscribers.PushBack(client)
	if c, ok := subscribers[key]; ok {
		return c, true
	} else {
		return SftpEntity{}, false
	}
}

// make a copy of SftpEntity matched with given key.
// return sftp.Client pointer or nil pointer.
func ForkSftpClient(key string) (*sftp.Client) {
	mutex.Lock()
	defer mutex.Unlock()
	//subscribers.PushBack(client)
	if c, ok := subscribers[key]; ok {
		return c.sftpClient
	} else {
		return nil
	}
}

// remove a sftp client by key.
func Leave(key string) {
	mutex.Lock()
	//subscribers.PushBack(client)
	if c, ok := subscribers[key]; ok {
		c.Close()                // close the client.
		delete(subscribers, key) // remove from map.
	}
	mutex.Unlock()
}
