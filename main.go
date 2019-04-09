package main

import (
	_ "github.com/genshen/ssh-web-console/src/routers"
	"github.com/genshen/ssh-web-console/src/utils"
	"log"
	"net/http"
)

func main() {
	if err := utils.InitConfig("conf/config.yaml"); err != nil {
		log.Fatal("config error,", err)
		return
	}
	// listen http
	if err := http.ListenAndServe(utils.Config.Site.ListenAddr, nil); err != nil {
		log.Fatal(err)
		return
	}
}
