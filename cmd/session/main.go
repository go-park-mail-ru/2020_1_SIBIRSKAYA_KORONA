package main

import (
	"flag"
	"log"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/microservices/session/server"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/config"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"

	"github.com/spf13/viper"
)

var opts struct {
	configPath string
}

func main() {
	flag.StringVar(&opts.configPath, "c", "", "path to configuration file")
	flag.StringVar(&opts.configPath, "config", "", "path to configuration file")

	flag.Parse()
	viper.SetConfigFile(opts.configPath)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	logger.InitLogger()
	configControl := config.CreateSessionConfigController()

	srv := &server.Server{
		IP:     configControl.GetServerIP(),
		Port:   configControl.GetServerPort(),
		Config: configControl,
	}
	srv.Run()
}
