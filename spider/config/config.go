package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

//Mongo conf
type Mongo struct {
	Host string `json:"host"`
}

//WebDav 结构
type WebDav struct {
	Host      string   `json:"host"`
	User      string   `json:"user"`
	Password  string   `json:"password"`
	DataPath  string   `json:"date_path"`
	Files     []string `json:"files"`
	RemoteDir string   `json:"remote_dir"`
}

//Service 服务信息
type Service struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

//Redis conf
type Redis struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	Auth string `json:"password"`
}

//Config conf
type Config struct {
	Mongo   Mongo   `json:"mongodb"`
	WebDav  WebDav  `json:"webdav"`
	Service Service `json:"service"`
	Redis   Redis   `json:"redis"`
}

var cfg *Config
var configFile = "../../config/config.json"

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err.Error)
	}
	return dir
}

func getCfgFile() string {
	curPath := getCurrentDirectory()
	switch runtime.GOOS {
	case "windows":
		configFile = strings.Replace(configFile, "/", "\\", -1)
		curPath = curPath + "\\"
	default:
		curPath = curPath + "/"
	}
	return curPath + configFile
}

func GetConfig() *Config {
	if cfg == nil {
		cfg = &Config{}
		if data, err := ioutil.ReadFile(getCfgFile()); err != nil {
			panic(err.Error())
		} else {
			if err := json.Unmarshal(data, cfg); err != nil {
				panic(err.Error)
			}
		}
	}
	return cfg
}
