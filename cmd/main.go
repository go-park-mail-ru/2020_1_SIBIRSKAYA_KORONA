package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/models"
	sessionRepository "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/session/repository"
	userRepository "github.com/go-park-mail-ru/2020_1_SIBIRSKAYA_KORONA/app/user"

	"github.com/gorilla/mux"
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
	io.WriteString(w, string(res))
}

// TODO: шаблонное чтение всех данных
func ReadUser(r *http.Request) (*models.User, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	var user models.User
	err = json.Unmarshal(body, &user)
	return &user, err
}

/***************** Handler **********************/

type Handler struct {
	userStore    *userRepository.UseCase
	sessionStore *sessionRepository.SessionStore
}

func SetHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", allowOriginUrl)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func (this *Handler) SetCookie(w http.ResponseWriter, nickname string) {
	SID := this.sessionStore.AddSession(nickname)
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
	if err != nil || user.Nickname == "" || user.Password == "" {
		SendMessage(w, http.StatusBadRequest)
		return
	}
	_, has := this.userStore.GetUser(user.Nickname)
	if has {
		SendMessage(w, http.StatusConflict)
	} else {
		this.userStore.AddUser(user)
		user.PathToAvatar = defaultUserImgPath
		this.SetCookie(w, user.Nickname)
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
	if err != nil || user.Nickname == "" || user.Password == "" {
		SendMessage(w, http.StatusBadRequest)
		return
	}
	realUser, has := this.userStore.GetUser(user.Nickname)
	if !has {
		SendMessage(w, http.StatusNotFound)
		return
	}
	if realUser.Password != user.Password {
		SendMessage(w, http.StatusPreconditionFailed)
		return
	}
	this.SetCookie(w, user.Nickname)
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
	nickname, has := this.GetCookie(r)
	if !has {
		SendMessage(w, http.StatusUnauthorized)
		return
	}
	realUser, has := this.userStore.GetUserAll(nickname)
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
	if newUser.Surname != "" {
		realUser.Surname = newUser.Surname
	}
	if newUser.Email != "" {
		realUser.Email = newUser.Email
	}
	pathImg, err := UploadAvatarToLocalStorage(r, newUser.Nickname)
	if err == nil {
		log.Println("change photo")
		realUser.PathToAvatar = pathImg
	} else {
		log.Println("change photo error", err)
	}
	SendMessage(w, http.StatusOK, Pair{"user", realUser})
}

func (this *Handler) GetSettingsUser(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w, r)
	nickname, hasCookie := this.GetCookie(r)
	if !hasCookie {
		SendMessage(w, http.StatusSeeOther)
		return
	}
	realUser, has := this.userStore.GetUser(nickname)
	if !has {
		SendMessage(w, http.StatusNotFound)
		return
	}
	SendMessage(w, http.StatusOK, Pair{"user", realUser})
}

func (this *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w, r)
	nickname, has := mux.Vars(r)["nickname"]
	if !has {
		SendMessage(w, http.StatusBadRequest)
		return
	}
	realUser, has := this.userStore.GetUser(nickname)
	if !has {
		SendMessage(w, http.StatusNotFound)
		return
	}
	SendMessage(w, http.StatusOK, Pair{"user", realUser})
}

/******************** ФОТО ***************/

func LocalStorageInit() error {
	os.RemoveAll(localStorage)
	return os.Mkdir(localStorage, os.ModePerm)
}

func ReadUserForUpdate(r *http.Request) (*models.User, string) {
	var user models.User
	user.Name = r.FormValue("newName")
	user.Surname = r.FormValue("newSurname")
	user.Nickname = r.FormValue("newNickname")
	user.Email = r.FormValue("newEmail")
	oldPassword := r.FormValue("oldPassword")
	user.Password = r.FormValue("newPassword")
	return &user, oldPassword
}

func UploadAvatarToLocalStorage(r *http.Request, nickname string) (string, error) {
	avatarSrc, _, err := r.FormFile("avatar")
	avatarExtension := r.FormValue("avatarExtension")
	if err != nil {
		return "", err
	}
	defer avatarSrc.Close()
	avatarFileName := fmt.Sprintf("%s.%s", nickname, avatarExtension)
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
		userStore:    userRepository.CreateUserStore(),
		sessionStore: sessionRepository.CreateSessionStore(),
	}

	router.HandleFunc("/settings", api.Join).Methods(http.MethodPost)           // создать аккаунт
	router.HandleFunc("/settings", api.GetSettingsUser).Methods(http.MethodGet) // получ все настройки
	router.HandleFunc("/settings", api.PutUser).Methods(http.MethodPut)         // изменить настройки юзера

	router.HandleFunc("/profile/{nickname}", api.GetUser).Methods(http.MethodGet) // получ не приватные настройки

	router.HandleFunc("/session", api.LogIn).Methods(http.MethodPost)     // login
	router.HandleFunc("/session", api.CheckLogIn).Methods(http.MethodGet) // залогинен или нет
	router.HandleFunc("/session", api.LogOut).Methods(http.MethodDelete)  // logout

	router.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowOriginUrl)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	})

	log.Println("start")
	//wg := &WaitGroup{}
	http.ListenAndServe(":"+port, router)
}
