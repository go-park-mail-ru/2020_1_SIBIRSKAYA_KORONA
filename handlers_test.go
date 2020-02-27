package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestCase struct {
	requestBody  map[string]interface{}
	responseBody map[string]interface{}
	httpMethod   string
	apiMethod    func(http.ResponseWriter, *http.Request)
}

func TestJoinHandler(t *testing.T) {
	t.Parallel()
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
				"email":    "aaa@mail.ru",
				"password": "lovelove",
			},

			responseBody: map[string]interface{}{
				"body": map[string]interface{}{
					"path": "/",
				},
				"status": 200,
			},

			httpMethod: "POST",
			apiMethod:  apiHandler.Join,
		},

		TestCase{
			requestBody: map[string]interface{}{
				"name":     "Антон",
				"surname":  "Гофер",
				"nickname": "Love",
				"email":    "aaa@mail.ru",
				"password": "lovelove",
			},

			// пользователь уже существует, соответствующий код ответа
			responseBody: map[string]interface{}{
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
	t.Parallel()
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
				"email":    "aaa@mail.ru",
				"password": "lovelove",
			},

			responseBody: map[string]interface{}{
				"body": map[string]interface{}{
					"path": "/",
				},
				"status": 200,
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
					"path": "/",
				},
				"status": 200,
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

// * регистрируем пользователя (получив куку и сохранив её у себя)
// * пытаемся получить информацию о пользователе с помощью выданной нам куки
func TestGetUserHandler(t *testing.T) {
	t.Parallel()
	apiHandler := Handler{
		userStore:    CreateUserStore(),
		sessionStore: CreateSessionStore(),
	}

	user := User{
		Name:         "Антон",
		SurName:      "Гофер",
		NickName:     "Love",
		Email:        "aaa@mail.ru",
		Password:     "lovelove",
		PathToAvatar: defaultUserImgPath,
	}

	/***************** Регистрация **********************/

	joinUrl := "http://localhost:8080/join"

	joinRespBodyMap := map[string]interface{}{
		"body": map[string]interface{}{
			"path": "/",
		},
		"status": 200,
	}

	joinReqBody, err := json.Marshal(user)
	if err != nil {
		t.Error(err)
	}

	joinRespBodyExpected, err := json.Marshal(joinRespBodyMap)
	if err != nil {
		t.Error(err)
	}

	joinRequest, err := http.NewRequest("POST", joinUrl, bytes.NewBuffer(joinReqBody))
	w := httptest.NewRecorder()
	apiHandler.Join(w, joinRequest)
	joinResponse := w.Result()

	// check cookie
	ourCookies := joinResponse.Cookies()
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
	joinResponseBody, _ := ioutil.ReadAll(joinResponse.Body)
	defer joinResponse.Body.Close()
	if string(joinResponseBody) != string(joinRespBodyExpected) {
		t.Errorf("GetUser scenario wrong response from joinHandler: got %+v, expected %+v",
			string(joinResponseBody), string(joinRespBodyExpected))
	}

	/***************** Регистрация **********************/

	/***************** Получаем данные через куку *******/
	{
		getUserUrl := "http://localhost:8080/profile"
		getRespBodyMap := map[string]interface{}{
			"body": map[string]interface{}{
				"user": user.GetInfo(),
			},
			"status": 200,
		}
		getRespBodyExpected, err := json.Marshal(getRespBodyMap)
		if err != nil {
			t.Error(err)
		}
		var tmp []byte
		getUserRequest, err := http.NewRequest("GET", getUserUrl, bytes.NewBuffer(tmp))
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
	}

	/***************** Получаем данные через квери стринг **********************/
	{
		getUserUrl := "http://localhost:8080/profile?nickname=Love"

		getRespBodyMap := map[string]interface{}{
			"body": map[string]interface{}{
				"user": user.GetInfo(),
			},
			"status": 200,
		}
		getRespBodyExpected, err := json.Marshal(getRespBodyMap)
		if err != nil {
			t.Error(err)
		}
		var tmp []byte
		getUserRequest, err := http.NewRequest("GET", getUserUrl, bytes.NewBuffer(tmp))
		wGet := httptest.NewRecorder()
		apiHandler.GetUser(wGet, getUserRequest)
		getResponse := wGet.Result()
		getResponseBody, _ := ioutil.ReadAll(getResponse.Body)
		defer getResponse.Body.Close()

		if string(getResponseBody) != string(getRespBodyExpected) {
			t.Errorf("GetUser scenario wrong response from getHandler: got %+v, expected %+v",
				string(getResponseBody), string(getRespBodyExpected))
		}
	}

	/***************** неверный квери стринг **********************/
	{
		getUserUrl := "http://localhost:8080/profile?nickname=sss&nickname=aaa"

		getRespBodyMap := map[string]interface{}{
			"status": http.StatusBadRequest,
		}
		getRespBodyExpected, err := json.Marshal(getRespBodyMap)
		if err != nil {
			t.Error(err)
		}
		var tmp []byte
		getUserRequest, err := http.NewRequest("GET", getUserUrl, bytes.NewBuffer(tmp))
		wGet := httptest.NewRecorder()
		apiHandler.GetUser(wGet, getUserRequest)
		getResponse := wGet.Result()
		getResponseBody, _ := ioutil.ReadAll(getResponse.Body)
		defer getResponse.Body.Close()

		if string(getResponseBody) != string(getRespBodyExpected) {
			t.Errorf("GetUser scenario wrong response from getHandler: got %+v, expected %+v",
				string(getResponseBody), string(getRespBodyExpected))
		}
	}

	/***************** нет юзера **********************/
	{
		getUserUrl := "http://localhost:8080/profile?nickname=sss"

		getRespBodyMap := map[string]interface{}{
			"status": 404,
		}
		getRespBodyExpected, err := json.Marshal(getRespBodyMap)
		if err != nil {
			t.Error(err)
		}
		var tmp []byte
		getUserRequest, err := http.NewRequest("GET", getUserUrl, bytes.NewBuffer(tmp))
		wGet := httptest.NewRecorder()
		apiHandler.GetUser(wGet, getUserRequest)
		getResponse := wGet.Result()
		getResponseBody, _ := ioutil.ReadAll(getResponse.Body)
		defer getResponse.Body.Close()

		if string(getResponseBody) != string(getRespBodyExpected) {
			t.Errorf("GetUser scenario wrong response from getHandler: got %+v, expected %+v",
				string(getResponseBody), string(getRespBodyExpected))
		}
	}

	/***************** нет куки **********************/
	{
		getUserUrl := "http://localhost:8080/profile"
		getRespBodyMap := map[string]interface{}{
			"body": map[string]interface{}{
				"path": "/login",
			},
			"status": http.StatusSeeOther,
		}
		getRespBodyExpected, err := json.Marshal(getRespBodyMap)
		if err != nil {
			t.Error(err)
		}
		var tmp []byte
		getUserRequest, err := http.NewRequest("GET", getUserUrl, bytes.NewBuffer(tmp))
		wGet := httptest.NewRecorder()
		apiHandler.GetUser(wGet, getUserRequest)
		getResponse := wGet.Result()
		getResponseBody, _ := ioutil.ReadAll(getResponse.Body)
		defer getResponse.Body.Close()

		if string(getResponseBody) != string(getRespBodyExpected) {
			t.Errorf("GetUser scenario wrong response from getHandler: got %+v, expected %+v",
				string(getResponseBody), string(getRespBodyExpected))
		}
	}
}

func TestLogOut(t *testing.T) {

	t.Parallel()
	apiHandler := Handler{
		userStore:    CreateUserStore(),
		sessionStore: CreateSessionStore(),
	}

	user := User{
		Name:         "Антон",
		SurName:      "Гофер",
		NickName:     "Love",
		Email:        "aaa@mail.ru",
		Password:     "lovelove",
		PathToAvatar: defaultUserImgPath,
	}

	/***************** Регистрация **********************/

	joinUrl := "http://localhost:8080/join"

	joinRespBodyMap := map[string]interface{}{
		"body": map[string]interface{}{
			"path": "/",
		},
		"status": 200,
	}

	joinReqBody, err := json.Marshal(user)
	if err != nil {
		t.Error(err)
	}

	joinRespBodyExpected, err := json.Marshal(joinRespBodyMap)
	if err != nil {
		t.Error(err)
	}

	joinRequest, err := http.NewRequest("POST", joinUrl, bytes.NewBuffer(joinReqBody))
	w := httptest.NewRecorder()
	apiHandler.Join(w, joinRequest)
	joinResponse := w.Result()

	// check cookie
	ourCookies := joinResponse.Cookies()
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

	joinResponseBody, _ := ioutil.ReadAll(joinResponse.Body)
	defer joinResponse.Body.Close()

	if string(joinResponseBody) != string(joinRespBodyExpected) {
		t.Errorf("GetUser scenario wrong response from joinHandler: got %+v, expected %+v",
			string(joinResponseBody), string(joinRespBodyExpected))
	}

	/***************** Регистрация **********************/

	/***************** LOGOUT **********************/

	logoutUrl := "http://localhost:8080/logout"

	// без куки
	{

		logoutRespBodyMap := map[string]interface{}{
			"body": map[string]interface{}{
				"path": "/login",
			},
			"status": http.StatusSeeOther,
		}
		logoutRespBodyExpected, err := json.Marshal(logoutRespBodyMap)
		if err != nil {
			t.Error(err)
		}
		var tmp []byte
		logoutRequest, err := http.NewRequest("DELETE", logoutUrl, bytes.NewBuffer(tmp))
		wDelete := httptest.NewRecorder()
		apiHandler.LogOut(wDelete, logoutRequest)
		logoutResponse := wDelete.Result()
		logoutResponseBody, _ := ioutil.ReadAll(logoutResponse.Body)
		defer logoutResponse.Body.Close()
		if string(logoutResponseBody) != string(logoutRespBodyExpected) {
			t.Errorf("GetUser scenario wrong response from getHandler: got %+v, expected %+v",
				string(logoutResponseBody), string(logoutRespBodyExpected))
		}

	}

	// с кукой
	{

		logoutRespBodyMap := map[string]interface{}{
			"body": map[string]interface{}{
				"path": "/login",
			},
			"status": 200,
		}
		logoutRespBodyExpected, err := json.Marshal(logoutRespBodyMap)
		if err != nil {
			t.Error(err)
		}
		var tmp []byte
		logoutRequest, err := http.NewRequest("DELETE", logoutUrl, bytes.NewBuffer(tmp))
		wDelete := httptest.NewRecorder()
		logoutRequest.AddCookie(ourCookie)
		apiHandler.LogOut(wDelete, logoutRequest)
		logoutResponse := wDelete.Result()
		logoutResponseBody, _ := ioutil.ReadAll(logoutResponse.Body)
		defer logoutResponse.Body.Close()
		if string(logoutResponseBody) != string(logoutRespBodyExpected) {
			t.Errorf("GetUser scenario wrong response from getHandler: got %+v, expected %+v",
				string(logoutResponseBody), string(logoutRespBodyExpected))
		}

	}

	// с удаленной кукой
	{

		logoutRespBodyMap := map[string]interface{}{
			"body": map[string]interface{}{
				"path": "/login",
			},
			"status": http.StatusSeeOther,
		}
		logoutRespBodyExpected, err := json.Marshal(logoutRespBodyMap)
		if err != nil {
			t.Error(err)
		}
		var tmp []byte
		logoutRequest, err := http.NewRequest("DELETE", logoutUrl, bytes.NewBuffer(tmp))
		wDelete := httptest.NewRecorder()
		logoutRequest.AddCookie(ourCookie)
		apiHandler.LogOut(wDelete, logoutRequest)
		logoutResponse := wDelete.Result()
		logoutResponseBody, _ := ioutil.ReadAll(logoutResponse.Body)
		defer logoutResponse.Body.Close()
		if string(logoutResponseBody) != string(logoutRespBodyExpected) {
			t.Errorf("GetUser scenario wrong response from getHandler: got %+v, expected %+v",
				string(logoutResponseBody), string(logoutRespBodyExpected))
		}

	}

}
