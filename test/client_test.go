package test

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"
)

const (
	serverURL = "http://localhost:8080"
)

// These tests clients make requests to the test server, which must be running!
// From project root:
//
// ./scripts/dev.sh configs/test_config.yaml

func TestUserSignupAndLogin(t *testing.T) {
	// define new HTTP client
	client := http.Client{}

	// define signup post request
	signupBody := bytes.NewBufferString(`{
    "name": "Herb",
    "password": "cookies"
  }`)
	req, err := http.NewRequest(http.MethodPost, serverURL+"/users/signup", signupBody)
	if err != nil {
		t.Fatalf("Failed to build HTTP request:\nreq: %+v\nerr: %v", req, err)
	}

	time.Sleep(time.Second * 5)

	// send signup post request
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request:\nreq: %+v\nerr: %v", req, err)
	}
	fmt.Printf("Resp:\n%+v\n", resp)
	// attempt to login with improper credentials
	// attempt to login with proper credentials
}

/*
func TestUserPostSet(t *testing.T) {
  // define new HTTP client
  // attempt to post set without credentials
  // login with existing user
  // post set with logged in user
}

func TestUserReadSet(t *testing.T) {
  // define new HTTP client
  // attempt to read existing set by id without credentials
  // login with existing user
  // attempt to read set for different user
  // read set for logged in user
}

func TestUserReadSets(t *testing.T) {
  // define new HTTP client
  // attempt to read sets without credentials
  // login with existing user
  // read sets for logged in user
  // verify that each set belongs to the logged in user
}

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
