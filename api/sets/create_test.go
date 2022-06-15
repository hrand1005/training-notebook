package sets

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

// TestCreateSet tests the API layer's Create method for the Sets resource.
// The test suite mocks the SetDB interface to test edge cases and error conditions.
func TestCreateSet(t *testing.T) {
	tests := []struct {
		name        string
		requestBody bytes.Buffer
		db          *data.MockSetDB
		wantCode    int
		wantResp    bytes.Buffer
	}{
		{
			name: "Valid request and db call returns StatusCreated",
			requestBody: *bytes.NewBufferString(` {
					"movement": "Barbell Curl",
					"volume": 1,
					"intensity": 100
			} `),
			db: &data.MockSetDB{
				AddSetStub: func(s *models.Set) (models.SetID, error) {
					return 1, nil
				},
			},
			wantCode: http.StatusCreated,
			wantResp: *bytes.NewBufferString(` {
					"set-id": 1,
					"movement": "Barbell Curl",
					"volume": 1,
					"intensity": 100
			} `),
		},
		{
			name: "DB Error returns InternalServerError",
			requestBody: *bytes.NewBufferString(` {
					"movement": "Barbell Curl",
					"volume": 1,
					"intensity": 100
			} `),
			db: &data.MockSetDB{
				AddSetStub: func(s *models.Set) (models.SetID, error) {
					return data.InvalidSetID, fmt.Errorf("Expected Error")
				},
			},
			wantCode: http.StatusInternalServerError,
			wantResp: *bytes.NewBufferString(` {
				"message": "Expected Error"
			} `),
		},
		{
			name: "Missing movement returns StatusBadRequest",
			requestBody: *bytes.NewBufferString(`{
					"volume": 1,
					"intensity": 100
			}`),
			wantCode: http.StatusBadRequest,
			wantResp: *bytes.NewBufferString(` {
				"message": "'Movement' field must use unicode characters."
			} `),
		},
		{
			name: "Missing volume returns StatusBadRequest",
			requestBody: *bytes.NewBufferString(`{
					"movement": "Squat",
					"intensity": 100
			}`),
			wantCode: http.StatusBadRequest,
			wantResp: *bytes.NewBufferString(` {
				"message": "'Volume' field must be greater than 0."
			} `),
		},
		{
			name: "Zero Volume returns StatusBadRequest",
			requestBody: *bytes.NewBufferString(`{
				 	"movement": "Press",
					"volume": 0,
					"intensity": 80
			}`),
			wantCode: http.StatusBadRequest,
			wantResp: *bytes.NewBufferString(` {
				"message": "'Volume' field must be greater than 0."
			} `),
		},
		{
			name: "Missing intensity returns StatusBadRequest",
			requestBody: *bytes.NewBufferString(`{
				 	"movement": "Press",
					"volume": 5
			}`),
			wantCode: http.StatusBadRequest,
			wantResp: *bytes.NewBufferString(` {
				"message": "'Intensity' field must be greater than 0."
			} `),
		},
		{
			name: "Zero intensity returns StatusBadRequest",
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
			name: "101 intensity returns StatusBadRequest",
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
		// configure test case with data and test context
		ts, err := New(v.db)
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
		ts.Create(c)

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
