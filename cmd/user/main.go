package main

import (
	"flag"
	"os"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/user/server"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/spf13/viper"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/config"
)

var opts struct {
	configPath string
}

func init() {
	flag.StringVar(&opts.configPath, "c", "", "path to configuration file")
	flag.StringVar(&opts.configPath, "config", "", "path to configuration file")
}

func main() {
	flag.Parse()
	viper.SetConfigFile(opts.configPath)
	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal(err)
	}
	logger.InitLogger()
	configControl := config.CreateUserConfigController()

	avatarDir := configControl.GetFrontendAvatarDir()

	_, avatarErr := os.Stat(avatarDir)
	if os.IsNotExist(avatarErr) {
		errDir := os.MkdirAll(avatarDir, os.ModePerm)
		if errDir != nil {
			logger.Fatal(errDir)
		}
	}

	srv := &server.Server{
		IP:     configControl.GetServerIP(),
		Port:   configControl.GetServerPort(),
		Config: configControl,
	}
	srv.Run()
}
