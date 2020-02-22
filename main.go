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
)

/***************** UserStore **********************/

// TODO avatar
type User struct {
	// Id       uint   `json: "id"`
	Name     string `json:"name"`
	SurName  string `json:"surname"`
	// NickName string `json:"nickname"` // ключ в мапе
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

func (this *UserStore) Add(nickName string, user *User) {
	this.mu.Lock()
	defer this.mu.Unlock()
	//size := len(this.users)
	//user.Id = size
	this.users[nickName] = user // TODO: обработать существующие
	

}

func (this *UserStore) Get(nickName string) (*User) {
	this.mu.Lock()
	defer this.mu.Unlock()
	return this.users[nickName] // TODO: обработать не существующие
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
	http.Error(w, `{"join"}`, 400)
}


// http://127.0.0.1:8080/login?nick=tim&password=keklol
func (this *Handler) Login(w http.ResponseWriter, r *http.Request) {
	nickName := r.FormValue("nick")
	user := this.userStore.Get(nickName)
	if user.Password != r.FormValue("password") {
		http.Error(w, `bad pass`, 400)
		return
	}
	SID := this.sessionStore.AddSession(nickName)
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   SID,
		Path:    "/",
		Expires: time.Now().Add(10 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
	w.Write([]byte(SID))
}

func (this *Handler) User(w http.ResponseWriter, r *http.Request) {
	//nickName := mux.Vars(r)["nickname"]
	//user := this.store.Get(nickName);
	user := &User{
		Name: "Tim",
		SurName: "Razumov",
		Password: "keklol",
	}
	this.userStore.Add("tim", user)
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
	router.HandleFunc("/join", api.Join)
	router.HandleFunc("/login", api.Login)
	router.HandleFunc("/profile/{nickname}", api.User)//.Methods(http.MethodsGet)

	log.Println("start")
	//wg := &WaitGroup{}
	http.ListenAndServe(":" + port, router)
}
