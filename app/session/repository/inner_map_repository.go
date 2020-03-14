package repository

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"sync"

	//"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session"
)

type SessionStore struct {
	sessions map[string]string
	mu       sync.Mutex
}

func CreateSessionStore() session.Repository {
	return &SessionStore{
		sessions: make(map[string]string),
		mu:       sync.Mutex{},
	}
}

func (this *SessionStore) AddSession(nickname string) string {
	this.mu.Lock()
	defer this.mu.Unlock()
	tmp := md5.Sum([]byte(nickname))
	SID := hex.EncodeToString(tmp[:])
	this.sessions[SID] = nickname
	return SID
}

func (this *SessionStore) GetSession(SID string) (string, bool) {
	this.mu.Lock()
	defer this.mu.Unlock()
	val, has := this.sessions[SID]
	return val, has
}

func (this *SessionStore) DeleteSession(SID string) (err error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if _, has := this.sessions[SID]; has {
		delete(this.sessions, SID)
	} else {
		err = errors.New("no key")
	}
	return err
}
