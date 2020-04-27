package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type SessionConfigController struct {
	origins          []string
	serverIp         string
	serverPort       uint
	memcachedConnect string
}

func CreateSessionConfigController() *SessionConfigController {
	return &SessionConfigController{
		origins:          viper.GetStringSlice("cors.allowed_origins"),
		serverIp:         viper.GetString("server.ip"),
		serverPort:       viper.GetUint("server.port"),
		memcachedConnect: fmt.Sprintf("%s:%d", viper.GetString("memcached.host"), viper.GetUint("memcached.port")),
	}
}

func (cc *SessionConfigController) GetOriginsSlice() []string {
	return cc.origins
}

func (cc *SessionConfigController) GetServerIP() string {
	return cc.serverIp
}

func (cc *SessionConfigController) GetServerPort() uint {
	return cc.serverPort
}

func (cc *SessionConfigController) GetMemcachedConnect() string {
	return cc.memcachedConnect
}
