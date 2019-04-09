package main

import (
	"github.com/genshen/ssh-web-console/src/routers"
	"github.com/genshen/ssh-web-console/src/utils"
	"log"
	"net/http"
)

func main() {
	if err := utils.InitConfig("conf/config.yaml"); err != nil {
		log.Fatal("config error,", err)
		return
	}
	routers.Register()
	log.Println("listening on port ",utils.Config.Site.ListenAddr)
	// listen http
	if err := http.ListenAndServe(utils.Config.Site.ListenAddr, nil); err != nil {
		log.Fatal(err)
		return
	}
}
