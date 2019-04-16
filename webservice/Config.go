package webservice

import (
	"io/ioutil"
	"encoding/json"
	"log"
)

// 程序配置
type Config struct {
	WsPort int `json:"wsPort"`
	WsReadTimeout int `json:"wsReadTimeout"`
	WsWriteTimeout int `json:"wsWriteTimeout"`
	WsInChannelSize int `json:"wsInChannelSize"`
	WsOutChannelSize int `json:"wsOutChannelSize"`
	WsHeartbeatInterval int `json:"wsHeartbeatInterval"`
	ServicePort int `json:"servicePort"`
	ServiceReadTimeout int `json:"serviceReadTimeout"`
	ServiceWriteTimeout int `json:"serviceWriteTimeout"`
	MaxJoinRoom int`json:"maxJoinRoom"`
	ServerPem string `json:"serverPem"`
	ServerKey string `json:"serverKey"`
}

var (
	G_config *Config
)

func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf Config
	)

	if content, err = ioutil.ReadFile(filename); err != nil {
		log.Fatalf("Config: err io.read : %v", err)
		return
	}

	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}

	G_config = &conf
	return
}