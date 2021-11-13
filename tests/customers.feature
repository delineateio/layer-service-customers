Feature: Customers

Scenario: customers
  Given path /customers
  When method GET
  Then status 200

Scenario: customers_again
  Given path /customers
  When method GET
  Then status 200
