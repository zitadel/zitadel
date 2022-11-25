Feature: Users can login via IDP

    Testing user login feature for external accounts
    Out of scope: Pre-login Self-service (register, password reset)

    Rule: Users can login with external idp OIDC - Google
      Background:
        Given A user with a linked Google  does exist
        And an application with redirect uri "http://localhost:8080/redirect-to-here" exists

      Example: Username with External IDP Google
        Given login policy has values '{"second_factors": [], "passwordless_type": "PASSWORDLESS_TYPE_NOT_ALLOWED", "mfa_init_skip_lifetime": 0, "allow_external_idp": true, "idps": [{"idp_id": "1", "name": "Google", "idp_type": "2"}]}'
        And a clear browser session
        And user navigates to authorize endpoint with redirect uri "http://localhost:8080/redirect-to-here"

        When user enters loginname
        And user authenticates in Google

        Then user is redirected to "http://localhost:8080/redirect-to-here"

    Rule: Users can login with external idp OIDC - Google after auto registration
      Background:
        Given A user with password does exist
        And an application with redirect uri "http://localhost:8080/redirect-to-here" exists

      Example: Username with External IDP Google Auto register
        Given login policy has values '{"second_factors": [], "passwordless_type": "PASSWORDLESS_TYPE_NOT_ALLOWED", "mfa_init_skip_lifetime": 0, "allow_external_idp": true, "idps": [{"idp_id": "1", "name": "Google", "idp_type": "2"}]}'
        And registered idp {"auto_register": true}
        And a clear browser session
        And user navigates to authorize endpoint with redirect uri "http://localhost:8080/redirect-to-here"

        When user enters loginname
        And user authenticates in Google

        Then user is created in ZITADEL
        And redirected to "http://localhost:8080/redirect-to-here"

    Rule: Users can login with external idp OIDC - Google after manually register
      Background:
        Given A user with password does exist
        And an application with redirect uri "http://localhost:8080/redirect-to-here" exists

      Example: Username with External IDP Google manual registration
        Given login policy has values '{"second_factors": [], "passwordless_type": "PASSWORDLESS_TYPE_NOT_ALLOWED", "mfa_init_skip_lifetime": 0, "allow_external_idp": true, "idps": [{"idp_id": "1", "name": "Google", "idp_type": "2"}]}'
        And registered idp {"auto_register": false}
        And a clear browser session
        And user navigates to authorize endpoint with redirect uri "http://localhost:8080/redirect-to-here"

        When user enters loginname
        And user authenticates in Google
        And chooses registration
        And fill register form

        Then user is created in ZITADEL
        And redirected to "http://localhost:8080/redirect-to-here"

    Rule: Users can login with external idp OIDC - Google after linking
      Background:
        Given A user with password "Password1" does exist
        And an application with redirect uri "http://localhost:8080/redirect-to-here" exists

      Example: Username with External IDP Google
        Given login policy has values '{"second_factors": [], "passwordless_type": "PASSWORDLESS_TYPE_NOT_ALLOWED", "mfa_init_skip_lifetime": 0, "allow_external_idp": true, "idps": [{"idp_id": "1", "name": "Google", "idp_type": "2"}]}'
        And registered idp {"auto_register": false}
        And a clear browser session
        And user navigates to authorize endpoint with redirect uri "http://localhost:8080/redirect-to-here"

        When user enters loginname
        And user authenticates in Google
        And chooses linking
        And user enters username
        And user enters password

        Then user is linked
        And redirected to "http://localhost:8080/redirect-to-here"

  Rule: Users can login with external idp OIDC - AzureAD after auto registration
    Background:
      Given A user with password does exist
      And an application with redirect uri "http://localhost:8080/redirect-to-here" exists

    Example: Username with External IDP Google Auto register
      Given login policy has values '{"second_factors": [], "passwordless_type": "PASSWORDLESS_TYPE_NOT_ALLOWED", "mfa_init_skip_lifetime": 0, "allow_external_idp": true, "idps": [{"idp_id": "1", "name": "AzureAD", "idp_type": "2"}]}'
      And registered idp {"auto_register": true}
      And a clear browser session
      And user navigates to authorize endpoint with redirect uri "http://localhost:8080/redirect-to-here"

      When user enters loginname
      And user authenticates in AzureAD

      Then user is created in ZITADEL
      And redirected to "http://localhost:8080/redirect-to-here"
