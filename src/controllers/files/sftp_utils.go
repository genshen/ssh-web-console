package files

import (
	"github.com/genshen/ssh-web-console/src/utils"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"log"
	"sync"
)

type SftpNode utils.Node // struct alias.

type SftpEntity struct {
	sshEntity  *utils.SSHShellSession // from utils/ssh_utils
	sftpClient *sftp.Client           // sftp session created by sshEntity.client..
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
	sshEntity := utils.SSHShellSession{
		Node: utils.NewSSHNode(user.Host, user.Port),
	}
	// init ssh connection.
	err := sshEntity.Connect(username, ssh.Password(password))
	if err != nil {
		return SftpEntity{}, err
	}

	// make a new sftp client
	if sshClient, err := sshEntity.GetClient(); err != nil {
		return SftpEntity{}, err
	} else {
		client, err := sftp.NewClient(sshClient)
		if err != nil {
			return SftpEntity{}, err
		}
		return SftpEntity{sshEntity: &sshEntity, sftpClient: client}, nil
	}
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
// return sftp.client pointer or nil pointer.
func ForkSftpClient(key string) *sftp.Client {
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
