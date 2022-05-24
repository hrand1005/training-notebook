package sets

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
)

func TestCreateSet(t *testing.T) {
	testCases := []struct {
		name        string
		data        []*data.Set
		requestBody bytes.Buffer
		wantCode    int
		wantResp    bytes.Buffer
		wantSet     *data.Set
	}{
		{
			name: "Valid set created returns 200",
			requestBody: *bytes.NewBufferString(` {
					"movement": "Barbell Curl",
					"volume": 1,
					"intensity": 100
			} `),
			wantCode: 201,
			wantResp: *bytes.NewBufferString(` {
					"id": 1,
					"movement": "Barbell Curl",
					"volume": 1,
					"intensity": 100
			} `),
			wantSet: &data.Set{
				ID:        1,
				Movement:  "Barbell Curl",
				Volume:    1,
				Intensity: 100,
			},
		},
		{
			name: "Invalid volume returns 400",
			requestBody: *bytes.NewBufferString(` {
					"movement": "Barbell Curl",
					"volume": 0,
					"intensity": 0.5
			} `),
			wantCode: 400,
			wantResp: *bytes.NewBufferString(` {
				"message": "Key: 'Set.Volume' Error:Field validation for 'Volume' failed on the 'gt' tag"
			} `),
		},
		{
			name: "0 intensity returns 400",
			requestBody: *bytes.NewBufferString(` {
					"movement": "Barbell Curl",
					"volume": 2,
					"intensity": 0
			} `),
			wantCode: 400,
			wantResp: *bytes.NewBufferString(` {
				"message": "Key: 'Set.Intensity' Error:Field validation for 'Intensity' failed on the 'gt' tag"
			} `),
		},
		{
			name: "101 intensity returns 400",
			requestBody: *bytes.NewBufferString(` {
					"movement": "Barbell Curl",
					"volume": 2,
					"intensity": 101
			} `),
			wantCode: 400,
			wantResp: *bytes.NewBufferString(` {
				"message": "Key: 'Set.Intensity' Error:Field validation for 'Intensity' failed on the 'lte' tag"
			} `),
		},
	}

	for _, v := range testCases {
		// configure test case with data and test context
		initialSetSize := len(v.data)
		testData := data.NewSetData(v.data)
		ts, err := New(testData)
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

		// check that set has been added if valid
		if v.wantSet != nil {
			sets, _ := testData.Sets()
			if len(sets) != initialSetSize+1 {
				t.Fatalf("Set not added to the data set\nData: %v\n", sets)
			}

			// compare retrieved set with expected
			gotSet, err := testData.SetByID(v.wantSet.ID)
			if err != nil {
				t.Fatalf("Set with wantSet.ID not found\nWanted set: %v", v.wantSet)
			}

			if !reflect.DeepEqual(v.wantSet, gotSet) {
				t.Fatalf("Wanted set: %+v\nGot set: %+v\n", v.wantSet, gotSet)
			}
		}
	}
}
