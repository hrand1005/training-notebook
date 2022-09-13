package users

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/hrand1005/training-notebook/internal/app"
	"github.com/hrand1005/training-notebook/internal/mocks"
)

func TestReadByID(t *testing.T) {
	const validUserID app.UserID = "valid-id"
	validTestUser := &app.User{
		ID:           validUserID,
		FirstName:    "valid-first",
		LastName:     "valid-last",
		Email:        "valid-user@yahoo.mail",
		PasswordHash: "12345password",
	}

	mockValidUserService := &mocks.UserService{
		ReadByIDStub: func(id app.UserID) (*app.User, error) {
			if id == validUserID {
				return validTestUser, nil
			}
			return nil, fmt.Errorf("%w", app.ErrNotFound)
		},
	}

	mockInvalidUserService := &mocks.UserService{
		ReadByIDStub: func(id app.UserID) (*app.User, error) {
			return nil, fmt.Errorf("unknown err")
		},
	}

	tests := []struct {
		name        string
		userID      app.UserID
		userService app.UserService
		wantStatus  int
		wantBody    []byte
	}{
		{
			name:        "Nominal success returns 200 and user response",
			userID:      validUserID,
			userService: mockValidUserService,
			wantStatus:  fiber.StatusOK,
			wantBody:    []byte(`{"data":{"id":"valid-id","type":"user","attributes":{"first-name":"valid-first","last-name":"valid-last","email":"valid-user@yahoo.mail"}}}`),
		},
		{
			name:        "Non-existent user returns 404 and error response",
			userID:      "invalid-id",
			userService: mockValidUserService, wantStatus: fiber.StatusNotFound,
			wantBody: []byte(`{"errors":[{"message":"user not found"}]}`),
		},
		{
			name:        "Service error returns 500 and error response",
			userID:      "valid-id",
			userService: mockInvalidUserService,
			wantStatus:  fiber.StatusInternalServerError,
			wantBody:    []byte(fmt.Sprintf(`{"errors":[{"message":"%v"}]}`, app.ErrServiceFailure.Error())),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewUserHandler(tc.userService, log.Default())

			testApp := fiber.New()
			testApp.Get("/test/:userID", handler.ReadByID)

			testURL := fmt.Sprintf("/test/%s", tc.userID)
			testReq, err := http.NewRequest(http.MethodGet, testURL, nil)
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
