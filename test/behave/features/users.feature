Feature: Users

  Background: 
    Given a valid user exists with id "1"

  Scenario: 
    Given the client is not authenticated
    When the client sends request "GET /users/1"
    Then the server responds with status code "403"

  Scenario:
    Given the client is authenticated 
    And the client is not the owner of the resource
    When the client sends request "GET /users/1"
    Then the server responds with status code "403"

  Scenario:
    Given the client is authenticated 
    And the client is the owner of the resource
    When the client sends request "GET /users/1"
    Then the server responds with status code "200"

