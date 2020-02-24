package main

import (
	//"io"
	//"fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	//"log"
	"bytes"
)

var apiHandler *Handler

type TestCase struct {
	requestBody  map[string]interface{}
	responseBody map[string]interface{}
	apiMethod func(http.ResponseWriter, *http.Request)
}

func TestJoinHandler(t *testing.T) {
	// в начале каждого тестового сценария принудительно сбрасываем хранилище
	apiHandler = &Handler {
		userStore: CreateUserStore(),
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
				"status": 308,
			},

			apiMethod: apiHandler.Join,
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
				"status": 409,
			},

			apiMethod: apiHandler.Join,
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

			apiMethod: apiHandler.Join,
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

		request, err := http.NewRequest("POST", joinUrl, bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		test.apiMethod(w, request)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		bodyStr := string(body)
		if bodyStr != string(resBody) {
			t.Errorf("[%d] wrong Response: got %+v, expected %+v",
				num, bodyStr, string(resBody))
		}

	}
}

func TestLoginHandler(t *testing.T) {
	apiHandler = &Handler {
		userStore: CreateUserStore(),
		sessionStore: CreateSessionStore(),
	}

	loginUrl := "http://localhost:8080/login"

	cases := []TestCase{
		TestCase {
			requestBody: map[string]interface{} {	
				"nickname" : "test",
				"password" : "testtest",
			},

			responseBody: map[string]interface{} {
				"status" : 404,
			},

			apiMethod : apiHandler.LogIn,
		},

		TestCase {
			requestBody: map[string]interface{} {	
				"name" : "Антон",
				"surname" : "Гофер",
				"nickname" : "Love",
				"password" : "lovelove",
			},

			responseBody: map[string]interface{} {
				"status" : 308,
			},

			apiMethod: apiHandler.Join,
		},

		TestCase {
			requestBody: map[string]interface{} {	
				"nickname" : "Love",
				"password" : "lovelove",
			},

			responseBody: map[string]interface{} {
				"status" : 308,
			},

			apiMethod : apiHandler.LogIn,
		},

	}

	for num, test := range cases {
		reqBody, err :=  json.Marshal(test.requestBody)
		if err != nil {
			t.Error(err)
		}
		resBody, err := json.Marshal(test.responseBody)
		if err != nil {
			t.Error(err)
		}

		request, err := http.NewRequest("POST", loginUrl, bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		test.apiMethod(w, request)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		bodyStr := string(body)
		if bodyStr != string(resBody) {
			t.Errorf("[%d] wrong Response: got %+v, expected %+v",
				num, bodyStr, string(resBody))
		}

	}
}
