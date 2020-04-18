package main

import (
	"flag"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/server"
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

	//TODO: синглтоны
	ApiConfigControl := config.CreateApiConfigController()
	UserConfigControl := config.CreateUserConfigController()

	srv := &server.Server{
		IP:         ApiConfigControl.GetServerIP(),
		Port:       ApiConfigControl.GetServerPort(),
		ApiConfig:  ApiConfigControl,
		UserConfig: UserConfigControl,
	}
	srv.Run()
}
