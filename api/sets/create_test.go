package sets

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
)

// TestCreateSet_WithMock allows us to test the API alone while mocking the DB, now that
// DB unit tests exist. NOTE: the idea is to eventually replace the above test with proper
// integration tests.
func TestCreateSet(t *testing.T) {
	tests := []struct {
		name        string
		requestBody bytes.Buffer
		db          *data.MockSetDB
		wantCode    int
		wantResp    bytes.Buffer
	}{
		{
			name: "Valid request and db call returns 200",
			requestBody: *bytes.NewBufferString(` {
					"movement": "Barbell Curl",
					"volume": 1,
					"intensity": 100
			} `),
			db: &data.MockSetDB{
				AddSetStub: func(s *data.Set) (data.SetID, error) {
					return 1, nil
				},
			},
			wantCode: http.StatusCreated,
			wantResp: *bytes.NewBufferString(` {
					"id": 1,
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
				AddSetStub: func(s *data.Set) (data.SetID, error) {
					return data.InvalidID, fmt.Errorf("Expected Error")
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
				"message": "Key: 'Set.Movement' Error:Field validation for 'Movement' failed on the 'movement' tag"
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
				"message": "Key: 'Set.Volume' Error:Field validation for 'Volume' failed on the 'gt' tag"
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
				"message": "Key: 'Set.Volume' Error:Field validation for 'Volume' failed on the 'gt' tag"
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
				"message": "Key: 'Set.Intensity' Error:Field validation for 'Intensity' failed on the 'gt' tag"
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
				"message": "Key: 'Set.Intensity' Error:Field validation for 'Intensity' failed on the 'gt' tag"
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
				"message": "Key: 'Set.Intensity' Error:Field validation for 'Intensity' failed on the 'lte' tag"
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
