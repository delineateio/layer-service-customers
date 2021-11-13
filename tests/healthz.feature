Feature: Healthz

Scenario: health
  Given path /healthz
  When method GET
  Then status 200
