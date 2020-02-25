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

/***************** SessionStore **********************/

type SessionStore struct {
	sessions map[string]string
	mu       sync.Mutex // RWMutex в лекции?
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

/*func ReadUsers(r *http.Request) (old, new *User, err error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, nil, err
	}
	defer r.Body.Close()
	var data map[string]string
	err = json.Unmarshal(body, &data)
	_, hasOld := data["old_user"]
	_, hasNew := data["new_user"]
	if err != nil || !hasOld || !hasNew {
		return nil, nil, err
	}
	err = json.Unmarshal(data["old_user"], old)
	if err != nil {
		return nil, nil, err
	}
	err = json.Unmarshal(data["new_user"], new)
	return old, new, err
}*/

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

func (this *Handler) GetCookie(r *http.Request) (string, bool) {
	nick := ""
	authorized := false
	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		nick, authorized = this.sessionStore.GetSession(session.Value)
	}
	return nick, authorized
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
	SetHeaders(w, r)
	/*if this.HasCookie(r) {
		this.LogOut()
	}*/
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
		Name:    "session_id",
		Value:   SID,
		Path:    "/",
		Expires: time.Now().Add(10 * time.Hour),
		// SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
	SendMessage(w, http.StatusPermanentRedirect, Pair{"path", "/"})
}

func (this *Handler) LogOut(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w, r)
}

func (this *Handler) PostUser(w http.ResponseWriter, r *http.Request) {
	/*SetHeaders(w, r)
	if _, has := this.GetCookie(r); !has {
		SendMessage(w, http.StatusUnauthorized)
		return
	}
	oldUser, newUser, err := ReadUsers(r)
	if err != nil || oldUser.NickName == "" {
		SendMessage(w, http.StatusBadRequest)
		return
	}
	realUser, has := this.userStore.Get(oldUser.NickName)
	if !has {
		SendMessage(w, http.StatusNotFound)
		return
	}
	if oldUser.Password != "" && newUser.Password != "" && oldUser.Password == realUser.Password {
		realUser.Password = newUser.Password
	} else if oldUser.Password != "" {
		SendMessage(w, http.StatusPreconditionFailed)
		return
	}*/
}

func (this *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w, r)
	nickQuery, hasNick := r.URL.Query()["nickname"]
	if !hasNick || len(nickQuery) != 1  {
		SendMessage(w, http.StatusBadRequest)
		return
	}
	nickName := string(nickQuery[0])
	authNickName, hasCookie := this.GetCookie(r)
	isU := true
	if !hasCookie || nickName != authNickName {
		isU = false
	}
	realUser, has := this.userStore.Get(nickName)
	if !has {
		SendMessage(w, http.StatusNotFound)
		return
	}
	SendMessage(w, http.StatusOK, Pair{"user", realUser.GetInfo()}, Pair{"is_u", isU})
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
	router.HandleFunc("/logout", api.LogOut).Methods(http.MethodPost)
	router.HandleFunc("/profile", api.PostUser).Methods(http.MethodPost)
	router.HandleFunc("/profile", api.GetUser).Methods(http.MethodGet)

	log.Println("start")
	//wg := &WaitGroup{}
	http.ListenAndServe(":"+port, router)
}
