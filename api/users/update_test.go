package users

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
	"github.com/hrand1005/training-notebook/models"
)

// TestUpdateUser tests the API layer's Update method for the Users resource.
// The test suite mocks the UserDB interface to test edge cases and error conditions.
func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name string
		db   *data.MockUserDB
		// id from URL, not part of body for updates
		id          string
		requestBody bytes.Buffer
		wantCode    int
		wantResp    bytes.Buffer
	}{
		{
			name: "Valid user updated returns StatusOK",
			db: &data.MockUserDB{
				UpdateUserStub: func(id models.UserID, s *models.User) error {
					return nil
				},
			},
			id: "1",
			requestBody: *bytes.NewBufferString(` {
				"name": "Keonwoo Oh"
			} `),
			wantCode: http.StatusOK,
			wantResp: *bytes.NewBufferString(` {
					"id": 1,
					"name": "Keonwoo Oh"
			} `),
		},
		{
			name: "Invalid id param returns StatusBadRequest",
			db: &data.MockUserDB{
				UpdateUserStub: func(id models.UserID, s *models.User) error {
					return data.ErrNotFound
				},
			},
			id: "-1",
			requestBody: *bytes.NewBufferString(` {
				"name": "someone"
			} `),
			wantCode: http.StatusBadRequest,
			wantResp: *bytes.NewBufferString(` {
				"message": "Invalid user ID"
			} `),
		},
		{
			name: "User not found returns StatusNotFound",
			db: &data.MockUserDB{
				UpdateUserStub: func(id models.UserID, s *models.User) error {
					return data.ErrNotFound
				},
			},
			id: "2",
			requestBody: *bytes.NewBufferString(` {
				"name": "unknown"
			} `),
			wantCode: http.StatusNotFound,
			wantResp: *bytes.NewBufferString(` {
				"message": "no such user with id 2"
			} `),
		},
		{
			name: "Invalid db call returns InternalServerError",
			db: &data.MockUserDB{
				UpdateUserStub: func(id models.UserID, s *models.User) error {
					return fmt.Errorf("Expected error")
				},
			},
			id: "2",
			requestBody: *bytes.NewBufferString(` {
				"name": "some user"
			} `),
			wantCode: http.StatusInternalServerError,
			wantResp: *bytes.NewBufferString(` {
				"message": "Expected error"
			} `),
		},
	}

	for _, v := range tests {
		ts, err := New(v.db)
		if err != nil {
			t.Fail()
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// user the body in the test context's Request
		bodyReader := bytes.NewReader(v.requestBody.Bytes())

		// add id to URL params
		c.AddParam("id", v.id)
		c.Request, _ = http.NewRequest("", "", bodyReader)

		// execute update with the test context
		ts.Update(c)

		// check response code
		if v.wantCode != w.Code {
			t.Fatalf("Wanted code: %v\nGot code: %v\n", v.wantCode, w.Code)
		}

		// check response body
		if equal, _ := JSONBytesEqual(v.wantResp.Bytes(), w.Body.Bytes()); !equal {
			t.Fatalf("Wanted body: %v\nGot body: %v\n", v.wantResp.String(), w.Body.String())
		}
	}
}
