package sets

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/api/users"
	"github.com/hrand1005/training-notebook/data"
	"github.com/hrand1005/training-notebook/models"
)

// TestUpdateSet tests the API layer's Update method for the Sets resource.
// The test suite mocks the SetDB interface to test edge cases and error conditions.
func TestUpdateSet(t *testing.T) {
	tests := []struct {
		name string
		db   *data.MockSetDB
		// id from URL, not part of body for updates
		setID       string
		userID      models.UserID
		requestBody bytes.Buffer
		wantCode    int
		wantResp    bytes.Buffer
	}{
		{
			name: "Valid set updated returns StatusOK",
			db: &data.MockSetDB{
				UpdateSetForUserStub: func(setID models.SetID, userID models.UserID, s *models.Set) error {
					return nil
				},
			},
			setID:  "1",
			userID: 1,
			requestBody: *bytes.NewBufferString(` {
					"movement": "Dumbbell Curl",
					"volume": 5,
					"intensity": 80
			} `),
			wantCode: http.StatusOK,
			wantResp: *bytes.NewBufferString(` {
					"set-id": 1,
					"user-id": 1,
					"movement": "Dumbbell Curl",
					"volume": 5,
					"intensity": 80
			} `),
		},
		{
			name: "Invalid id param returns StatusBadRequest",
			db: &data.MockSetDB{
				UpdateSetForUserStub: func(setID models.SetID, userID models.UserID, s *models.Set) error {
					return data.ErrNotFound
				},
			},
			setID:  "-1",
			userID: 1,
			requestBody: *bytes.NewBufferString(` {
					"movement": "Dumbbell Curl",
					"volume": 5,
					"intensity": 80
			} `),
			wantCode: http.StatusBadRequest,
			wantResp: *bytes.NewBufferString(` {
				"message": "Invalid set ID"
			} `),
		},
		{
			name: "Set not found returns StatusNotFound",
			db: &data.MockSetDB{
				UpdateSetForUserStub: func(setID models.SetID, userID models.UserID, s *models.Set) error {
					return data.ErrNotFound
				},
			},
			setID:  "2",
			userID: 1,
			requestBody: *bytes.NewBufferString(` {
					"movement": "Dumbbell Curl",
					"volume": 5,
					"intensity": 80
			} `),
			wantCode: http.StatusNotFound,
			wantResp: *bytes.NewBufferString(` {
				"message": "no such set with id 2"
			} `),
		},
		{
			name: "Invalid db call returns InternalServerError",
			db: &data.MockSetDB{
				UpdateSetForUserStub: func(setID models.SetID, userID models.UserID, s *models.Set) error {
					return fmt.Errorf("Expected error")
				},
			},
			setID:  "2",
			userID: 1,
			requestBody: *bytes.NewBufferString(` {
					"movement": "Dumbbell Curl",
					"volume": 5,
					"intensity": 80
			} `),
			wantCode: http.StatusInternalServerError,
			wantResp: *bytes.NewBufferString(` {
				"message": "Expected error"
			} `),
		},
		{
			name:   "Missing volume returns StatusBadRequest",
			setID:  "1",
			userID: 1,
			requestBody: *bytes.NewBufferString(` {
					"movement": "Barbell Curl",
					"intensity": 0.5
			} `),
			wantCode: http.StatusBadRequest,
			wantResp: *bytes.NewBufferString(` {
				"message": "'Volume' field must be greater than 0."
			} `),
		},
		{
			name:   "Zero Volume returns StatusBadRequest",
			setID:  "1",
			userID: 1,
			requestBody: *bytes.NewBufferString(` {
					"movement": "Barbell Curl",
					"volume": 0,
					"intensity": 40
			} `),
			wantCode: http.StatusBadRequest,
			wantResp: *bytes.NewBufferString(` {
				"message": "'Volume' field must be greater than 0."
			} `),
		},
		{
			name:   "Missing intensity returns StatusBadRequest",
			setID:  "1",
			userID: 1,
			requestBody: *bytes.NewBufferString(` {
					"movement": "Barbell Curl",
					"volume": 2
			} `),
			wantCode: http.StatusBadRequest,
			wantResp: *bytes.NewBufferString(` {
				"message": "'Intensity' field must be greater than 0."
			} `),
		},
		{
			name:   "Zero intensity returns StatusBadRequest",
			setID:  "1",
			userID: 1,
			requestBody: *bytes.NewBufferString(`{
				 	"movement": "Press",
					"volume": 5,
					"intensity": 0
			}`),
			wantCode: http.StatusBadRequest,
			wantResp: *bytes.NewBufferString(` {
				"message": "'Intensity' field must be greater than 0."
			} `),
		},
		{
			name:   "101 intensity returns StatusBadRequest",
			setID:  "1",
			userID: 1,
			requestBody: *bytes.NewBufferString(`{
				 	"movement": "Press",
					"volume": 5,
					"intensity": 101
			}`),
			wantCode: http.StatusBadRequest,
			wantResp: *bytes.NewBufferString(` {
				"message": "'Intensity' field must be no more than 100."
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

		// set the body in the test context's Request
		bodyReader := bytes.NewReader(v.requestBody.Bytes())

		// add id to URL params
		c.AddParam(SetIDFromParamsKey, v.setID)
		c.Request, _ = http.NewRequest("", "", bodyReader)
		c.Set(users.UserIDFromContextKey, v.userID)

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
