Feature: Users can login via ZITADEL

    Testing user login feature for local accounts
    Out of scope: Pre-login Self-service (register, password reset)

    Rule: Users can login via username and password

        Background:
            Given A user with password "Password1!" and verified email does exist
            And an application with redirect uri "http://localhost:8080/redirect-to-here" exists

        Example: Username Password
            Given login policy has values '{"second_factors": [], "passwordless_type": "PASSWORDLESS_TYPE_NOT_ALLOWED", "mfa_init_skip_lifetime": 0}'
            And a clear browser session
            And user navigates to authorize endpoint with redirect uri "http://localhost:8080/redirect-to-here"

            When user enters loginname
            And user enters password "Password1!"

            Then user is redirected to "http://localhost:8080/redirect-to-here"