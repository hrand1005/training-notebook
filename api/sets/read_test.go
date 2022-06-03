package sets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
	"github.com/hrand1005/training-notebook/models"
)

// TestReadAllSets tests the API layer's Read method for the Sets resource.
// The test suite mocks the SetDB interface to test edge cases and error conditions.
func TestReadSingleSet(t *testing.T) {
	tests := []struct {
		name     string
		db       *data.MockSetDB
		id       string
		wantCode int
		wantResp bytes.Buffer
	}{
		{
			name: "Set found with valid db call returns StatusOK",
			db: &data.MockSetDB{
				SetByIDStub: func(s models.SetID) (*models.Set, error) {
					return &models.Set{
						ID:        1,
						Movement:  "Squat",
						Volume:    5,
						Intensity: 80,
					}, nil
				},
			},
			id:       "1",
			wantCode: http.StatusOK,
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
			name: "Set not found returns StatusNotFound",
			db: &data.MockSetDB{
				SetByIDStub: func(s models.SetID) (*models.Set, error) {
					return nil, data.ErrNotFound
				},
			},
			id:       "4",
			wantCode: http.StatusNotFound,
			wantResp: *bytes.NewBufferString(`{
				"message": "no such set with id 4"
			}`),
		},
		{
			name: "Invalid db query returns InternalServerError",
			db: &data.MockSetDB{
				SetByIDStub: func(s models.SetID) (*models.Set, error) {
					return nil, fmt.Errorf("Expected error")
				},
			},
			id:       "4",
			wantCode: http.StatusInternalServerError,
			wantResp: *bytes.NewBufferString(`{
				"message": "Expected error"
			}`),
		},
		{
			name:     "Invalid params returns StatusBadRequest",
			id:       "-1",
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
		c.AddParam("id", v.id)

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

// TestReadAllSets tests the API layer's ReadAll method for the Sets resource.
// The test suite mocks the SetDB interface to test edge cases and error conditions.
func TestReadAllSets(t *testing.T) {
	tests := []struct {
		name     string
		db       *data.MockSetDB
		wantCode int
		wantResp bytes.Buffer
	}{
		{
			name: "Valid db call with multiple sets returns StatusOK",
			db: &data.MockSetDB{
				SetsStub: func() ([]*models.Set, error) {
					return []*models.Set{
						{
							ID:        1,
							Movement:  "Squat",
							Volume:    5,
							Intensity: 80,
						},
						{
							ID:        2,
							Movement:  "Deadlift",
							Volume:    4,
							Intensity: 85,
						},
					}, nil
				},
			},
			wantCode: http.StatusOK,
			wantResp: *bytes.NewBufferString(`[
				{
					"id": 1,
					"movement": "Squat",
					"volume": 5,
					"intensity": 80
				},
				{
					"id": 2,
					"movement": "Deadlift",
					"volume": 4,
					"intensity": 85
				}
			]`),
		},
		{
			name: "Valid db call with empty set returns StatusOK",
			db: &data.MockSetDB{
				SetsStub: func() ([]*models.Set, error) {
					return nil, nil
				},
			},
			wantCode: http.StatusOK,
			wantResp: *bytes.NewBufferString(`[]`),
		},
		{
			name: "Invalid db call returns InternalServerError",
			db: &data.MockSetDB{
				SetsStub: func() ([]*models.Set, error) {
					return nil, fmt.Errorf("Expected Error")
				},
			},
			wantCode: http.StatusInternalServerError,
			wantResp: *bytes.NewBufferString(`{
				"message": "Expected Error"
			}`),
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
		log.Printf("Problem unmarshalling json a: %v\nError: %v\n", a, err)
		return false, err
	}
	if err := json.Unmarshal(b, &j2); err != nil {
		log.Printf("Problem unmarshalling json b: %v\nError: %v\n", b, err)
		return false, err
	}
	return reflect.DeepEqual(j2, j), nil
}
