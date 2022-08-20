Feature: sets

  Scenario: 
    Given a set exists with id "_set_id"
    But the client is not authenticated
    When the client sends GET request to "/sets/_set_id"
    Then the server responds with status code "403"

  Scenario:
    Given no set exists with id "_set_id"
    But the client is authenticated
    When the client sends GET request to "/sets/_set_id"
    Then the server responds with status code "404"

  Scenario:
    Given a set exists with id "_set_id"
    But the authenticated client is not the owner of the set
    When the client sends GET request to "/sets/_set_id"
    Then the server responds with status code "403"

  Scenario:
    Given a set exists with id "_set_id"
    And the authenticated client is the owner of the set
    When the client sends GET request to "/sets/_set_id"
    Then the server responds with status code "200"
