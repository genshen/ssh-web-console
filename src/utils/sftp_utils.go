package utils

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"log"
	"sync"
)

type SftpNode Node // struct alias.

type SftpEntity struct {
	sshEntity  *SSHShellSession // from utils/ssh_utils
	sftpClient *sftp.Client     // sftp session created by sshEntity.client..
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
	sftpMutex       = new(sync.RWMutex)
	subscribers = make(map[string]SftpEntity)
)

func NewSftpEntity(user SftpNode, username string, auth ssh.AuthMethod) (SftpEntity, error) {
	sshEntity := SSHShellSession{
		Node: NewSSHNode(user.Host, user.Port),
	}
	// init ssh connection.
	err := sshEntity.Connect(username, auth)
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
	sftpMutex.Lock()
	//subscribers.PushBack(client)
	if c, ok := subscribers[key]; ok {
		c.Close() // if client have exists, close the client.
	}
	subscribers[key] = sftpEntity // store sftpEntity.
	sftpMutex.Unlock()
}

// make a copy of SftpEntity matched with given key.
// return sftpEntity and exist flag (bool).
func Fork(key string) (SftpEntity, bool) {
	sftpMutex.Lock()
	defer sftpMutex.Unlock()
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
	sftpMutex.Lock()
	defer sftpMutex.Unlock()
	//subscribers.PushBack(client)
	if c, ok := subscribers[key]; ok {
		return c.sftpClient
	} else {
		return nil
	}
}

// remove a sftp client by key.
func Leave(key string) {
	sftpMutex.Lock()
	//subscribers.PushBack(client)
	if c, ok := subscribers[key]; ok {
		c.Close()                // close the client.
		delete(subscribers, key) // remove from map.
	}
	sftpMutex.Unlock()
}
