package webservice

import (
	"encoding/json"
	"io/ioutil"
)

// 程序配置
type Config struct {
	BackServiceTls      bool   `json:"backServiceTls"`
	WsPort              int    `json:"wsPort"`
	WsReadTimeout       int    `json:"wsReadTimeout"`
	WsWriteTimeout      int    `json:"wsWriteTimeout"`
	WsInChannelSize     int    `json:"wsInChannelSize"`
	WsOutChannelSize    int    `json:"wsOutChannelSize"`
	WsHeartbeatInterval int    `json:"wsHeartbeatInterval"`
	ServicePort         int    `json:"servicePort"`
	ServiceReadTimeout  int    `json:"serviceReadTimeout"`
	ServiceWriteTimeout int    `json:"serviceWriteTimeout"`
	MaxJoinRoom         int    `json:"maxJoinRoom"`
	ServerPem           string `json:"serverPem"`
	ServerKey           string `json:"serverKey"`
}

var (
	G_config *Config
)

func defaultConfig() {
	G_config = &Config{
		WsPort:              123,
		WsReadTimeout:       2000,
		WsWriteTimeout:      2000,
		WsInChannelSize:     1000,
		WsOutChannelSize:    1000,
		WsHeartbeatInterval: 60,
		ServicePort:         456,
		ServiceReadTimeout:  2000,
		ServiceWriteTimeout: 2000,
		MaxJoinRoom:         5,
		BackServiceTls:      false,
	}
}

func init() {
	var (
		content []byte
		conf    Config
		err     error
	)
	if content, err = ioutil.ReadFile("application.json"); err != nil {
		defaultConfig()
		return
	}
	if err = json.Unmarshal(content, &conf); err == nil {
		G_config = &conf
	}
	return
}
