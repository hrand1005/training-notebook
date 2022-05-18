package sets

import (
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
)

func TestReadSingleSet(t *testing.T) {

	// test cases:
	// id not found
	// invalid params or smth
	// success case

	testCases := []struct {
		name     string
		data     []*data.Set
		params   gin.Params
		wantCode int
		//		wantResp *bytes.Buffer
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
			//wantResp:
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
			//wantResp:
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
			//wantResp:
		},
	}

	for _, v := range testCases {
		testData := data.NewSetData(v.data)

		// create test set object
		ts, err := New(testData)
		if err != nil {
			t.Fail()
		}

		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = v.params

		ts.Read(c)

		if w.Code != v.wantCode {
			// TODO: print body
			t.Fatalf("Wanted code: %v\nGot code: %v\n", v.wantCode, w.Code)
		}
		// TODO: compare body
	}
}

// TODO: Are there failure cases for ReadAll?
func TestReadAllSets(t *testing.T) {
	testCases := []struct {
		name     string
		data     []*data.Set
		wantCode int
		//		wantResp *bytes.Buffer
	}{
		{
			name: "No params returns 200",
			data: []*data.Set{
				{
					ID:        1,
					Movement:  "Squat",
					Volume:    5,
					Intensity: 80,
				},
			},
			wantCode: 200,
			//wantResp:
		},
	}
	for _, v := range testCases {
		testData := data.NewSetData(v.data)

		ts, err := New(testData)
		if err != nil {
			t.Fail()
		}

		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		//c.Params = v.params

		ts.ReadAll(c)
		log.Printf("ReadAll Result: +%v\n", w)

		if w.Code != v.wantCode {
			// TODO: print body
			t.Fatalf("Wanted code: %v\nGot code: %v\n", v.wantCode, w.Code)
		}
		// TODO: compare body
	}

}
