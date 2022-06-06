package users

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

// TestReadAllUsers tests the API layer's Read method for the Users resource.
// The test suite mocks the UserDB interface to test edge cases and error conditions.
func TestReadSingleUser(t *testing.T) {
	tests := []struct {
		name     string
		db       *data.MockUserDB
		id       string
		wantCode int
		wantResp bytes.Buffer
	}{
		{
			name: "User found with valid db call returns StatusOK",
			db: &data.MockUserDB{
				UserByIDStub: func(s models.UserID) (*models.User, error) {
					return &models.User{
						ID:   1,
						Name: "Ray",
					}, nil
				},
			},
			id:       "1",
			wantCode: http.StatusOK,
			wantResp: *bytes.NewBufferString(`
				{
					"id": 1,
					"name": "Ray"
				}
			`),
		},
		{
			name: "User not found returns StatusNotFound",
			db: &data.MockUserDB{
				UserByIDStub: func(s models.UserID) (*models.User, error) {
					return nil, data.ErrNotFound
				},
			},
			id:       "4",
			wantCode: http.StatusNotFound,
			wantResp: *bytes.NewBufferString(`{
				"message": "no such user with id 4"
			}`),
		},
		{
			name: "Invalid db query returns InternalServerError",
			db: &data.MockUserDB{
				UserByIDStub: func(s models.UserID) (*models.User, error) {
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
			}`, ErrInvalidUserID)),
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

// TestReadAllUsers tests the API layer's Read method for the Users resource.
// The test suite mocks the UserDB interface to test edge cases and error conditions.
func TestReadAllUsers(t *testing.T) {
	tests := []struct {
		name     string
		db       *data.MockUserDB
		wantCode int
		wantResp *bytes.Buffer
	}{
		{
			name: "Valid db call with multiple users returns StatusOK",
			db: &data.MockUserDB{
				UsersStub: func() ([]*models.User, error) {
					return []*models.User{
						{
							ID:   1,
							Name: "nendo",
						},
						{
							ID:   2,
							Name: "yert",
						},
					}, nil
				},
			},
			wantCode: http.StatusOK,
			wantResp: bytes.NewBufferString(`[
				{
					"id": 1,
					"name": "nendo"
				},
				{
					"id": 2,
					"name": "yert"
				}
			]`),
		},
	}
	for _, v := range tests {
		// configure test case with data and test context
		u, err := New(v.db)
		if err != nil {
			t.Fail()
		}

		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// execute ReadAll with test context
		u.ReadAll(c)

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
