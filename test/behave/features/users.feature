Feature: Users

  Scenario: Read existing user
    Given "user-id" exists
    When the client sends "GET" request to "/users/user-id"
    Then the server responds with status code "200" 

