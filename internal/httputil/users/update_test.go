package users

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/hrand1005/training-notebook/internal/app"
	"github.com/hrand1005/training-notebook/internal/mocks"
)

func TestUpdateByID(t *testing.T) {
	var validUserID app.UserID = "valid-id"

	tests := []struct {
		name        string
		userService app.UserService
		userID      app.UserID
		requestBody []byte
		wantStatus  int
		wantBody    []byte
	}{
		{
			name: "Nominal case returns 200 and user response",
			userService: &mocks.UserService{
				UpdateByIDStub: func(id app.UserID, u *app.User) error {
					return nil
				},
			},
			userID:      validUserID,
			requestBody: []byte(`{"data":{"type":"user","attributes":{"first-name":"herb","last-name":"rand","email":"herb@yahoo.mail"}}}`),
			wantStatus:  fiber.StatusOK,
			wantBody:    []byte(`{"data":{"id":"valid-id","type":"user","attributes":{"first-name":"herb","last-name":"rand","email":"herb@yahoo.mail"}}}`),
		},
		{
			name: "Service error returns 500 and error response",
			userService: &mocks.UserService{
				UpdateByIDStub: func(id app.UserID, u *app.User) error {
					return fmt.Errorf("Expected Error")
				},
			},
			userID:      validUserID,
			requestBody: []byte(`{"data":{"type":"user","attributes":{"first-name":"herb","last-name":"rand","email":"herb@yahoo.mail"}}}`),
			wantStatus:  fiber.StatusInternalServerError,
			wantBody:    []byte(`{"errors":[{"message":"the service failed due to internal reasons"}]}`),
		},
		{
			name: "Invalid json returns 400 and error response",
			userService: &mocks.UserService{
				UpdateByIDStub: func(id app.UserID, u *app.User) error {
					return nil
				},
			},
			userID:      validUserID,
			requestBody: []byte(`{"data":`),
			wantStatus:  fiber.StatusBadRequest,
			wantBody:    []byte(`{"errors":[{"message":"invalid json: unexpected end of JSON input"}]}`),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewUserHandler(tc.userService, log.Default())

			testApp := fiber.New()
			testApp.Put("/test/:userID", handler.UpdateByID)

			testURL := fmt.Sprintf("/test/%s", tc.userID)
			testReq, err := http.NewRequest(http.MethodPut, testURL, bytes.NewReader(tc.requestBody))
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
