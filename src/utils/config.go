package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var Config struct {
	Site struct {
		AppName    string `yaml:"app_name"`
		RunMode    string `yaml:"runmode"`
		DeployHost string `yaml:"deploy_host"`
		ListenAddr string `yaml:"listen_addr"`
	} `yaml:"site"`
	Prod struct {
		StaticPrefix string `yaml:"static_prefix"` // http prefix of static and views files
	} `yaml:"prod"`
	Dev struct {
		StaticPrefix string `yaml:"static_prefix"` // https prefix of only static files
		//StaticPrefix string `yaml:"static_prefix"` // prefix of static files in dev mode.
		// redirect static files requests to this address, redirect "StaticPrefix" to "StaticRedirect + StaticPrefix"
		// for example, StaticPrefix is "static", StaticRedirect is "localhost:8080/dist",
		// this will redirect all requests having prefix "static" to "localhost:8080/dist/"
		StaticRedirect string `yaml:"static_redirect"`
		// http server will read static file from this dir if StaticRedirect is empty
		StaticDir   string `yaml:"static_dir"`
		ViewsPrefix string `yaml:"views_prefix"` // https prefix of only views files
		// path of view files (we can not redirect view files) to be served.
		ViewsDir string `yaml:"views_dir"` // todo
	} `yaml:"dev"`
	SSH struct {
		BufferCheckerCycleTime int `yaml:"buffer_checker_cycle_time"`
	} `yaml:"ssh"`
	Jwt struct {
		Secret        string `yaml:"jwt_secret"`
		TokenLifetime int64  `yaml:"token_lifetime"`
		Issuer        string `yaml:"issuer"`
		QueryTokenKey string `yaml:"query_token_key"`
	} `yaml:"jwt"`
}

func InitConfig(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &Config)
	if err != nil {
		return err
	}
	return nil
}
