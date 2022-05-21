package sets

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
)

func TestCreateSet(t *testing.T) {

	// test cases:
	// Bad Request
	// success case

	testCases := []struct {
		name     string
		data     []*data.Set
		wantCode int
		set      data.Set
		wantResp bytes.Buffer
	}{
		{
			name:     "Valid set created returns 200",
			wantCode: 201,
			set: data.Set{
				Movement:  "Barbell Curl",
				Volume:    12,
				Intensity: 0.5,
			},
			wantResp: *bytes.NewBufferString(` {
					"id": 1,
					"movement": "Barbell Curl",
					"volume": 12,
					"intensity": 0.5
				} `),
		},
	}

	for _, v := range testCases {
		// configure test case with data and test context
		testData := data.NewSetData(v.data)
		ts, err := New(testData)
		if err != nil {
			t.Fail()
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// set 'newSet' in the context
		c.Set("newSet", v.set)

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

// TODO: Decide how to Test Set Validator Middleware
// {
// 	name:     "Invalid set returns 400",
// 	wantCode: 400,
// 	set:      data.Set{},
// 	wantResp: *bytes.NewBufferString(` {
// 		"message": "invalid request parameters"
// 		} `),
// },
