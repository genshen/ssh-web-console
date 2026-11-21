package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/genshen/ssh-web-console/src/routers"
	"github.com/genshen/ssh-web-console/src/utils"
	"io/fs"
	"log"
	"net/http"
)

// The frontend build (web/build), copied to ./build before `go build`.
//
//go:embed all:build
var embeddedStatic embed.FS

var confFilePath string
var version bool

func init() {
	flag.StringVar(&confFilePath, "config", "conf/config.yaml", "filepath of config file.")
	flag.BoolVar(&version, "version", false, "show current version.")
}

func main() {
	flag.Parse()
	if version {
		fmt.Println("v0.4.0")
		return
	}
	if err := utils.InitConfig(confFilePath); err != nil {
		log.Fatal("config error,", err)
		return
	}
	staticFS, err := fs.Sub(embeddedStatic, "build")
	if err != nil {
		log.Fatal(err)
	}
	routers.Register(staticFS)
	log.Println("listening on port ", utils.Config.Site.ListenAddr)
	// listen http
	if err := http.ListenAndServe(utils.Config.Site.ListenAddr, nil); err != nil {
		log.Fatal(err)
		return
	}
}
