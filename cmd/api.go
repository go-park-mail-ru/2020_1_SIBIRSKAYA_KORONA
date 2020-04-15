package main

import (
	"flag"
	"log"
	"os"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/spf13/viper"
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
		log.Fatal(err)
	}

	logger.InitLogger()

	avatarDir := viper.GetString("frontend.public_dir") + viper.GetString("frontend.avatar_dir")

	_, avatarErr := os.Stat(avatarDir)
	if os.IsNotExist(avatarErr) {
		errDir := os.MkdirAll(avatarDir, os.ModePerm)
		if errDir != nil {
			logger.Fatal(errDir)
		}
	}
	logger.Info("Avatar static storage up!")

	srv := &server.Server{
		IP:   viper.GetString("server.ip"),
		Port: uint(viper.GetInt("server.port")),
	}
	srv.Run()
}
