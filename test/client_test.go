package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"

	//"net/url"
	"testing"

	"github.com/hrand1005/training-notebook/models"
)

const (
	serverURL = "http://localhost:8080"
)

// These tests clients make requests to the test server, which must be running!

func TestUserSignupAndLogin(t *testing.T) {
	// define new HTTP client
	client := newHTTPClientWithCookieJar()

	// define signup post request
	signupBody := bytes.NewBufferString(`{
    "name": "Herb",
    "password": "cookies"
  }`)
	signupReq, err := http.NewRequest(http.MethodPost, serverURL+"/users/signup", signupBody)
	if err != nil {
		t.Fatalf("Failed to build signup request:\nreq: %+v\nerr: %v", signupReq, err)
	}

	// send signup post request
	signupResp, err := client.Do(signupReq)
	if err != nil {
		t.Fatalf("Failed to send signup request:\nreq: %+v\nerr: %v", signupReq, err)
	}
	defer signupResp.Body.Close()

	// signup should yield 201 Created
	if signupResp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201 Created but got %v", signupResp.StatusCode)
	}

	// attempt to login with improper credentials
	invalidLoginBody := bytes.NewBufferString(`{
		"user-id": -1,
		"password": "cookies"
	}`)
	invalidLoginReq, err := http.NewRequest(http.MethodPost, serverURL+"/users/invalidLogin", invalidLoginBody)
	if err != nil {
		t.Fatalf("Failed to build invalidLogin request:\nreq: %+v\nerr: %v", invalidLoginReq, err)
	}

	// send invalidLogin post request
	invalidLoginResp, err := client.Do(invalidLoginReq)
	if err != nil {
		t.Fatalf("Failed to send invalid login request:\nreq: %+v\nerr: %v", invalidLoginReq, err)
	}
	defer invalidLoginResp.Body.Close()

	// invalidLogin should yield 401 Unauthorized
	if invalidLoginResp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected status 404 Unauthorized but got %v", invalidLoginResp.StatusCode)
	}

	// attempt to login with proper credentials
	user := &models.User{}
	if err := DecodeJSON(signupResp.Body, user); err != nil {
		t.Fatalf("Failed to decode signup response to user:\nerr: %v", err)
	}

	loginBody := bytes.NewBufferString(fmt.Sprintf(`{
		"user-id": %v,
		"password": "cookies"
	}`, user.ID))
	loginReq, err := http.NewRequest(http.MethodPost, serverURL+"/users/login", loginBody)
	if err != nil {
		t.Fatalf("Failed to build login request:\nreq: %+v\nerr: %v", loginReq, err)
	}

	// send login post request
	loginResp, err := client.Do(loginReq)
	if err != nil {
		t.Fatalf("Failed to send login request:\nreq: %+v\nerr: %v", loginReq, err)
	}
	defer loginResp.Body.Close()

	// login should yield 200 OK
	if loginResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 OK but got %v", loginResp.StatusCode)
	}

	// TODO: check that cookies are present
	//fmt.Println(loginResp.Cookies())
}

func TestUserPostSet(t *testing.T) {
	// define new HTTP client
	client := newHTTPClientWithCookieJar()

	// define set post request
	setBody := bytes.NewBufferString(`{
		"movement": "Barbell Curl",
		"volume": 1,
		"intensity": 100
  }`)
	setReq, err := http.NewRequest(http.MethodPost, serverURL+"/sets/", setBody)
	if err != nil {
		t.Fatalf("Failed to build set post request:\nreq: %+v\nerr: %v", setReq, err)
	}

	// send set post request before logging in (should fail)
	invalidResp, err := client.Do(setReq)
	if err != nil {
		t.Fatalf("Failed to send set post request:\nreq: %+v\nerr: %v", setReq, err)
	}
	defer invalidResp.Body.Close()

	// should fail with 401 Unauthorized
	if invalidResp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected status 401 Unauthorized but got %v", invalidResp.StatusCode)
	}

	// login with existing user
	userID := LoginWithValidUser(client)

	// re-create earlier set post request
	setBody = bytes.NewBufferString(`{
		"movement": "Barbell Curl",
		"volume": 1,
		"intensity": 100
  }`)
	setReq, err = http.NewRequest(http.MethodPost, serverURL+"/sets/", setBody)
	if err != nil {
		t.Fatalf("Failed to build set post request:\nreq: %+v\nerr: %v", setReq, err)
	}
	// post set with logged in user
	setResp, err := client.Do(setReq)
	if err != nil {
		t.Fatalf("Failed to send set post request:\nreq: %+v\nerr: %v", setReq, err)
	}
	defer setResp.Body.Close()

	if setResp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(setResp.Body)
		t.Fatalf("Expected status 201 Created but got %v\nBody: %s\n", setResp.StatusCode, bodyBytes)
	}

	set := &models.Set{}
	if err := DecodeJSON(setResp.Body, set); err != nil {
		t.Fatalf("Failed to decode set post response to set:\nerr: %v", err)
	}

	if set.UID != userID {
		t.Fatalf("Expected posted set to have UID matching logged in user:\nLogged in user-id: %v\nSet user-id: %v", set.UID, userID)
	}

	// bodyBytes, _ := io.ReadAll(setResp.Body)
	// fmt.Printf("Response body: %v", string(bodyBytes))
	// TODO
	// -- verify that any defined user-id is overwritten with the logged in user's id
}

func TestUserReadSet(t *testing.T) {
	// define HTTP clients to test valid and invalid cases
	clientValid := newHTTPClientWithCookieJar()
	clientInvalid := newHTTPClientWithCookieJar()

	// logs in and creates set with given client, returns setID
	setID := CreateUserAndPostTestSet(clientValid)

	// attempt to read existing set by id without credentials
	endpoint := fmt.Sprintf("%s/sets/%v", serverURL, setID)
	setReq, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		t.Fatalf("Failed to build set read request:\nreq: %+v\nerr: %v", setReq, err)
	}
	invalidResp, err := clientInvalid.Do(setReq)
	if err != nil {
		t.Fatalf("Failed to send read request:\nreq: %+v\nerr: %v", setReq, err)
	}
	defer invalidResp.Body.Close()

	if invalidResp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected response with status code 401 Unauthorized but got %v", invalidResp.StatusCode)
	}

	validResp, err := clientValid.Do(setReq)
	if err != nil {
		t.Fatalf("Failed to send read request:\nreq: %+v\nerr: %v", setReq, err)
	}
	defer validResp.Body.Close()

	if validResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(validResp.Body)
		t.Fatalf("Expected response with status 200 OK but got %v\nBody: %s", validResp.StatusCode, bodyBytes)
	}
	// login with existing user
	// attempt to read set for different user
	// read set for logged in user
}

func TestUserReadSets(t *testing.T) {
	// define HTTP clients to test valid and invalid cases
	clientValid := newHTTPClientWithCookieJar()
	clientInvalid := newHTTPClientWithCookieJar()

	// logs in and creates set with given client, returns setID
	firstSetID := CreateUserAndPostTestSet(clientValid)
	secondSetID := CreateUserAndPostTestSet(clientValid)

	// attempt to read sets without any credentials
	setReq, err := http.NewRequest(http.MethodGet, serverURL+"/sets/", nil)
	if err != nil {
		t.Fatalf("Failed to build set read request:\nreq: %+v\nerr: %v", setReq, err)
	}
	invalidResp, err := clientInvalid.Do(setReq)
	if err != nil {
		t.Fatalf("Failed to send read request:\nreq: %+v\nerr: %v", setReq, err)
	}
	defer invalidResp.Body.Close()

	if invalidResp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected response with status code 401 Unauthorized but got %v", invalidResp.StatusCode)
	}

	validResp, err := clientValid.Do(setReq)
	if err != nil {
		t.Fatalf("Failed to send read request:\nreq: %+v\nerr: %v", setReq, err)
	}
	defer validResp.Body.Close()

	if validResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(validResp.Body)
		t.Fatalf("Expected response with status 200 OK but got %v\nBody: %s", validResp.StatusCode, bodyBytes)
	}

	sets := []*models.Set{}
	err = DecodeJSON(validResp.Body, &sets)
	if err != nil {
		t.Fatalf("Failed to read response into sets: %v", err)
	}

	// verify that each setID is found
	for _, v := range sets {
		if v.ID != firstSetID && v.ID != secondSetID {
			t.Fatalf("Expected sets response to contain posted set ids:\nExpected to find %v and %v\nSets: %+v", firstSetID, secondSetID, sets)
		}
	}

	// login with existing user
	// attempt to read set for different user
	// read set for logged in user
	// define new HTTP client
	// attempt to read sets without credentials
	// login with existing user
	// read sets for logged in user
	// verify that each set belongs to the logged in user

}

/*
func TestUserUpdateSet(t *testing.T) {
  // define new HTTP client
  // attempt to update existing set without credentials
  // login with existing user
  // attempt to update set for different user
  // attempt to update invalid fields
  // update one of the user's sets
}

func TestUserDeleteSet(t *testing.T) {
  // define new HTTP client
  // attempt to delete a set without credentials
  // login with existing user
  // attempt to delete set for different user
  // delete set for logged in user
}
*/

// LoginWithValidUser is a test utility that logs in the given client with a valid
// user without having to rewrite signup/login logic in each use case. It may be assumed
// that the login succeeds upon return. This function returns the user-id of the logged
// in user.
func LoginWithValidUser(c *http.Client) models.UserID {
	const testUserName, testUserPassword = "TestUserName", "TestPassword"
	sBody := bytes.NewBufferString(fmt.Sprintf(`{
    "name": "%s",
    "password": "%s"
  }`, testUserName, testUserPassword))
	sReq, err := http.NewRequest(http.MethodPost, serverURL+"/users/signup", sBody)
	if err != nil {
		panic("LoginWithValidUser: " + err.Error())
	}

	sResp, err := c.Do(sReq)
	if err != nil {
		panic("LoginWithValidUser: " + err.Error())
	}
	defer sResp.Body.Close()

	if sResp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(sResp.Body)
		panic(fmt.Sprintf("LoginWithValidUser: wrong status code in signup response: %v\nBody: %s", sResp.StatusCode, bodyBytes))
	}

	user := &models.User{}
	if err := DecodeJSON(sResp.Body, user); err != nil {
		panic("LoginWithValidUser: " + err.Error())
	}

	lBody := bytes.NewBufferString(fmt.Sprintf(`{
		"user-id": %v,
		"password": "%s"
	}`, user.ID, testUserPassword))
	lReq, err := http.NewRequest(http.MethodPost, serverURL+"/users/login", lBody)
	if err != nil {
		panic("LoginWithValidUser: " + err.Error())
	}

	lResp, err := c.Do(lReq)
	if err != nil {
		panic("LoginWithValidUser: " + err.Error())
	}
	defer lResp.Body.Close()

	// login should yield 201 Created
	if lResp.StatusCode != http.StatusOK {
		panic("LoginWithValidUser: " + err.Error())
	}

	return user.ID
}

func CreateUserAndPostTestSet(c *http.Client) models.SetID {
	LoginWithValidUser(c)

	testSetRequest := bytes.NewBufferString(`{
		"movement": "testMovement",
		"volume": 12,
		"intensity": 65
	}`)

	setReq, err := http.NewRequest(http.MethodPost, serverURL+"/sets/", testSetRequest)
	if err != nil {
		panic("CreateUserAndPostTestSet: " + err.Error())
	}

	setResp, err := c.Do(setReq)
	if err != nil {
		panic("CreateUserAndPostTestSet: Failed to post setRequest: " + err.Error())
	}
	defer setResp.Body.Close()

	if setResp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(setResp.Body)
		panic(fmt.Sprintf("CreateUserAndPostTestSet: Expected 201 status created but got %v\nBody: %s", setResp.StatusCode, bodyBytes))
	}

	set := &models.Set{}
	err = DecodeJSON(setResp.Body, set)
	if err != nil {
		panic("CreateUserAndPostTestSet: " + err.Error())
	}

	return set.ID
}

// DecodeJSON decodes the given Reader to the target type
func DecodeJSON(source io.Reader, target interface{}) error {
	d := json.NewDecoder(source)
	return d.Decode(target)
}

func newHTTPClientWithCookieJar() *http.Client {
	jar, _ := cookiejar.New(nil)
	return &http.Client{
		Jar: jar,
	}
}
