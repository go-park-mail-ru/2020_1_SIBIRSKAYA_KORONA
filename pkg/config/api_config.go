package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type ApiConfigController struct {
	origins      []string
	serverIp     string
	serverPort   uint
	db           string
	dbConnection string
}

func CreateApiConfigController() *ApiConfigController {
	dbHost := viper.GetString("database.host")
	dbUser := viper.GetString("database.user")
	dbPass := viper.GetString("database.password")
	dbName := viper.GetString("database.name")
	dbMode := viper.GetString("database.sslmode")

	return &ApiConfigController{
		origins:      viper.GetStringSlice("cors.allowed_origins"),
		serverIp:     viper.GetString("server.ip"),
		serverPort:   viper.GetUint("server.port"),
		db:           viper.GetString("database.dbms"),
		dbConnection: fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", dbHost, dbUser, dbPass, dbName, dbMode),
	}
}

func (cc *ApiConfigController) GetOriginsSlice() []string {
	return cc.origins
}

func (cc *ApiConfigController) GetServerIP() string {
	return cc.serverIp
}

func (cc *ApiConfigController) GetServerPort() uint {
	return cc.serverPort
}

func (cc *ApiConfigController) GetDB() string {
	return cc.db
}

func (cc *ApiConfigController) GetDBConnection() string {
	return cc.dbConnection
}
