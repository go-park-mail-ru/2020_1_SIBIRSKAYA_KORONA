package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
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

//const frontendAbsolutePublicDir = "/home/gavroman/tp/2sem/tp_front/2020_1_SIBIRSKAYA_KORONA/public"
const frontendAbsolutePublicDir = "/home/timofey/2020_1_SIBIRSKAYA_KORONA/public"

const frontendUrl = "http://localhost:5757"

//const frontendAbsolutePublicDir = "/home/ubuntu/frontend/public" // (or absolute path to public folder in frontend)
//const frontendUrl = "http://89.208.197.150:5757"

const frontendAvatarStorage = frontendUrl + "/img/avatar"
const defaultUserImgPath = frontendUrl + "/img/default_avatar.png"
const localStorage = frontendAbsolutePublicDir + "/img/avatar"
const allowOriginUrl = frontendUrl

/***************** UserStore **********************/

// TODO avatar
type User struct {
	Name         string `json:"name"`
	SurName      string `json:"surname"`
	NickName     string `json:"nickname"`
	Email        string `json:"email"`
	PathToAvatar string `json:"avatar"`
	Password     string `json:"password,omitempty"`
}

func (this *User) GetInfo() User {
	return User{
		Name:         this.Name,
		SurName:      this.SurName,
		NickName:     this.NickName,
		Email:        this.Email,
		PathToAvatar: this.PathToAvatar,
		Password:     "",
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

/***************** Handler **********************/

type Handler struct {
	userStore    *UserStore
	sessionStore *SessionStore
}

func SetHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", allowOriginUrl)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
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

func (this *Handler) Join(w http.ResponseWriter, r *http.Request) {
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
	_, has := this.userStore.Get(user.NickName)
	if has {
		SendMessage(w, http.StatusConflict)
	} else {
		this.userStore.Add(user)
		user.PathToAvatar = defaultUserImgPath
		this.SetCookie(w, user.NickName)
		SendMessage(w, http.StatusOK)
	}
}

func (this *Handler) LogIn(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w, r)
	if _, hasCookie := this.GetCookie(r); hasCookie {
		SendMessage(w, http.StatusPermanentRedirect)
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
		SendMessage(w, http.StatusOK)
	} else {
		SendMessage(w, http.StatusSeeOther)
	}
}

func (this *Handler) CheckLogIn(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w, r)
	if _, hasCookie := this.GetCookie(r); hasCookie {
		SendMessage(w, http.StatusOK)
	} else {
		SendMessage(w, http.StatusUnauthorized)
	} 
}

func (this *Handler) PutUser(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w, r)
	nickName, has := this.GetCookie(r);
	if !has {
		SendMessage(w, http.StatusUnauthorized)
		return
	}
	realUser, has := this.userStore.Get(nickName)
	if !has {
		SendMessage(w, http.StatusNotFound)
		return
	}
	newUser, oldPassword := ReadUserForUpdate(r)
	log.Println("old", oldPassword)
	log.Println("real", realUser.Password)
	log.Println("new", newUser.Password)
	if oldPassword != "" && newUser.Password != "" {
		if realUser.Password == oldPassword {
			realUser.Password = newUser.Password
			SendMessage(w, http.StatusOK)
		} else {
			SendMessage(w, http.StatusPreconditionFailed)
		}
		return
	}
	if newUser.Name != "" {
		realUser.Name = newUser.Name
	}
	if newUser.SurName != "" {
		realUser.SurName = newUser.SurName
	}
	if newUser.Email != "" {
		realUser.Email = newUser.Email
	}
	pathImg, err := UploadAvatarToLocalStorage(r, newUser.NickName)
	if err == nil {
		log.Println("change photo")
		realUser.PathToAvatar = pathImg
	} else {
		log.Println("change photo error", err)
	}
	SendMessage(w, http.StatusOK, Pair{"user", realUser.GetInfo()})
}

func (this *Handler) GetSettingsUser(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w, r)
	nickName, hasCookie := this.GetCookie(r)
	if !hasCookie {
		SendMessage(w, http.StatusSeeOther)
		return
	}
	realUser, has := this.userStore.Get(nickName)
	if !has {
		SendMessage(w, http.StatusNotFound)
		return
	}
	SendMessage(w, http.StatusOK, Pair{"user", realUser.GetInfo()})
}

func (this *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w, r)
	nickName, has := mux.Vars(r)["nickname"]
	if !has {
		SendMessage(w, http.StatusBadRequest)
		return
	}
	realUser, has := this.userStore.Get(nickName)
	if !has {
		SendMessage(w, http.StatusNotFound)
		return
	}
	SendMessage(w, http.StatusOK, Pair{"user", realUser.GetInfo()})
}

/******************** ФОТО ***************/

func LocalStorageInit() error {
	os.RemoveAll(localStorage)
	return os.Mkdir(localStorage, os.ModePerm)
}

func ReadUserForUpdate(r *http.Request) (*User, string) {
	var user User
	user.Name = r.FormValue("newName")
	user.SurName = r.FormValue("newSurname")
	user.NickName = r.FormValue("newNickname")
	user.Email = r.FormValue("newEmail")
	oldPassword := r.FormValue("oldPassword")
	user.Password = r.FormValue("newPassword")
	return &user, oldPassword
}

func UploadAvatarToLocalStorage(r *http.Request, nickName string) (string, error) {
	avatarSrc, _, err := r.FormFile("avatar")
	avatarExtension := r.FormValue("avatarExtension")
	if err != nil {
		return "", err
	}
	defer avatarSrc.Close()
	avatarFileName := fmt.Sprintf("%s.%s", nickName, avatarExtension)
	avatarPath := fmt.Sprintf("%s/%s", localStorage, avatarFileName)
	avatarDst, err := os.Create(avatarPath)
	if err != nil {
		return "", err
	}
	defer avatarDst.Close()
	_, err = io.Copy(avatarDst, avatarSrc)

	frontEndAvatarUrl := fmt.Sprintf("%s/%s", frontendAvatarStorage, avatarFileName)
	return frontEndAvatarUrl, err
}

/*********************** ФОТО ********************/

func main() {
	err := LocalStorageInit()
	if err != nil {
		log.Fatal("Local storage init failed: ", err)
	}

	port := "8080"
	if len(os.Args) == 2 {
		port = os.Args[1]
	}
	router := mux.NewRouter()

	api := &Handler{
		userStore:    CreateUserStore(),
		sessionStore: CreateSessionStore(),
	}

	router.HandleFunc("/settings", api.Join).Methods(http.MethodPost) // создать аккаунт
	router.HandleFunc("/settings", api.GetSettingsUser).Methods(http.MethodGet) // получ все настройки
	router.HandleFunc("/settings", api.PutUser).Methods(http.MethodPut) // изменить настройки юзера

	router.HandleFunc("/profile/{nickname}", api.GetUser).Methods(http.MethodGet) // получ не приватные настройки

    router.HandleFunc("/session", api.LogIn).Methods(http.MethodPost) // login
    router.HandleFunc("/session", api.CheckLogIn).Methods(http.MethodGet) // залогинен или нет
    router.HandleFunc("/session", api.LogOut).Methods(http.MethodDelete) // logout

	router.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowOriginUrl)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	})

	log.Println("start")
	//wg := &WaitGroup{}
	http.ListenAndServe(":"+port, router)
}
