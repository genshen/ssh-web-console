package utils

import (
	"errors"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"strconv"
)

const (
	SSH_IO_MODE_CHANNEL = 0
	SSH_IO_MODE_SESSION = 1
)

type SSHConnInterface interface {
	// close ssh connection
	Close()
	// connect using username and password
	Connect(username string, auth ssh.AuthMethod) error
	// config connection after connected and may also create a ssh session.
	Config(cols, rows uint32) (*ssh.Session, error)
}

type Node struct {
	Host   string // host, e.g: ssh.example.com
	Port   int    //port,default value is 22
	client *ssh.Client
}

func NewSSHNode(host string, port int) Node {
	return Node{Host: host, Port: port, client: nil}
}

func (node *Node) GetClient() (*ssh.Client, error) {
	if node.client == nil {
		return nil, errors.New("client is not set")
	}
	return node.client, nil
}

//see: http://www.nljb.net/default/Go-SSH-%E4%BD%BF%E7%94%A8/
// establish a ssh connection. if success return nil, than can operate ssh connection via pointer Node.client in struct Node.
func (node *Node) Connect(username string, auth ssh.AuthMethod) error {
	//var hostKey ssh.PublicKey

	// An SSH client is represented with a ClientConn.
	//
	// To authenticate with the remote server you must pass at least one
	// implementation of AuthMethod via the Auth field in ClientConfig,
	// and provide a HostKeyCallback.
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			auth,
		},
		//HostKeyCallback: ssh.FixedHostKey(hostKey),
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", node.Host+":"+strconv.Itoa(node.Port), config)
	if err != nil {
		return err
	}
	node.client = client
	return nil
}

// connect to ssh server using ssh session.
type SSHShellSession struct {
	Node
	// calling Write() to write data to ssh server
	StdinPipe io.WriteCloser
	// Write() be called to receive data from ssh server
	WriterPipe io.Writer
	session    *ssh.Session
}

// setup ssh shell session
// set SSHShellSession.session and StdinPipe from created session here.
// and Session.Stdout and Session.Stderr are also set for outputting.
// Return value is a pointer of ssh session which is created by ssh client for shell interaction.
// If it has error in this func, ssh session will be nil.
func (s *SSHShellSession) Config(cols, rows uint32) (*ssh.Session, error) {
	session, err := s.client.NewSession()
	if err != nil {
		return nil, err
	}
	s.session = session

	// we set stdin, then we can write data to ssh server via this stdin.
	// but, as for reading data from ssh server, we can set Session.Stdout and Session.Stderr
	// to receive data from ssh server, and write back to somewhere.
	if stdin, err := session.StdinPipe(); err != nil {
		log.Fatal("failed to set IO stdin: ", err)
		return nil, err
	} else {
		// in fact, stdin it is channel.
		s.StdinPipe = stdin
	}

	// set writer, such the we can receive ssh server's data and write the data to somewhere specified by WriterPipe.
	if s.WriterPipe == nil {
		return nil, errors.New("WriterPipe is nil")
	}
	session.Stdout = s.WriterPipe
	session.Stderr = s.WriterPipe

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echo
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Request pseudo terminal
	if err := session.RequestPty("xterm", int(rows), int(cols), modes); err != nil {
		log.Fatal("request for pseudo terminal failed: ", err)
		return nil, err
	}
	// Start remote shell
	if err := session.Shell(); err != nil {
		log.Fatal("failed to start shell: ", err)
		return nil, err
	}
	return session,nil
}

func (s *SSHShellSession) Close() {
	if s.session != nil {
		s.session.Close()
	}

	if s.client != nil {
		s.client.Close()
	}
}

// deprecated, use session SSHShellSession instead
// connect to ssh server using channel.
type SSHShellChannel struct {
	Node
	Channel ssh.Channel
}

type ptyRequestMsg struct {
	Term     string
	Columns  uint32
	Rows     uint32
	Width    uint32
	Height   uint32
	Modelist string
}

func (ch *SSHShellChannel) Config(cols, rows uint32) error {
	channel, requests, err := ch.client.Conn.OpenChannel("session", nil)
	if err != nil {
		return err
	}
	ch.Channel = channel

	go func() {
		for req := range requests {
			if req.WantReply {
				req.Reply(false, nil)
			}
		}
	}()

	//see https://github.com/golang/crypto/blob/master/ssh/example_test.go
	modes := ssh.TerminalModes{ //todo configure
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	var modeList []byte
	for k, v := range modes {
		kv := struct {
			Key byte
			Val uint32
		}{k, v}
		modeList = append(modeList, ssh.Marshal(&kv)...)
	}
	modeList = append(modeList, 0)
	req := ptyRequestMsg{ //todo configure
		Term:     "xterm",
		Columns:  cols,
		Rows:     rows,
		Width:    cols * 8,
		Height:   rows * 8,
		Modelist: string(modeList),
	}

	ok, err := channel.SendRequest("pty-req", true, ssh.Marshal(&req))
	if !ok || err != nil {
		return errors.New("error sending pty-request" +
			func() (string) {
				if err == nil {
					return ""
				}
				return err.Error()
			}())
	}

	ok, err = channel.SendRequest("shell", true, nil)
	if !ok || err != nil {
		return errors.New("error sending shell-request" +
			func() (string) {
				if err == nil {
					return ""
				}
				return err.Error()
			}())
	}
	return nil
}
