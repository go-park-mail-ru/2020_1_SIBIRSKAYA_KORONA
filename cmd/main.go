package main

import (
	"flag"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/server"
	"github.com/spf13/viper"
)

var opts struct {
	configPath string
}

// TODO: вынести в конфиг
//const frontendAbsolutePublicDir = "/home/gavroman/tp/2sem/tp_front/2020_1_SIBIRSKAYA_KORONA/public"
const frontendAbsolutePublicDir = "/home/timofey/2020_1_SIBIRSKAYA_KORONA/public"

const frontendUrl = "http://localhost:5757"

//const frontendAbsolutePublicDir = "/home/ubuntu/frontend/public" // (or absolute path to public folder in frontend)
//const frontendUrl = "http://89.208.197.150:5757"

const frontendAvatarStorage = frontendUrl + "/img/avatar"
const defaultUserImgPath = frontendUrl + "/img/default_avatar.png"
const localStorage = frontendAbsolutePublicDir + "/img/avatar"
const allowOriginUrl = frontendUrl

func init() {
	flag.StringVar(&opts.configPath, "c", "", "path to configuration file")
	flag.StringVar(&opts.configPath, "config", "", "path to configuration file")
	flag.Parse()

	// к этому моменту уже создан инстанс этого объекта в функции init() соответсвующего пакета
	// в дальнейшем можно обращаться к конфигу через его методы в остальных пакетах
	viper.SetConfigFile(opts.configPath)
	err := viper.ReadInConfig()
	if err != nil {
		// в main() можно даже не переходить
		panic(err)
	}
}

func main() {
	flag.Parse()

	srv := &server.Server{
		IP:   viper.GetString("server.ip"),
		Port: viper.GetInt("server.port"),
	}
	srv.Run()
}
