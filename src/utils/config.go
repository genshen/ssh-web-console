package utils

import (
	"os"
	"log"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

const (
	KEY_SSH_IO_MODE = "ssh_io_mode"
)

var Config struct {
	Site struct {
		AppName      string `yaml:"app_name"`
		RunMode      string `yaml:"runmode"`
		DeployHost   string `yaml:"deploy_host"`
		ListenAddr   string `yaml:"listen_addr"`
		StaticPrefix string `yaml:"static_prefix"`
		StaticDir    string `yaml:"static_dir"`
		ViewsDir     string `yaml:"views_dir"`
	} `yaml:"site"`
	VPN struct {
		Enable bool `yaml:"enable"`
	} `yaml:"vpn_juniper"`
	SSH struct {
		IOMode                 int `yaml:"io_mode"`
		BufferCheckerCycleTime int `yaml:"buffer_checker_cycle_time"`
	} `yaml:"ssh"`
	Jwt struct {
		Secret        string `yaml:"jwt_secret"`
		TokenLifetime int64  `yaml:"token_lifetime"`
		Issuer        string `yaml:"issuer"`
		QueryTokenKey string `yaml:"query_token_key"`
	} `yaml:"jwt"`
}

func init() {
	f, err := os.Open("conf/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(content, &Config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
