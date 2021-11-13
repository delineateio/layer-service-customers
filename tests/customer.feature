Feature: Customer

Scenario: wrong_verb
  Given path /customer
  When method GET
  Then status 404

Scenario: no_body
  Given path /customer
  When method POST
  Then status 400

Scenario: bad_request
  Given path /customer
    And request customer/bad_request.json
  When method POST
  Then status 400

Scenario: good_request
  Given path /customer
    And request customer/good_request.json
  When method POST
  Then status 201

Scenario: duplicate_request
  Given path /customer
    And request customer/good_request.json
  When method POST
  Then status 409
