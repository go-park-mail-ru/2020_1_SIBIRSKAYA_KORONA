package main

import (
	"flag"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/server"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
	"github.com/spf13/viper"

	"log"
	"os"
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

	publicDirPath, exists := os.LookupEnv("DRELLO_PUBLIC_DIR")
	if !exists {
		log.Fatal("DRELLO_PUBLIC_DIR environment variable not exist")
	}

	avatarDir := publicDirPath + viper.GetString("frontend.avatar_dir")

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
