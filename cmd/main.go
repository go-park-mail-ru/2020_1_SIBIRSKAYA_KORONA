package main

import (
	"flag"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/server"
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

	avatarDir := viper.GetString("frontend.public_dir") + viper.GetString("frontend.avatar_dir")
	if removeErr := os.RemoveAll(avatarDir); removeErr != nil {
		log.Fatal(err)
	}
	if mkdirErr := os.Mkdir(avatarDir, os.ModePerm); mkdirErr != nil {
		log.Fatal(err)
	}

	srv := &server.Server{
		IP:   viper.GetString("server.ip"),
		Port: uint(viper.GetInt("server.port")),
	}
	srv.Run()
}
