package main

import (
	//"io"
	//"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
	//"log"
	"bytes"
)


type TestCase struct {
	requestBody map[string]string
	responseBody map[string]string
	statusCode int
}



func TestJoinHandler(t *testing.T) {
	joinUrl := "http://localhost:8080/join"

	cases := []TestCase{
		TestCase{
			requestBody: map[string]string{	"name" : "Антон",
								"surname" : "Пися",
								"nickname" : "АнтонLove",
								"password" : "lovelove",
			},

			responseBody: map[string]string{
				"answer" : `{"join"}`,
			},

			statusCode: 300,
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

		request, err := http.NewRequest("POST", joinUrl, bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		api := &Handler{
			userStore:    CreateUserStore(),
			sessionStore: CreateSessionStore(),
		}

		api.Join(w, request)

		if w.Code != test.statusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				num, w.Code, test.statusCode)
		}		

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