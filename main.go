package main

import (
	"flag"
	"fmt"
	"github.com/genshen/ssh-web-console/src/routers"
	"github.com/genshen/ssh-web-console/src/utils"
	"log"
	"net/http"
)

var confFilePath string
var version bool

func init() {
	flag.StringVar(&confFilePath, "config", "conf/config.yaml", "filepath of config file.")
	flag.BoolVar(&version, "version", false, "show current version.")
}

func main() {
	flag.Parse()
	if version {
		fmt.Println("v0.2.2")
		return
	}
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
