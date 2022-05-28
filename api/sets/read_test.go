package sets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
)

func TestReadSingleSet(t *testing.T) {
	testCases := []struct {
		name     string
		data     []*data.Set
		params   gin.Params
		wantCode int
		wantResp bytes.Buffer
	}{
		{
			name:   "Set found returns 200",
			params: []gin.Param{{Key: "id", Value: "1"}},
			data: []*data.Set{
				{
					ID:        1,
					Movement:  "Squat",
					Volume:    5,
					Intensity: 80,
				},
			},
			wantCode: 200,
			wantResp: *bytes.NewBufferString(`
				{
					"id": 1,
					"movement": "Squat",
					"volume": 5,
					"intensity": 80
				}
			`),
		},
		{
			name:   "Set not found returns 404",
			params: []gin.Param{{Key: "id", Value: "4"}},
			data: []*data.Set{
				{
					ID:        1,
					Movement:  "Squat",
					Volume:    5,
					Intensity: 80,
				},
			},
			wantCode: 404,
			wantResp: *bytes.NewBufferString(`{
				"message": "resource not found"
			}`),
		},
		{
			name:   "Invalid params returns 400",
			params: []gin.Param{{Key: "bad", Value: "request"}},
			data: []*data.Set{
				{
					ID:        1,
					Movement:  "Squat",
					Volume:    5,
					Intensity: 80,
				},
			},
			wantCode: 400,
			wantResp: *bytes.NewBufferString(fmt.Sprintf(`{
				"message": %q
			}`, ErrInvalidSetID)),
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
		c.Params = v.params

		// execute Read on test context
		ts.Read(c)

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

func TestReadAllSets(t *testing.T) {
	testCases := []struct {
		name     string
		data     []*data.Set
		wantCode int
		wantResp bytes.Buffer
	}{
		{
			name: "1 set returns 200",
			data: []*data.Set{
				{
					ID:        1,
					Movement:  "Squat",
					Volume:    5,
					Intensity: 80,
				},
			},
			wantCode: 200,
			wantResp: *bytes.NewBufferString(`[
				{
					"id": 1,
					"movement": "Squat",
					"volume": 5,
					"intensity": 80
				}
			]`),
		},
		{
			name:     "No sets returns 200",
			wantCode: 200,
			wantResp: *bytes.NewBufferString(`[]`),
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

		// execute ReadAll with test context
		ts.ReadAll(c)

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

// JSONBytesEqual compares the JSON in two byte slices.
func JSONBytesEqual(a, b []byte) (bool, error) {
	var j, j2 interface{}
	if err := json.Unmarshal(a, &j); err != nil {
		return false, err
	}
	if err := json.Unmarshal(b, &j2); err != nil {
		return false, err
	}
	return reflect.DeepEqual(j2, j), nil
}
