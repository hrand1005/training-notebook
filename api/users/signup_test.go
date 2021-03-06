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

// TestSignup tests the API layer's Signup method for the Users resource.
// The test suite mocks the UserDB interface to test edge cases and error conditions.

func TestSignup(t *testing.T) {
	tests := []struct {
		name        string
		requestBody bytes.Buffer
		db          *data.MockUserDB
		wantCode    int
		wantResp    bytes.Buffer
	}{
		{
			name: "Valid request and DB call returns StatusCreated",
			requestBody: *bytes.NewBufferString(` {
				"name": "hildegard",
				"password": "12345"
			} `),
			db: &data.MockUserDB{
				AddUserStub: func(s *models.User) (models.UserID, error) {
					return 1, nil
				},
			},
			wantCode: http.StatusCreated,
			wantResp: *bytes.NewBufferString(` {
					"user-id": 1,
					"name": "hildegard"
			} `),
		},
		{
			name: "DB Error returns InternalServerError",
			requestBody: *bytes.NewBufferString(` {
				"name": "John",
				"password": "iamjohn44"
			} `),
			db: &data.MockUserDB{
				AddUserStub: func(s *models.User) (models.UserID, error) {
					return data.InvalidUserID, fmt.Errorf("Expected Error")
				},
			},
			wantCode: http.StatusInternalServerError,
			wantResp: *bytes.NewBufferString(` {
				"message": "Expected Error"
			} `),
		},
	}
	for _, v := range tests {
		// configure test case with data and test context
		u, err := New(v.db)
		if err != nil {
			t.Fail()
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// set the body in the test context's Request
		bodyReader := bytes.NewReader(v.requestBody.Bytes())
		// method/uri parsing exceed the scope of this test
		c.Request, _ = http.NewRequest("", "", bodyReader)

		// execute create with the test context
		u.Signup(c)

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
