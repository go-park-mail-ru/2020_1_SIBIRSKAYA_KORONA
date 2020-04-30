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
	metricURL        string
	service          string
}

func CreateSessionConfigController() *SessionConfigController {
	return &SessionConfigController{
		origins:          viper.GetStringSlice("cors.allowed_origins"),
		serverIp:         viper.GetString("server.ip"),
		serverPort:       viper.GetUint("server.port"),
		memcachedConnect: fmt.Sprintf("%s:%d", viper.GetString("memcached.host"), viper.GetUint("memcached.port")),
		metricURL:        viper.GetString("metrics.url"),
		service:          viper.GetString("metrics.service"),
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

func (cc *SessionConfigController) GetMetricsURL() string {
	return cc.metricURL
}

func (cc *SessionConfigController) GetServiceName() string {
	return cc.service
}
