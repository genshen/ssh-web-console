package main

import (
	"flag"
	"github.com/genshen/ssh-web-console/src/routers"
	"github.com/genshen/ssh-web-console/src/utils"
	"log"
	"net/http"
)

var confFilePath string

func init() {
	flag.StringVar(&confFilePath, "config", "conf/config.yaml", "filepath of config file.")
}

func main() {
	flag.Parse()
	if err := utils.InitConfig(confFilePath); err != nil {
		log.Fatal("config error,", err)
		return
	}
	routers.Register()
	log.Println("listening on port ", utils.Config.Site.ListenAddr)
	// listen http
	if err := http.ListenAndServe(utils.Config.Site.ListenAddr, nil); err != nil {
		log.Fatal(err)
		return
	}
}
