package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"bytes"
)

type TestCase struct {
	requestBody  map[string]interface{}
	responseBody map[string]interface{}
	httpMethod   string
	apiMethod    func(http.ResponseWriter, *http.Request)
}

func TestJoinHandler(t *testing.T) {
	// в начале каждого тестового сценария принудительно сбрасываем хранилище
	apiHandler := Handler{
		userStore:    CreateUserStore(),
		sessionStore: CreateSessionStore(),
	}

	joinUrl := "http://localhost:8080/join"

	cases := []TestCase{
		TestCase{
			requestBody: map[string]interface{}{
				"name":     "Антон",
				"surname":  "Гофер",
				"nickname": "Love",
				"password": "lovelove",
			},

			responseBody: map[string]interface{}{
				"body": map[string]interface{}{
					"path": "/login",
				},
				"status": 308,
			},

			httpMethod: "POST",
			apiMethod:  apiHandler.Join,
		},

		TestCase{
			requestBody: map[string]interface{}{
				"name":     "Антон",
				"surname":  "Гофер",
				"nickname": "Love",
				"password": "lovelove",
			},

			// пользователь уже существует, соответствующий код ответа
			responseBody: map[string]interface{}{
				"body": map[string]interface{}{
					"path": "/login",
				},
				"status": 409,
			},

			httpMethod: "POST",
			apiMethod:  apiHandler.Join,
		},

		TestCase{
			//json body не соответствует модели
			requestBody: map[string]interface{}{
				"field": "Тимофей",
				"text":  "Гофер",
			},

			responseBody: map[string]interface{}{
				"body": map[string]interface{}{
					"path": "/login",
				},
				"status": 400,
			},

			httpMethod: "POST",
			apiMethod:  apiHandler.Join,
		},
	}

	for num, test := range cases {
		reqBody, err := json.Marshal(test.requestBody)
		if err != nil {
			t.Error(err)
		}

		resBody, err := json.Marshal(test.responseBody)
		if err != nil {
			t.Error(err)
		}

		request, err := http.NewRequest(test.httpMethod, joinUrl, bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		test.apiMethod(w, request)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		bodyStr := string(body)
		if bodyStr != string(resBody) {
			t.Errorf("[%d] wrong response by JoinHandler: got %+v, expected %+v",
				num, bodyStr, string(resBody))
		}

	}
}

func TestLoginHandler(t *testing.T) {
	apiHandler := Handler{
		userStore:    CreateUserStore(),
		sessionStore: CreateSessionStore(),
	}

	loginUrl := "http://localhost:8080/login"

	cases := []TestCase{
		TestCase{
			requestBody: map[string]interface{}{
				"nickname": "test",
				"password": "testtest",
			},

			responseBody: map[string]interface{}{
				"body": map[string]interface{}{
					"path": "/login",
				},
				"status": 404,
			},

			httpMethod: "POST",
			apiMethod:  apiHandler.LogIn,
		},

		TestCase{
			requestBody: map[string]interface{}{
				"name":     "Антон",
				"surname":  "Гофер",
				"nickname": "Love",
				"password": "lovelove",
			},

			responseBody: map[string]interface{}{
				"status": 308,
			},

			httpMethod: "POST",
			apiMethod:  apiHandler.Join,
		},

		TestCase{
			requestBody: map[string]interface{}{
				"nickname": "Love",
				"password": "lovelove",
			},

			responseBody: map[string]interface{}{
				"body": map[string]interface{}{
					"path": "/login",
				},
				"status": 308,
			},

			httpMethod: "POST",
			apiMethod:  apiHandler.LogIn,
		},
	}

	for num, test := range cases {
		reqBody, err := json.Marshal(test.requestBody)
		if err != nil {
			t.Error(err)
		}
		resBody, err := json.Marshal(test.responseBody)
		if err != nil {
			t.Error(err)
		}

		request, err := http.NewRequest(test.httpMethod, loginUrl, bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		test.apiMethod(w, request)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		bodyStr := string(body)
		if bodyStr != string(resBody) {
			t.Errorf("[%d] wrong response by LogInHandler: got %+v, expected %+v",
				num, bodyStr, string(resBody))
		}

	}
}

// реализуем сложный сценарий, в котором мы
// * регистрируем пользователя
// * пытаемся авторизоваться (получив куку и сохранив её у себя)
// * пытаемся получить информацию о пользователе с помощью выданной нам куки
func TestGetUserHandler(t *testing.T) {
	apiHandler := Handler{
		userStore:    CreateUserStore(),
		sessionStore: CreateSessionStore(),
	}

	//getUserUrl := "http://localhost:8080/login" // ??

	/***************** Регистрация **********************/

	joinReqBodyMap := map[string]interface{}{
		"name":     "Антон",
		"surname":  "Гофер",
		"nickname": "Love",
		"password": "lovelove",
	}

	joinRespBodyMap := map[string]interface{}{
		"body": map[string]interface{}{
			"path": "/login",
		},
		"status": 308,
	}

	joinReqBody, err := json.Marshal(joinReqBodyMap)
	if err != nil {
		t.Error(err)
	}

	joinRespBodyExpected, err := json.Marshal(joinRespBodyMap)
	if err != nil {
		t.Error(err)
	}

	joinRequest, err := http.NewRequest("POST", "fakeurl", bytes.NewBuffer(joinReqBody))
	w := httptest.NewRecorder()

	apiHandler.Join(w, joinRequest)

	joinResponse := w.Result()
	joinResponseBody, _ := ioutil.ReadAll(joinResponse.Body)
	defer joinResponse.Body.Close()

	if string(joinResponseBody) != string(joinRespBodyExpected) {
		t.Errorf("GetUser scenario wrong response from joinHandler: got %+v, expected %+v",
			string(joinResponseBody), string(joinRespBodyExpected))
	}

	/***************** Регистрация **********************/

	/***************** Авторизация **********************/

	loginReqBodyMap := map[string]interface{}{
		"nickname": "Love",
		"password": "lovelove",
	}

	loginRespBodyMap := map[string]interface{}{
		"body": map[string]interface{}{
			"path": "/login",
		},
		"status": 308,
	}

	loginReqBody, err := json.Marshal(loginReqBodyMap)
	if err != nil {
		t.Error(err)
	}

	loginRespBodyExpected, err := json.Marshal(loginRespBodyMap)
	if err != nil {
		t.Error(err)
	}

	loginRequest, err := http.NewRequest("POST", "fakeurl", bytes.NewBuffer(loginReqBody))
	wLogin := httptest.NewRecorder()

	apiHandler.LogIn(wLogin, loginRequest)

	loginResponse := wLogin.Result()

	ourCookies := loginResponse.Cookies()

	var ourCookie *http.Cookie = nil

	for idCookie, _ := range ourCookies {
		if ourCookies[idCookie].Name == "session_id" {
			ourCookie = ourCookies[idCookie]
			break
		}
	}

	if ourCookie == nil {
		log.Fatal("Cant find session_id cookie")
	}

	loginResponseBody, _ := ioutil.ReadAll(loginResponse.Body)
	defer loginResponse.Body.Close()

	if string(loginResponseBody) != string(loginRespBodyExpected) {
		t.Errorf("GetUser scenario wrong response from loginHandler: got %+v, expected %+v",
			string(loginResponseBody), string(loginRespBodyExpected))
	}

	/***************** Авторизация **********************/

	/***************** Получаем данные через куку **********************/

	getReqBodyMap := map[string]interface{}{
		"nickname": "Love",
	}

	getRespBodyMap := map[string]interface{}{
		"body": map[string]interface{}{
			"user": map[string]interface{}{
				"name":     "Антон",
				"surname":  "Гофер",
				"nickname": "Love",
			},
		},
		"status": 200,
	}

	getReqBody, err := json.Marshal(getReqBodyMap)
	if err != nil {
		t.Error(err)
	}

	getRespBodyExpected, err := json.Marshal(getRespBodyMap)
	if err != nil {
		t.Error(err)
	}

	getUserRequest, err := http.NewRequest("GET", "fakeurl", bytes.NewBuffer(getReqBody))
	wGet := httptest.NewRecorder()

	getUserRequest.AddCookie(ourCookie)

	apiHandler.GetUser(wGet, getUserRequest)

	getResponse := wGet.Result()

	getResponseBody, _ := ioutil.ReadAll(getResponse.Body)
	defer getResponse.Body.Close()

	if string(getResponseBody) != string(getRespBodyExpected) {
		t.Errorf("GetUser scenario wrong response from getHandler: got %+v, expected %+v",
			string(getResponseBody), string(getRespBodyExpected))
	}

	/***************** Получаем данные через куку **********************/

}
