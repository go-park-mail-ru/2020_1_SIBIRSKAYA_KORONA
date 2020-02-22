package main

import (
	// "flag"
	// "encoding/json"
	"github.com/gorilla/mux"
	"log"
	"sync"
	"net/http"
)

// TODO avatar
type User struct {
	Name     string `json:"name"`
	SurName  string `json:"surname"`
	// NickName string `json:"nickname"` ключ в мапе
	Password string `json:"password"`
}

type UserStore struct {
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


type UserHandler struct {
	store *UserStore
}

func (this *UserHandler) Main(w http.ResponseWriter, r *http.Request) {
	http.Error(w, `{"main"}`, 400)
}

func (this *UserHandler) Join(w http.ResponseWriter, r *http.Request) {
	http.Error(w, `{"join"}`, 400)
}

func (this *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	http.Error(w, `{"login"}`, 400)
}

func (this *UserHandler) User(w http.ResponseWriter, r *http.Request) {
	http.Error(w, `{"password"}`, 400)
}

/*func Main(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	    jMsg, err := json.Marshal("aa: 1")
		if err != nil {
			log.Println(err)
			return
		}
		r.Write(jMsg)
	//http.Error(w, `{"error":"bad id"}`, 400)
}*/

func main() {
	router := mux.NewRouter()

	api := &UserHandler{
		store: CreateUserStore(),
	}

	router.HandleFunc("/", api.Main)
	router.HandleFunc("/join", api.Join)
	router.HandleFunc("/login", api.Login)
	router.HandleFunc("/profile/{nickname}", api.User)//.Methods(http.MethodsGet)

	log.Println("start serving :8080")
	http.ListenAndServe(":8080", router)
}
