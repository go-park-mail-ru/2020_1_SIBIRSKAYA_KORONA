package repository

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/session"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/logger"
)

type SessionStore struct {
	DB *memcache.Client
}

func CreateRepository(db *memcache.Client) session.Repository {
	return &SessionStore{db}
}

func (sessionStore *SessionStore) Create(ses models.Session) error {
	err := sessionStore.DB.Set(&memcache.Item{
		Key:        ses.SID,
		Value:      []byte(fmt.Sprintf("%d", ses.ID)),
		Expiration: ses.Expires,
	})
	if err != nil {
		logger.Error(err)
		return errors.ErrDbBadOperation
	}
	return nil
}

func (sessionStore *SessionStore) Get(sid string) (uint, error) {
	idByte, err := sessionStore.DB.Get(sid)
	if err != nil {
		logger.Error(err)
		return 0, errors.ErrDbBadOperation
	}
	var id uint
	_, err = fmt.Sscan(string(idByte.Value), &id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (sessionStore *SessionStore) Delete(sid string) error {
	err := sessionStore.DB.Delete(sid)
	if err != nil {
		logger.Error(err)
		return errors.ErrDbBadOperation
	}
	return nil
}
