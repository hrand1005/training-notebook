package sets

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
)

func TestDeleteSet(t *testing.T) {
	testCases := []struct {
		name string
		data []*data.Set
		// id from URL, not part of body for updates
		id       string
		wantCode int
		wantResp bytes.Buffer
	}{
		{
			name: "Valid delete returns 204",
			data: []*data.Set{
				{
					ID:        1,
					Movement:  "Barbell Curl",
					Volume:    1,
					Intensity: 100,
				},
			},
			id:       "1",
			wantCode: 204,
		},
		{
			name: "Nonexistent set returns 404",
			data: []*data.Set{
				{
					ID:        1,
					Movement:  "Barbell Curl",
					Volume:    1,
					Intensity: 100,
				},
			},
			id:       "2",
			wantCode: 404,
			wantResp: *bytes.NewBufferString(`{
				"message": "set not found"
			}`),
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

		// add id to URL params
		c.Params = append(c.Params, gin.Param{
			Key:   "id",
			Value: v.id,
		})

		// execute delete with the test context
		ts.Delete(c)

		// check response code
		if v.wantCode != w.Code {
			t.Fatalf("Wanted code: %v\nGot code: %v\n", v.wantCode, w.Code)
		}

		// DELETE success case should be StatusNoContent only
		if w.Code == http.StatusNoContent {
			// check that set has been added if valid
			sets, _ := testData.Sets()
			if len(sets) != initialSetSize-1 {
				t.Fatalf("Length of the data set did not decrease after deleting\nData: %v\n", sets)
			}

			if w.Body.String() != v.wantResp.String() {
				t.Fatalf("Body should not be set for StatusNoContent responses")
			}
		} else {
			// no data changes should occur in the failure case
			sets, _ := testData.Sets()
			if len(sets) != initialSetSize {
				t.Fatalf("Data changes should not occur when delete fails")
			}
			// check response body
			if equal, _ := JSONBytesEqual(v.wantResp.Bytes(), w.Body.Bytes()); !equal {
				t.Fatalf("Wanted body: %v\nGot body: %v\n", v.wantResp.String(), w.Body.String())
			}
		}
	}
}
