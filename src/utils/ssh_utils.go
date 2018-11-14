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

type Node struct {
	Host string // host, e.g: ssh.example.com
	Port int    //port,default value is 22
}

type SSH struct {
	Node Node
	IO   struct {
		StdIn  io.WriteCloser
		StdOut io.Reader
		StdErr io.Reader
	}
	Client     *ssh.Client
	Channel    ssh.Channel // used only in channel mode.
	hasChannel bool // can be true only in channel mode.
	Session    *ssh.Session // used only in session mode.
}

//see: http://www.nljb.net/default/Go-SSH-%E4%BD%BF%E7%94%A8/
// establish a ssh connection. if success return nil, than can operate ssh connection via pointer SSH.Client in struct SSH.
func (s *SSH) Connect(username, password string) error {
	//var hostKey ssh.PublicKey

	// An SSH client is represented with a ClientConn.
	//
	// To authenticate with the remote server you must pass at least one
	// implementation of AuthMethod via the Auth field in ClientConfig,
	// and provide a HostKeyCallback.
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		//HostKeyCallback: ssh.FixedHostKey(hostKey),
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil;
		},
	}

	client, err := ssh.Dial("tcp", s.Node.Host+":"+strconv.Itoa(s.Node.Port), config)
	if err != nil {
		return err
	}
	s.Client = client
	return nil
}

func (s *SSH) Close() {
	if s.hasChannel {
		s.Channel.Close()
	}

	if s.Session != nil {
		s.Session.Close()
	}

	if s.Client != nil {
		s.Client.Close()
	}
}

type ptyRequestMsg struct {
	Term     string
	Columns  uint32
	Rows     uint32
	Width    uint32
	Height   uint32
	Modelist string
}

// deprecated.
// use session ConfigShellSession instead.
func (this *SSH) ConfigShellChannel(cols, rows uint32) (ssh.Channel, error) {
	channel, requests, err := this.Client.Conn.OpenChannel("session", nil)
	if err != nil {
		return nil, err
	}

	this.hasChannel = true
	this.Channel = channel
	this.IO.StdIn = channel
	this.IO.StdOut = channel
	this.IO.StdErr = channel

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
		return nil, errors.New("error sending pty-request" +
			func() (string) {
				if err == nil {
					return ""
				}
				return err.Error()
			}())
	}

	ok, err = channel.SendRequest("shell", true, nil)
	if !ok || err != nil {
		return nil, errors.New("error sending shell-request" +
			func() (string) {
				if err == nil {
					return ""
				}
				return err.Error()
			}())
	}

	return channel, nil
}

// setup ssh shell session
func (this *SSH) ConfigShellSession(cols, rows int) (*ssh.Session, error) {
	session, err := this.Client.NewSession()
	if err != nil {
		return nil, err
	}

	this.Session = session

	err = this.setSessionInputOutput()
	if err != nil {
		log.Fatal("failed to set IO: ", err)
		return nil, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echo
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Request pseudo terminal
	if err := session.RequestPty("xterm", rows, cols, modes); err != nil {
		log.Fatal("request for pseudo terminal failed: ", err)
		return nil, err
	}
	// Start remote shell
	if err := session.Shell(); err != nil {
		log.Fatal("failed to start shell: ", err)
		return nil, err
	}
	return session, nil
}

func (this *SSH) setSessionInputOutput() (error) {
	stdin, err := this.Session.StdinPipe() // in fact, it is channel.
	if err != nil {
		return err
	}
	this.IO.StdIn = stdin

	stdout, err := this.Session.StdoutPipe() // in fact, it is channel.
	if err != nil {
		return err
	}
	this.IO.StdOut = stdout

	stderr, _ := this.Session.StderrPipe() // in fact, it is channel.
	if err != nil {
		return err
	}
	this.IO.StdErr = stderr
	return nil
}
