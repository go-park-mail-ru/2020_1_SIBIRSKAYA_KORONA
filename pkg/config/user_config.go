package config

import (
	"fmt"
	"os"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/spf13/viper"
)

type UserConfigController struct {
	origins         []string
	frontStorageURL string
	frontPublicDir  string
	avatarDir       string
	serverIp        string
	serverPort      uint
	db              string
	dbConnection    string
	metricURL       string
	service         string
}

func CreateUserConfigController() *UserConfigController {
	publicDirPath, exists := os.LookupEnv("DRELLO_PUBLIC_DIR")
	if !exists {
		logger.Fatal("DRELLO_PUBLIC_DIR environment variable not exist")
	}

	dbHost := viper.GetString("database.host")
	dbUser := viper.GetString("database.user")
	dbPass := viper.GetString("database.password")
	dbName := viper.GetString("database.name")
	dbMode := viper.GetString("database.sslmode")

	return &UserConfigController{
		origins: viper.GetStringSlice("cors.allowed_origins"),
		frontStorageURL: fmt.Sprintf("%s://%s:%s%s",
			viper.GetString("frontend.protocol"),
			viper.GetString("frontend.ip"),
			viper.GetString("frontend.port"),
			viper.GetString("frontend.avatar_dir")),
		frontPublicDir: publicDirPath,
		avatarDir:      viper.GetString("frontend.avatar_dir"),
		serverIp:       viper.GetString("server.ip"),
		serverPort:     viper.GetUint("server.port"),
		db:             viper.GetString("database.dbms"),
		dbConnection:   fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", dbHost, dbUser, dbPass, dbName, dbMode),
		metricURL:      viper.GetString("metrics.url"),
		service:        viper.GetString("metrics.service"),
	}
}

func (cc *UserConfigController) GetFrontendPublicDir() string {
	return cc.frontPublicDir
}

func (cc *UserConfigController) GetFrontendAvatarDir() string {
	return cc.frontPublicDir + cc.avatarDir
}

func (cc *UserConfigController) GetFrontendStorageURL() string {
	return cc.frontStorageURL
}

func (cc *UserConfigController) GetOriginsSlice() []string {
	return cc.origins
}

func (cc *UserConfigController) GetServerIP() string {
	return cc.serverIp
}

func (cc *UserConfigController) GetServerPort() uint {
	return cc.serverPort
}

func (cc *UserConfigController) GetDB() string {
	return cc.db
}

func (cc *UserConfigController) GetDBConnection() string {
	return cc.dbConnection
}

func (cc *UserConfigController) GetMetricsURL() string {
	return cc.metricURL
}

func (cc *UserConfigController) GetServiceName() string {
	return cc.service
}
