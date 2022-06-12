package users

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
	"github.com/hrand1005/training-notebook/models"
)

func TestLogin(t *testing.T) {
  testUser := &models.User{
    ID: 2,
    Name: "TestUser",
    Password: hashPassword("12345"),
  }
  tests := []struct{
    name string
		requestBody bytes.Buffer
		db          *data.MockUserDB
		wantCode    int
		wantToken string
  }{
    {
      name: "Valid request and DB call returns StatusOK",
			requestBody: *bytes.NewBufferString(` {
				"user-id": 2,
				"password": "12345"
			} `),
			db: &data.MockUserDB{
        UserByIDStub: func(id models.UserID) (*models.User, error) {
          return testUser, nil
        },
			},
			wantCode: http.StatusOK,
			wantToken: buildToken(testUser),
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

		// set the body in the test context's Request
		bodyReader := bytes.NewReader(v.requestBody.Bytes())
		// method/uri parsing exceed the scope of this test
		c.Request, _ = http.NewRequest("", "", bodyReader)

		// execute create with the test context
		u.Login(c)

    fmt.Printf("Wanted token: %s\nGot response: %v", v.wantToken, w.Body)
		// check response code
		if v.wantCode != w.Code {
			t.Fatalf("Wanted code: %v\nGot code: %v\n", v.wantCode, w.Code)
		}

    fmt.Printf("Wanted token: %s\nGot response: %v", v.wantToken, w.Body)
  }
}
