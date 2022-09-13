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

func TestDeleteByID(t *testing.T) {
	const validUserID app.UserID = "validUserID"

	mockValidUserService := &mocks.UserService{
		DeleteByIDStub: func(id app.UserID) error {
			if id == validUserID {
				return nil
			}
			return fmt.Errorf("%w", app.ErrNotFound)
		},
	}

	mockInvalidUserService := &mocks.UserService{
		DeleteByIDStub: func(id app.UserID) error {
			return fmt.Errorf("unknown err")
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
			name:        "Nominal success returns 204",
			userID:      validUserID,
			userService: mockValidUserService,
			wantStatus:  fiber.StatusNoContent,
			wantBody:    nil,
		},
		{
			name:        "Non-existent user returns 404 and error response",
			userID:      "invalid-user-id",
			userService: mockValidUserService,
			wantStatus:  fiber.StatusNotFound,
			wantBody:    []byte(`{"errors":[{"message":"user not found"}]}`),
		},
		{
			name:        "Service error returns 500 and error response",
			userID:      validUserID,
			userService: mockInvalidUserService,
			wantStatus:  fiber.StatusInternalServerError,
			wantBody:    []byte(fmt.Sprintf(`{"errors":[{"message":"%v"}]}`, app.ErrServiceFailure.Error())),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewUserHandler(tc.userService, log.Default())

			testApp := fiber.New()
			testApp.Delete("/test/:userID", handler.DeleteByID)

			testURL := fmt.Sprintf("/test/%s", tc.userID)
			testReq, err := http.NewRequest(http.MethodDelete, testURL, nil)
			if err != nil {
				t.Fatalf("Failed to build test request:\n%#v", testReq)
			}

			testReq.Header.Add("Content-Type", "application/json")

			resp, _ := testApp.Test(testReq, 1)
			gotBody, _ := io.ReadAll(resp.Body)

			if tc.wantStatus != resp.StatusCode {
				t.Fatalf("Expected status code: %v\nGot status code: %v\nResp Body: %s", tc.wantStatus, resp.StatusCode, string(gotBody))
			}

			if tc.wantStatus != fiber.StatusNoContent {
				if !JSONBytesEqual(tc.wantBody, gotBody) {
					t.Fatalf("Expected response body: %v\nGot body: %v", string(tc.wantBody), string(gotBody))
				}
			}
		})
	}
}
