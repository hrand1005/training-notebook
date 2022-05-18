package sets

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
)

// TODO: setup mock db
func TestCreateSet(t *testing.T) {

	// test cases:
	// Bad Request
	// success case

	testCases := []struct {
		name     string
		wantCode int
		set      data.Set
		//		wantResp *bytes.Buffer
	}{
		{
			name:     "Set created returns 200",
			wantCode: 201,
			set: data.Set{
				Movement:  "Barbell Curl",
				Volume:    12,
				Intensity: 0.5,
			},
			//wantResp:
		},
		// TODO: Bad request case
	}

	for _, v := range testCases {
		// start with empty setData
		testData := data.NewSetData(nil)

		ts, err := New(testData)
		if err != nil {
			t.Fail()
		}

		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)

		c.Set("newSet", v.set)

		ts.Create(c)

		if w.Code != v.wantCode {
			// TODO: print body
			t.Fatalf("Wanted code: %v\nGot code: %v\n", v.wantCode, w.Code)
		}
		// TODO: compare body
		// TODO: ADD CLEANUP!!!
	}
}
