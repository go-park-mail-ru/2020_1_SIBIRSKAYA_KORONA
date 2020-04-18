package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/services/session"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/pkg/errors"

	"github.com/bradfitz/gomemcache/memcache"
)

type SessionStore struct {
	DB *memcache.Client
}

func CreateRepository(db *memcache.Client) session.Repository {
	return &SessionStore{db}
}

func (sessionStore *SessionStore) Create(ses models.Session) error {
	return sessionStore.DB.Set(&memcache.Item{
		Key:        ses.SID,
		Value:      []byte(fmt.Sprintf("%d", ses.ID)),
		Expiration: int32(ses.Expires),
	})
}

func (sessionStore *SessionStore) Get(sid string) (uint, error) {
	idByte, err := sessionStore.DB.Get(sid)
	if err != nil {
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
	return sessionStore.DB.Delete(sid)
}
