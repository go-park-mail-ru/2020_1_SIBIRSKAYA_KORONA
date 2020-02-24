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

type TestCase struct {
	requestBody  map[string]interface{}
	responseBody map[string]interface{}
}

func TestJoinHandler(t *testing.T) {
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
		},
	}

	// в рамках тестового сценария одного хэндлера используем постоянное хранилище
	api := &Handler{
		userStore:    CreateUserStore(),
		sessionStore: CreateSessionStore(),
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

		api.Join(w, request)

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

// func TestLoginHandler(t *testing.T) {

// }
