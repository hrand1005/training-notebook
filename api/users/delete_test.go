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

// TestDeleteUser tests the API layer's Delete method for the Users resource.
// The test suite mocks the UserDB interface to test edge cases and error conditions.
func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name     string
		db       *data.MockUserDB
		id       string
		wantCode int
		wantResp bytes.Buffer
	}{
		{
			name: "User found with valid db call returns StatusNoContent",
			db: &data.MockUserDB{
				DeleteUserStub: func(id models.UserID) error {
					return nil
				},
			},
			id:       "1",
			wantCode: http.StatusNoContent,
		},
		{
			name: "User not found returns StatusNotFound",
			db: &data.MockUserDB{
				DeleteUserStub: func(id models.UserID) error {
					return data.ErrNotFound
				},
			},
			id:       "4",
			wantCode: http.StatusNotFound,
			wantResp: *bytes.NewBufferString(`{
				"message": "no such user with id 4"
			}`),
		},
		{
			name: "Invalid db query returns InternalServerError",
			db: &data.MockUserDB{
				DeleteUserStub: func(id models.UserID) error {
					return fmt.Errorf("Expected error")
				},
			},
			id:       "4",
			wantCode: http.StatusInternalServerError,
			wantResp: *bytes.NewBufferString(`{
				"message": "Expected error"
			}`),
		},
		{
			name:     "Invalid params returns StatusBadRequest",
			id:       "-1",
			wantCode: http.StatusBadRequest,
			wantResp: *bytes.NewBufferString(fmt.Sprintf(`{
				"message": %q
			}`, ErrInvalidUserID)),
		},
	}

	for _, v := range tests {
		// configure test case with data and test context
		ts, err := New(v.db)
		if err != nil {
			t.Fail()
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.AddParam(UserIDFromParamsKey, v.id)

		// execute Delete on test context
		ts.Delete(c)

		// check response code
		if v.wantCode != w.Code {
			t.Fatalf("Wanted code: %v\nGot code: %v\n", v.wantCode, w.Code)
		}

		if v.wantCode != http.StatusNoContent {
			// check response body
			if equal, _ := JSONBytesEqual(v.wantResp.Bytes(), w.Body.Bytes()); !equal {
				t.Fatalf("Wanted body: %v\nGot body: %v\n", v.wantResp.String(), w.Body.String())
			}
		}
	}
}
