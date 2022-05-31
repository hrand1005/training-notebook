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

// TestDeleteSet tests the API layer's Delete method for the Sets resource.
// The test suite mocks the SetDB interface to test edge cases and error conditions.
func TestDeleteSet(t *testing.T) {
	tests := []struct {
		name     string
		db       *data.MockSetDB
		params   gin.Params
		wantCode int
		wantResp bytes.Buffer
	}{
		{
			name: "Set found with valid db call returns StatusNoContent",
			db: &data.MockSetDB{
				DeleteSetStub: func(id data.SetID) error {
					return nil
				},
			},
			params:   []gin.Param{{Key: "id", Value: "1"}},
			wantCode: http.StatusNoContent,
		},
		{
			name: "Set not found returns StatusNotFound",
			db: &data.MockSetDB{
				DeleteSetStub: func(id data.SetID) error {
					return data.ErrNotFound
				},
			},
			params:   []gin.Param{{Key: "id", Value: "4"}},
			wantCode: http.StatusNotFound,
			wantResp: *bytes.NewBufferString(`{
				"message": "no such set with id 4"
			}`),
		},
		{
			name: "Invalid db query returns InternalServerError",
			db: &data.MockSetDB{
				DeleteSetStub: func(id data.SetID) error {
					return fmt.Errorf("Expected error")
				},
			},
			params:   []gin.Param{{Key: "id", Value: "4"}},
			wantCode: http.StatusInternalServerError,
			wantResp: *bytes.NewBufferString(`{
				"message": "Expected error"
			}`),
		},
		{
			name:     "Invalid params returns StatusBadRequest",
			params:   []gin.Param{{Key: "bad", Value: "request"}},
			wantCode: http.StatusBadRequest,
			wantResp: *bytes.NewBufferString(fmt.Sprintf(`{
				"message": %q
			}`, ErrInvalidSetID)),
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
		c.Params = v.params

		// execute Delete on test context
		ts.Delete(c)

		// check response code
		if v.wantCode != w.Code {
			t.Fatalf("Wanted code: %v\nGot code: %v\n", v.wantCode, w.Code)
		}

		if v.wantCode != http.StatusNoContent {
			// check response body
			if equal, _ := JSONBytesEqual(v.wantResp.Bytes(), w.Body.Bytes()); !equal {
				t.Fatalf("Wanted body: %v\nGot body: %v\n", v.wantResp.String(), w.Body.String())
			}
		}
	}
}
