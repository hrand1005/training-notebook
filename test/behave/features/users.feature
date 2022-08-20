Feature: Users

  Scenario: 
    Given a user exists with id "_user_id"
    But the client is not authenticated
    When the client sends GET request to "/users/_user_id"
    Then the server responds with status code "403"

  Scenario:
    Given a user exists with id "_user_id"
    But the authenticated client is not that user
    When the client sends GET request to "/users/_user_id"
    Then the server responds with status code "403"

  Scenario:
    Given a user exists with id "_user_id"
    And the authenticated client is that user
    When the client sends GET request to "/users/_user_id"
    Then the server responds with status code "200"

