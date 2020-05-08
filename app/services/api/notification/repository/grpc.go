package repository

import (

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models/proto"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/api/user"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/config"
)

type NotificationStore struct {
}

func CreateRepository(clt proto.UserClient, Config_ *config.UserConfigController) user.Repository {
	return &UserStore{
		clt:    clt,
		Config: Config_,
	}
}