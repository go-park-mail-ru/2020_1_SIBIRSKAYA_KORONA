package main

import (
	"flag"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
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

	_, avatarErr := os.Stat(avatarDir)
	if os.IsNotExist(avatarErr) {
		errDir := os.MkdirAll(avatarDir, os.ModePerm)
		if errDir != nil {
			log.Fatal(models.ErrBadAvatarUpload, errDir)
		}
	}
	log.Println("Avatar static storage up!")

	srv := &server.Server{
		IP:   viper.GetString("server.ip"),
		Port: uint(viper.GetInt("server.port")),
	}
	srv.Run()
}
