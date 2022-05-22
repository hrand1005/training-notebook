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

func TestUpdateSet(t *testing.T) {
	testCases := []struct {
		name string
		data []*data.Set
		// id from URL, not part of body for updates
		id          string
		requestBody bytes.Buffer
		wantCode    int
		wantResp    bytes.Buffer
		wantSet     *data.Set
	}{
		{
			name: "Valid set updated returns 200",
			data: []*data.Set{
				{
					ID:        1,
					Movement:  "Barbell Curl",
					Volume:    1,
					Intensity: 100,
				},
			},
			id: "1",
			requestBody: *bytes.NewBufferString(` {
					"movement": "Dumbbell Curl",
					"volume": 5,
					"intensity": 80
			} `),
			wantCode: 200,
			wantResp: *bytes.NewBufferString(` {
					"id": 1,
					"movement": "Dumbbell Curl",
					"volume": 5,
					"intensity": 80
			} `),
			wantSet: &data.Set{
				ID:        1,
				Movement:  "Dumbbell Curl",
				Volume:    5,
				Intensity: 80,
			},
		},
		{
			name: "Non-existent id returns 404",
			data: []*data.Set{
				{
					ID:        1,
					Movement:  "Barbell Curl",
					Volume:    1,
					Intensity: 100,
				},
			},
			id: "2",
			requestBody: *bytes.NewBufferString(` {
					"movement": "Dumbbell Curl",
					"volume": 5,
					"intensity": 80
			} `),
			wantCode: 404,
			wantResp: *bytes.NewBufferString(` {
				"message": "set not found"
			} `),
		},
		{
			name: "Invalid volume returns 400",
			id:   "1",
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
			id:   "1",
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
			id:   "1",
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

		// add id to URL params
		c.Params = append(c.Params, gin.Param{
			Key:   "id",
			Value: v.id,
		})
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

		// check that set has been added if valid
		if v.wantSet != nil {
			if len(testData.Sets()) != initialSetSize {
				t.Fatalf("Length of the data set has changed\nData: %v\n", testData.Sets())
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
