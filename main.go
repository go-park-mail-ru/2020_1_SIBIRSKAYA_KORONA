package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

/***************** UserStore **********************/

// TODO avatar
type User struct {
	Name     string `json:"name"`
	SurName  string `json:"surname"`
	NickName string `json:"nickname"`
	Password string `json:"password,omitempty"`
}

func (this *User) GetInfo() User {
	return User{
		Name:     this.Name,
		SurName:  this.SurName,
		NickName: this.NickName,
		Password: "",
	}
}

type UserStore struct {
	users map[string]*User
	mu    sync.Mutex // RWMutex в лекции?
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

func (this *UserStore) GetAll() map[string]*User {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.users
}

/***************** SessionStore **********************/

type SessionStore struct {
	sessions map[string]bool
	mu       sync.Mutex // RWMutex в лекции?
}

func CreateSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]bool),
		mu:       sync.Mutex{},
	}
}

func (this *SessionStore) AddSession(nickname string) string {
	this.mu.Lock()
	defer this.mu.Unlock()
	tmp := md5.Sum([]byte(nickname))
	SID := hex.EncodeToString(tmp[:])
	this.sessions[SID] = true
	return SID
}

func (this *SessionStore) HasSession(SID string) bool {
	this.mu.Lock()
	defer this.mu.Unlock()
	_, has := this.sessions[SID]
	return has
}

func (this *SessionStore) DeleteSession(SID string) {
	this.mu.Lock()
	defer this.mu.Unlock()
	delete(this.sessions, SID)
}

/***************** Transfer **********************/

type Pair struct {
	name string
	data interface{}
}

func SendMessage(w http.ResponseWriter, status uint, bodyData ...Pair) {
	msg := make(map[string]interface{})
	msg["status"] = status
	if len(bodyData) != 0 {
		bodyMap := make(map[string]interface{})
		for _, elem := range bodyData {
			bodyMap[elem.name] = elem.data
		}
		jbodyData, _ := json.Marshal(bodyMap)
	    msg["body"] = jbodyData
	}
	res, _ := json.Marshal(msg)
	io.WriteString(w, string(res))
}

func ReadUser(r *http.Request) (*User, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	var user User
	err = json.Unmarshal(body, &user)
	return &user, err
}

/***************** Handler **********************/

type Handler struct {
	userStore    *UserStore
	sessionStore *SessionStore
}

func (this *Handler) hasCookie(r *http.Request) bool {
	authorized := false
	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		authorized = this.sessionStore.HasSession(session.Value)
	}
	return authorized
}

func (this *Handler) Main(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if this.hasCookie(r) {
		w.Write([]byte("ты доска"))
	} else {
		w.Write([]byte("ты не доска"))
	}
}

func (this *Handler) Join(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	user, err := ReadUser(r)
	if err != nil || user.Name == "" || user.SurName == "" ||
		user.NickName == "" || user.Password == "" {
		SendMessage(w, http.StatusBadRequest)
		return
	}
	_, has := this.userStore.Get(user.NickName)
	if has {
		SendMessage(w, http.StatusConflict)
	} else {
		this.userStore.Add(user)
		SendMessage(w, http.StatusPermanentRedirect, Pair{"path", "/login"})
	}
}

func (this *Handler) LogIn(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if this.hasCookie(r) {
		//this.LogOut()
	}
	user, err := ReadUser(r)
	if err != nil || user.NickName == "" || user.Password == "" {
		SendMessage(w, http.StatusBadRequest)
		return
	}
	realUser, has := this.userStore.Get(user.NickName)
	if !has {
		SendMessage(w, http.StatusNotFound)
		return
	}
	if realUser.Password != user.Password {
		SendMessage(w, http.StatusPreconditionFailed)
		return
	}
	SID := this.sessionStore.AddSession(user.NickName)
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    SID,
		Path:     "/",
		Expires:  time.Now().Add(10 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
	SendMessage(w, http.StatusPermanentRedirect, Pair{"path", "/"})
}

func (this *Handler) LogOut(w http.ResponseWriter, r *http.Request) {
}

func (this *Handler) PostUser(w http.ResponseWriter, r *http.Request) {
	this.GetUser(w, r)
}

func (this *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if !this.hasCookie(r) {
		SendMessage(w, http.StatusUnauthorized)
		return
	}
	user, err := ReadUser(r)
	if err != nil || user.NickName == "" {
		SendMessage(w, http.StatusBadRequest)
		return
	}
	realUser, has := this.userStore.Get(user.NickName)
	if !has {
		SendMessage(w, http.StatusNotFound)
		return
	}
	SendMessage(w, http.StatusOK, Pair{"user", realUser.GetInfo()})
}

func main() {
	port := "8080"
	if len(os.Args) == 2 {
		port = os.Args[1]
	}
	router := mux.NewRouter()

	api := &Handler{
		userStore:    CreateUserStore(),
		sessionStore: CreateSessionStore(),
	}

	router.HandleFunc("/", api.Main)
	router.HandleFunc("/join", api.Join).Methods(http.MethodPost)
	router.HandleFunc("/login", api.LogIn).Methods(http.MethodPost)
	router.HandleFunc("/profile", api.PostUser).Methods(http.MethodPost)
	router.HandleFunc("/profile", api.GetUser).Methods(http.MethodGet)

	log.Println("start")
	//wg := &WaitGroup{}
	http.ListenAndServe(":"+port, router)
}
