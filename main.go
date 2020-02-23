package main

import (
	"os"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"sync"
	"net/http"
	"time"
	"io"
	"io/ioutil"
)

/***************** UserStore **********************/

// TODO avatar
type User struct {
	// Id       uint   `json: "id"`
	Name     string `json:"name"`
	SurName  string `json:"surname"`
	NickName string `json:"nickname"` // ключ в мапе
	Password string `json:"password"`
}

type UserStore struct {
	sessions map[string]string
	users map[string]*User
	mu    sync.Mutex    // RWMutex в лекции?
}

func CreateUserStore() *UserStore {
	return &UserStore{
		users: make(map[string]*User),
		mu:    sync.Mutex{},
	}
}

func (this *UserStore) Add(user *User) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.users[user.NickName] = user
}

func (this *UserStore) Get(nickName string) (*User, bool) {
	this.mu.Lock()
	defer this.mu.Unlock()
	user, has := this.users[nickName]
	return user, has
}

func (this *UserStore) GetAll() (map[string]*User) {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.users
}

/***************** SessionStore **********************/

type SessionStore struct {
	sessions map[string]string
	mu    sync.Mutex    // RWMutex в лекции?
}

func CreateSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]string),
		mu:    sync.Mutex{},
	}
}

func (this *SessionStore) AddSession(login string) string {
	this.mu.Lock()
	defer this.mu.Unlock()
	tmp := md5.Sum([]byte(login))
	SID :=	hex.EncodeToString(tmp[:])
	this.sessions[SID] = login
	return SID
}

func (this *SessionStore) HasSession(SID string) bool {
	this.mu.Lock()
	defer this.mu.Unlock()
	_, has := this.sessions[SID]
	return has
}

func (this *SessionStore) DeleteSession() {
}

/***************** Handler **********************/

type Handler struct {
	userStore *UserStore
	sessionStore *SessionStore
}

func (this *Handler) Main(w http.ResponseWriter, r *http.Request) {
	authorized := false
	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		authorized = this.sessionStore.HasSession(session.Value)
	}

	if authorized {
		w.Write([]byte("autrorized"))
	} else {
		w.Write([]byte("not autrorized"))
	}
}

func (this *Handler) Join(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		io.WriteString(w, `{"status":400}`)
		return
	}
	defer r.Body.Close()
	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		io.WriteString(w, `{"status":400}`)
		return
	}
	_, has := this.userStore.Get(user.NickName)
	if has {
		io.WriteString(w, `{"status":409}`)
	} else {
		this.userStore.Add(&user)
		io.WriteString(w, `{"status":308}`)
			                // TODO:`{"status": 301, "body": {"path": "/login"}}`
	}
}

func (this *Handler) Login(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		io.WriteString(w, `{"status":400}`)
		return
	}
	defer r.Body.Close()
	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		io.WriteString(w, `{"status":400}`)
		return
	}
	realUser, has := this.userStore.Get(user.NickName)
	if !has {
		io.WriteString(w, `{"status":404}`)
		return
	}
	if realUser.Password != user.Password {
		io.WriteString(w, `{"status":412}`)
		return
	}
	SID := this.sessionStore.AddSession(user.NickName)
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   SID,
		Path:    "/",
		Expires: time.Now().Add(10 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
	io.WriteString(w, `{"status":308}`)
}

func (this *Handler) User(w http.ResponseWriter, r *http.Request) {
	user := &User{
		Name: "Tim",
		SurName: "Razumov",
		Password: "keklol",
		NickName: "tim",
	}
	this.userStore.Add(user)
	jsonData, _ := json.Marshal(user)	
	w.Write(jsonData)
}

func main() {
	port := "8080"
	if len(os.Args) == 2 {
		port = os.Args[1]
	}
	router := mux.NewRouter()

	api := &Handler{
		userStore: CreateUserStore(),
		sessionStore: CreateSessionStore(),
	}

	router.HandleFunc("/", api.Main)
	router.HandleFunc("/join", api.Join).Methods(http.MethodPost)
	router.HandleFunc("/login", api.Login).Methods(http.MethodPost)
	router.HandleFunc("/profile/{nickname}", api.User)//.Methods(http.MethodsGet)

	log.Println("start")
	//wg := &WaitGroup{}
	http.ListenAndServe(":" + port, router)
}
