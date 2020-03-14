package main

import (
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/server"
)

//const frontendAbsolutePublicDir = "/home/gavroman/tp/2sem/tp_front/2020_1_SIBIRSKAYA_KORONA/public"
const frontendAbsolutePublicDir = "/home/timofey/2020_1_SIBIRSKAYA_KORONA/public"

const frontendUrl = "http://localhost:5757"

//const frontendAbsolutePublicDir = "/home/ubuntu/frontend/public" // (or absolute path to public folder in frontend)
//const frontendUrl = "http://89.208.197.150:5757"

const frontendAvatarStorage = frontendUrl + "/img/avatar"
const defaultUserImgPath = frontendUrl + "/img/default_avatar.png"
const localStorage = frontendAbsolutePublicDir + "/img/avatar"
const allowOriginUrl = frontendUrl

func main() {
	server.Run()
}
