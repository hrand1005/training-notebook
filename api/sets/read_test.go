package sets

import (
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestReadSingleSet(t *testing.T) {

	// test cases:
	// id not found
	// invalid params or smth
	// success case

	testCases := []struct {
		name     string
		params   gin.Params
		wantCode int
		//		wantResp *bytes.Buffer
	}{
		{
			name:     "Set found returns 200",
			params:   []gin.Param{{Key: "id", Value: "1"}},
			wantCode: 200,
			//wantResp:
		},
		{
			name:     "Set not found returns 404",
			params:   []gin.Param{{Key: "id", Value: "2"}},
			wantCode: 404,
			//wantResp:
		},
		{
			name:     "Invalid params returns 400",
			params:   []gin.Param{{Key: "bad", Value: "request"}},
			wantCode: 400,
			//wantResp:
		},
	}

	for _, v := range testCases {
		ts, err := New()
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
		wantCode int
		//		wantResp *bytes.Buffer
	}{
		{
			name:     "No params returns 200",
			wantCode: 200,
			//wantResp:
		},
	}
	for _, v := range testCases {
		ts, err := New()
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
