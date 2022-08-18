Feature: Sets

  Scenario: Create new set
    When the client sends "POST" request to endpoint "/sets" with body:
    """
    {
      "data": {
        "type": "set",
        "attributes": {
          "movement": "curl", 
          "volume": 12,
          "intensity": 65.0
        }
      }
    }
    """
    Then the server responds with status code "204"

  Scenario: Read existing set 
    Given "set-id" exists
    When the client sends "GET" request to endpoint "/sets/set-id"
    Then the server responds with status code "200" 

