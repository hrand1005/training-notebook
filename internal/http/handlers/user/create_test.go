package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/hrand1005/training-notebook/internal/app"
	"github.com/hrand1005/training-notebook/internal/mocks"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		name        string
		userService app.UserService
		requestBody []byte
		wantStatus  int
		wantBody    []byte
	}{
		{
			name: "Nominal case returns 201 and users response",
			userService: &mocks.UserService{
				CreateStub: func(u *app.User) (app.UserID, error) {
					return "TestUserID", nil
				},
			},
			requestBody: []byte(`{"data":{"type":"user","attributes":{"first-name":"herb","last-name":"rand","email":"herb@yahoo.mail"}}}`),
			wantStatus:  fiber.StatusCreated,
			wantBody:    []byte(`{"data":{"id":"TestUserID","type":"","attributes":{"first-name":"herb","last-name":"rand","email":"herb@yahoo.mail"}}}`),
		},
		{
			name: "Service error returns 500 and error response",
			userService: &mocks.UserService{
				CreateStub: func(u *app.User) (app.UserID, error) {
					return "", fmt.Errorf("Expected Error")
				},
			},
			requestBody: []byte(`{"data":{"type":"user","attributes":{"first-name":"herb","last-name":"rand","email":"herb@yahoo.mail"}}}`),
			wantStatus:  fiber.StatusInternalServerError,
			wantBody:    []byte(`{"error":"Expected Error","message":"failed to create user"}`),
		},
		{
			name: "Invalid json returns 400 and error response",
			userService: &mocks.UserService{
				CreateStub: func(u *app.User) (app.UserID, error) {
					return "TestUserID", nil
				},
			},
			requestBody: []byte(`{"data":`),
			wantStatus:  fiber.StatusBadRequest,
			wantBody:    []byte(`{"error":"unexpected end of JSON input","message":"invalid json"}`),
		},
		{
			name: "Invalid field values returns 400 and error response",
			userService: &mocks.UserService{
				CreateStub: func(u *app.User) (app.UserID, error) {
					return "TestUserID", nil
				},
			},
			requestBody: []byte(`{"data":{"type":"user","attributes":{"first-name":"h","last-name":"rand","email":"herb@yahoo.mail"}}}`),
			wantStatus:  fiber.StatusBadRequest,
			wantBody:    []byte(`[{"invalid-field":"RequestBody.Data.Attributes.FirstName","tag":"min","value":"2"}]`),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewUserHandler(tc.userService, log.Default())

			testApp := fiber.New()
			testApp.Post("/test", handler.Create)

			testReq, err := http.NewRequest(http.MethodPost, "/test", bytes.NewReader(tc.requestBody))
			if err != nil {
				t.Fatalf("Failed to build test request:\n%#v", testReq)
			}

			testReq.Header.Add("Content-Type", "application/json")

			resp, _ := testApp.Test(testReq, 1)
			gotBody, _ := io.ReadAll(resp.Body)

			if tc.wantStatus != resp.StatusCode {
				t.Fatalf("Expected status code: %v\nGot status code: %v\nResp Body: %s", tc.wantStatus, resp.StatusCode, string(gotBody))
			}

			if !JSONBytesEqual(tc.wantBody, gotBody) {
				t.Fatalf("Expected response body: %v\nGot body: %v", string(tc.wantBody), string(gotBody))
			}
		})
	}
}

// JSONBytesEqual is a utility function for comparing responses.
func JSONBytesEqual(a, b []byte) bool {
	var aDecode, bDecode interface{}
	if err := json.Unmarshal(a, &aDecode); err != nil {
		// log.Printf("Encountered error decoding:\n%v\nErr:%v", string(a), err)
		return false
	}
	if err := json.Unmarshal(b, &bDecode); err != nil {
		// log.Printf("Encountered error decoding:\n%v\nErr:%v", string(b), err)
		return false
	}

	return reflect.DeepEqual(aDecode, bDecode)
}
