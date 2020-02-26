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
	//"strings"
	"errors"
	"sync"
	"time"
)

/***************** UserStore **********************/

// TODO avatar
type User struct {
	Name     string `json:"name"`
	SurName  string `json:"surname"`
	NickName string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

func (this *User) GetInfo() User {
	return User{
		Name:     this.Name,
		SurName:  this.SurName,
		NickName: this.NickName,
		Email:    this.Email,
		Password: "",
	}
}

func (this *User) Empty() bool {
	return this.Name == "" || this.SurName == "" ||
		this.NickName == "" || /*this.Email == "" ||*/ this.Password == ""
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

/***************** SessionStore **********************/

type SessionStore struct {
	sessions map[string]string
	mu       sync.Mutex
}

func CreateSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]string),
		mu:       sync.Mutex{},
	}
}

func (this *SessionStore) AddSession(nickName string) string {
	this.mu.Lock()
	defer this.mu.Unlock()
	tmp := md5.Sum([]byte(nickName))
	SID := hex.EncodeToString(tmp[:])
	this.sessions[SID] = nickName
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
		msg["body"] = bodyMap
	}
	res, _ := json.Marshal(msg)
	// log.Println(string(res), err)
	io.WriteString(w, string(res))
}

// TODO: шаблонное чтение всех данных
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

func ReadChangeUser(r *http.Request) (*User, string, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, "", err
	}
	defer r.Body.Close()
	var user User
	err = json.Unmarshal(body, &user)
	// TODO
	var tmp map[string]interface{}
	err = json.Unmarshal(body, &tmp)
	return &user, tmp["old_password"].(string), err
}

/***************** Handler **********************/

type Handler struct {
	userStore    *UserStore
	sessionStore *SessionStore
}

func SetHeaders(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
}

func (this *Handler) SetCookie(w http.ResponseWriter, nickName string) {
	SID := this.sessionStore.AddSession(nickName)
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   SID,
		Path:    "/",
		Expires: time.Now().Add(10 * time.Hour),
		// SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}

func (this *Handler) GetCookie(r *http.Request) (string, bool) {
	nick := ""
	authorized := false
	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		nick, authorized = this.sessionStore.GetSession(session.Value)
	}
	return nick, authorized
}

func (this *Handler) DeleteCookie(w http.ResponseWriter, r *http.Request) error {
	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		err = this.sessionStore.DeleteSession(session.Value)
		session.Expires = time.Now().AddDate(0, 0, -2)
		http.SetCookie(w, session)
	}
	return err
}

func (this *Handler) Main(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w, r)
	if _, has := this.GetCookie(r); has {
		w.Write([]byte("ты доска"))
	} else {
		SendMessage(w, http.StatusPermanentRedirect, Pair{"path", "/login"})
	}
}

func (this *Handler) Join(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w, r)
	user, err := ReadUser(r)
	if err != nil || user.Empty() {
		SendMessage(w, http.StatusBadRequest)
		return
	}
	_, has := this.userStore.Get(user.NickName)
	if has {
		SendMessage(w, http.StatusConflict)
	} else {
		this.userStore.Add(user)
		this.SetCookie(w, user.NickName)
		SendMessage(w, http.StatusOK, Pair{"path", "/"})
	}
}

func (this *Handler) LogIn(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w, r)
	if _, hasCookie := this.GetCookie(r); hasCookie {
		SendMessage(w, http.StatusPermanentRedirect, Pair{"path", "/"})
		return
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
	this.SetCookie(w, user.NickName)
	SendMessage(w, http.StatusOK, Pair{"path", "/"})
}

func (this *Handler) LogOut(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w, r)
	if err := this.DeleteCookie(w, r); err == nil {
		SendMessage(w, http.StatusOK, Pair{"path", "/login"})
	} else {
		SendMessage(w, http.StatusSeeOther)
	}
}

func (this *Handler) PutUser(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w, r)
	newUser, oldPassword, err := ReadChangeUser(r)
	if err != nil || newUser.Empty() {
		SendMessage(w, http.StatusBadRequest)
		return
	}
	// TODO: через uid
	if nickSession, has := this.GetCookie(r); !has || nickSession != newUser.NickName {
		SendMessage(w, http.StatusForbidden)
		return
	}
	realUser, has := this.userStore.Get(newUser.NickName)
	if !has {
		SendMessage(w, http.StatusNotFound)
		return
	}
	if realUser.Password == oldPassword {
		realUser = newUser
	} else {
		SendMessage(w, http.StatusForbidden)
	}
	SendMessage(w, http.StatusOK, Pair{"user", realUser.GetInfo()})
}

func (this *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w, r)
	nickName := ""
	if nickQuery, hasNick := r.URL.Query()["nickname"]; hasNick {
		if len(nickQuery) != 1 {
			SendMessage(w, http.StatusBadRequest)
			return
		}
		nickName = string(nickQuery[0])
	} else {
		tmp, hasCookie := this.GetCookie(r)
		if !hasCookie {
			SendMessage(w, http.StatusSeeOther, Pair{"path", "/login"})
			return
		}
		nickName = tmp
	}
	realUser, has := this.userStore.Get(nickName)
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
	router.HandleFunc("/logout", api.LogOut).Methods(http.MethodDelete/*, http.MethodOptions*/)
	router.HandleFunc("/profile", api.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/profile", api.PutUser).Methods(http.MethodPut)

	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length,"+
			" Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers,"+
			" Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		return
	})

	log.Println("start")
	//wg := &WaitGroup{}
	http.ListenAndServe(":"+port, router)
}
