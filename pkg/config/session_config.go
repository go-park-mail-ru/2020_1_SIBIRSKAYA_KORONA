package config

import (
	"github.com/spf13/viper"
)

type SessionConfigController struct {
	origins    []string
	serverIp   string
	serverPort uint
}

func CreateSessionConfigController() *SessionConfigController {
	return &SessionConfigController{
		origins:    viper.GetStringSlice("cors.allowed_origins"),
		serverIp:   viper.GetString("server.ip"),
		serverPort: viper.GetUint("server.port"),
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
