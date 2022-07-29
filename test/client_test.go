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
	serverAPI = "http://localhost:8080/api"
)

var (
	testUser = &models.User{
		Name:     "DummyUserName",
		Password: "DummyPassword",
	}
	testSet = &models.Set{
		Movement:  "Big Squat",
		Volume:    20,
		Intensity: 70,
	}
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
	signupReq, err := http.NewRequest(http.MethodPost, serverAPI+"/signup", signupBody)
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
	invalidLoginReq, err := http.NewRequest(http.MethodPost, serverAPI+"/login", invalidLoginBody)
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
	loginReq, err := http.NewRequest(http.MethodPost, serverAPI+"/login", loginBody)
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
	setReq, err := http.NewRequest(http.MethodPost, serverAPI+"/sets/", setBody)
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
	userID := CreateUserAndLogin(client, testUser)

	// re-create earlier set post request
	setBody = bytes.NewBufferString(`{
		"movement": "Barbell Curl",
		"volume": 1,
		"intensity": 100
  }`)
	setReq, err = http.NewRequest(http.MethodPost, serverAPI+"/sets/", setBody)
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

	// TODO
	// -- verify that any defined user-id is overwritten with the logged in user's id
}

func TestUserReadSet(t *testing.T) {
	// define HTTP clients to test valid and invalid cases
	clientValid := newHTTPClientWithCookieJar()
	clientInvalid := newHTTPClientWithCookieJar()

	// logs in and creates set with given client
	userID := CreateUserAndLogin(clientValid, testUser)
	// posts the test set using the logged in client
	setID := CreateAndPostSet(clientValid, testSet)

	// attempt to read existing set by id without credentials on the invalid client
	endpoint := fmt.Sprintf("%s/sets/%v", serverAPI, setID)
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

	set := &models.Set{}
	err = DecodeJSON(validResp.Body, set)
	if err != nil {
		t.Fatalf("Failed to decode set read response to set:\nerr: %v", err)
	}

	if userID != set.UID {
		t.Fatalf("Expected to retrieve a set with user-id %v but got %v.", userID, set.UID)
	}

	if setID != set.ID {
		t.Fatalf("Expected to retrieve a set with set-id %v but got %v.", setID, set.ID)
	}
	// login with existing user
	// attempt to read set for different user
	// read set for logged in user
}

func TestUserReadSets(t *testing.T) {
	// define HTTP client for invalid use cases
	clientInvalid := newHTTPClientWithCookieJar()

	// attempt to read sets without any credentials
	setReq, err := http.NewRequest(http.MethodGet, serverAPI+"/sets/", nil)
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

	// create and login with a user and get validUserID for success case
	clientValid := newHTTPClientWithCookieJar()
	validUserID := CreateUserAndLogin(clientValid, testUser)
	firstSetID := CreateAndPostSet(clientValid, testSet)
	secondSetID := CreateAndPostSet(clientValid, testSet)

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

	if len(sets) != 2 {
		t.Fatalf("Expected sets response to contain 2 sets.\nResponse sets count: %v\nSets: %+v", len(sets), sets)
	}

	// verify that each set retrieved is as expected
	if firstSetID != sets[0].ID && firstSetID != sets[1].ID {
		t.Fatalf("Expected one set in response to match set-id %v", firstSetID)
	}
	if secondSetID != sets[0].ID && secondSetID != sets[1].ID {
		t.Fatalf("Expected one set in response to match set-id %v", secondSetID)
	}
	for _, v := range sets {
		if firstSetID != v.ID {
			if secondSetID != v.ID {
				t.Fatalf("Expected set in response to match one of the following set-ids:\n%v, %v\nGot set-id: %v", firstSetID, secondSetID, v.ID)
			}
		}
		if v.UID != validUserID {
			t.Fatalf("Retrieved set doesn't match id of logged in user:\nExpected user-id: %v\nGot set with user-id: %v", validUserID, v.ID)
		}
	}

	// verify that the valid client's data can't be retrieved by another logged in user
	CreateUserAndLogin(clientInvalid, testUser)

	setReq, err = http.NewRequest(http.MethodGet, serverAPI+"/sets/", nil)
	if err != nil {
		t.Fatalf("Failed to build set read request:\nreq: %+v\nerr: %v", setReq, err)
	}
	invalidResp, err = clientInvalid.Do(setReq)
	if err != nil {
		t.Fatalf("Failed to send read request:\nreq: %+v\nerr: %v", setReq, err)
	}
	defer invalidResp.Body.Close()

	// we expect this user to retrieve an empty list of sets
	invalidSets := []*models.Set{}
	DecodeJSON(invalidResp.Body, invalidSets)

	if len(invalidSets) != 0 {
		t.Fatalf("Expected invalid case to retrieve no sets but got set response:\n%+v", invalidSets)
	}
}

func TestUserUpdateSet(t *testing.T) {
	client := newHTTPClientWithCookieJar()

	// login with user and create a set
	userID := CreateUserAndLogin(client, testUser)
	setID := CreateAndPostSet(client, testSet)

	// define new set fields to update the existing set
	updateSet := &models.Set{
		ID:        setID,
		UID:       userID,
		Movement:  "Update Movement",
		Volume:    1,
		Intensity: 1,
	}
	setBody := bytes.NewBufferString(fmt.Sprintf(`{
  	"movement": "%s",
  	"volume": %v,
  	"intensity": %v
  }`, updateSet.Movement, updateSet.Volume, updateSet.Intensity))

	// update the set with new fields for the logged in user
	endpoint := fmt.Sprintf("%s/sets/%v", serverAPI, setID)
	setReq, err := http.NewRequest(http.MethodPut, endpoint, setBody)
	if err != nil {
		t.Fatalf("Failed to build set request, err: %v", err)
	}
	setResp, err := client.Do(setReq)
	if err != nil {
		t.Fatalf("Failed to send set update request, err: %v", err)
	}

	if setResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(setResp.Body)
		t.Fatalf("Expected 200 Status OK but got %v\nSet Resp: %s", setResp.StatusCode, bodyBytes)
	}

	gotSet := &models.Set{}
	DecodeJSON(setResp.Body, gotSet)

	if gotSet.UID != userID || gotSet.Movement != updateSet.Movement || gotSet.Volume != updateSet.Volume || gotSet.Intensity != updateSet.Intensity {
		t.Fatalf("Fields not updated as expected:\nGot Set: %+v\nExpected set: %+v", gotSet, updateSet)
	}

	// attempt to update set for different user
	clientInvalid := newHTTPClientWithCookieJar()
	CreateUserAndLogin(clientInvalid, testUser)
	invalidBody := bytes.NewBufferString(fmt.Sprintf(`{
  	"movement": "%s",
  	"volume": %v,
  	"intensity": %v
  }`, updateSet.Movement, updateSet.Volume, updateSet.Intensity))
	invalidReq, err := http.NewRequest(http.MethodPut, endpoint, invalidBody)
	if err != nil {
		t.Fatalf("Failed to build set request, err: %v", err)
	}
	invalidResp, err := clientInvalid.Do(invalidReq)
	if err != nil {
		t.Fatalf("Failed to send set update request, err: %v", err)
	}

	if invalidResp.StatusCode != http.StatusNotFound {
		bodyBytes, _ := io.ReadAll(invalidResp.Body)
		t.Fatalf("Expected 404 Not Found for invalid user but got %v\nResp: %s", invalidResp.StatusCode, bodyBytes)
	}
	// TODO:
	// attempt to update invalid fields
}

func TestUserDeleteSet(t *testing.T) {
	// define HTTP clients to test valid and invalid cases
	clientValid := newHTTPClientWithCookieJar()

	// logs in and creates set with given client
	CreateUserAndLogin(clientValid, testUser)
	// posts the test set using the logged in client
	setID := CreateAndPostSet(clientValid, testSet)

	// attempt to read existing set by id without credentials with invalid client
	endpoint := fmt.Sprintf("%s/sets/%v", serverAPI, setID)
	setReq, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		t.Fatalf("Failed to build set read request:\nreq: %+v\nerr: %v", setReq, err)
	}
	clientInvalid := newHTTPClientWithCookieJar()
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

	if validResp.StatusCode != http.StatusNoContent {
		bodyBytes, _ := io.ReadAll(validResp.Body)
		t.Fatalf("Expected response with status 204 NoContent but got %v\nBody: %s", validResp.StatusCode, bodyBytes)
	}

	// attempting to delete the set again should result in 404
	validResp, err = clientValid.Do(setReq)
	if err != nil {
		t.Fatalf("Failed to send read request:\nreq: %+v\nerr: %v", setReq, err)
	}
	defer validResp.Body.Close()

	if validResp.StatusCode != http.StatusNotFound {
		bodyBytes, _ := io.ReadAll(validResp.Body)
		t.Fatalf("Expected response with status 404 Not Found but got %v\nBody: %s", validResp.StatusCode, bodyBytes)
	}
}

// CreateUserAndLogin is a testing utility which creates the logged in user on the given client, returning the
// UserID of the created user, and logging in (thus setting jwt cookie enabling authenticated operations)
func CreateUserAndLogin(c *http.Client, u *models.User) models.UserID {
	sBody := bytes.NewBufferString(fmt.Sprintf(`{
		"name": "%s",
		"password": "%s"
	}`, u.Name, u.Password))
	sReq, err := http.NewRequest(http.MethodPost, serverAPI+"/signup", sBody)
	if err != nil {
		panic("CreateUserAndLogin: " + err.Error())
	}

	sResp, err := c.Do(sReq)
	if err != nil {
		panic("CreateUserAndLogin: " + err.Error())
	}
	defer sResp.Body.Close()

	if sResp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(sResp.Body)
		panic(fmt.Sprintf("CreateUserAndLogin: wrong status code in signup response: %v\nBody: %s", sResp.StatusCode, bodyBytes))
	}

	user := &models.User{}
	if err := DecodeJSON(sResp.Body, user); err != nil {
		panic("CreateUserAndLogin: " + err.Error())
	}

	lBody := bytes.NewBufferString(fmt.Sprintf(`{
		"user-id": %v,
		"password": "%s"
	}`, user.ID, u.Password))
	lReq, err := http.NewRequest(http.MethodPost, serverAPI+"/login", lBody)
	if err != nil {
		panic("CreateUserAndLogin: " + err.Error())
	}

	lResp, err := c.Do(lReq)
	if err != nil {
		panic("CreateUserAndLogin: " + err.Error())
	}
	defer lResp.Body.Close()

	// login should yield 201 Created
	if lResp.StatusCode != http.StatusOK {
		panic("CreateUserAndLogin: " + err.Error())
	}

	return user.ID
}

// PostTestSet posts the given set using the given client's credentials.
// Returns the id of the created set. Panics upon failure.
func CreateAndPostSet(c *http.Client, s *models.Set) models.SetID {
	setBody := bytes.NewBufferString(fmt.Sprintf(`{
		"movement": "%s",
		"volume": %v,
		"intensity": %v
	}`, s.Movement, s.Volume, s.Intensity))

	setReq, err := http.NewRequest(http.MethodPost, serverAPI+"/sets/", setBody)
	if err != nil {
		panic("CreateAndPostSet: " + err.Error())
	}

	setResp, err := c.Do(setReq)
	if err != nil {
		panic("CreateAndPostSet: Failed to post setRequest: " + err.Error())
	}
	defer setResp.Body.Close()

	if setResp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(setResp.Body)
		panic(fmt.Sprintf("CreateAndPostSet: Expected 201 status created but got %v\nBody: %s", setResp.StatusCode, bodyBytes))
	}

	set := &models.Set{}
	err = DecodeJSON(setResp.Body, set)
	if err != nil {
		panic("CreateAndPostSet: " + err.Error())
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
