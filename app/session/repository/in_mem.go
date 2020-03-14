package repository

import (
	"errors"
	"sync"

	//"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session"
)

type SessionStore struct {
	sessions map[string]uint
	mu       sync.Mutex
}

func CreateSessionStore() session.Repository {
	return &SessionStore{
		sessions: make(map[string]uint),
		mu:       sync.Mutex{},
	}
}

func (sessionStore *SessionStore) Create(id uint) (string, error) {
	sessionStore.mu.Lock()
	defer sessionStore.mu.Unlock()
	// TODO: норм хэширование
	sid := string(id)
	sessionStore.sessions[sid] = id
	return sid, nil
}

func (sessionStore *SessionStore) Has(sid string) bool {
	sessionStore.mu.Lock()
	defer sessionStore.mu.Unlock()
	_, has := sessionStore.sessions[sid]
	return has
}

func (sessionStore *SessionStore) Delete(sid string) (err error) {
	sessionStore.mu.Lock()
	defer sessionStore.mu.Unlock()
	if _, has := sessionStore.sessions[sid]; has {
		delete(sessionStore.sessions, sid)
	} else {
		err = errors.New("no key")
	}
	return err
}