package main

import (
	_ "github.com/genshen/webConsole/src/routers"
	"github.com/astaxie/beego"
	"os"
	"golang.org/x/crypto/ssh"
	"github.com/genshen/webConsole/src/utils"
	"log"
)

func main() {
	beego.Run()
	//setupSSH()
}

func setupSSH() {
	check := func(err error, msg string) {
		if err != nil {
			log.Fatalf("%s error: %v", msg, err)
		}
	}

	sshEntity := utils.SSH{
		Node: utils.Node{
			Host: "ssh.hpc.gensh.me",
			Port: 22,
		},
	}
	_, err := sshEntity.Connect("genshen", "genshen1234")
	check(err,"connect")
	defer sshEntity.Client.Close()

	session, err := sshEntity.Client.NewSession()
	check(err, "new session")
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	err = session.RequestPty("xterm", 25, 100, modes)
	check(err, "request pty")

	err = session.Shell()
	check(err, "start shell")

	err = session.Wait()
	check(err, "return")
}