![ZITADEL](/docs/img/zitadel-logo-oneline-lightdesign@2x.png "Zitadel Logo")

***
# Before you begin
This documentation has been carefully compiled to accommodate the requirements of three personas.

| Persona | Role |
|:-|:-|
| User | Seeing Zitadel not as a product per se but as a means of login-experience. … |
| Administrator | Creating new logins, groups, and grants. … |
| Developer | Extending or fixing Zitadel. … |

When navigating the documentation you should always be able to find all information needed to accomplish the tasks of one of the personas—it should be “persona-complete”. This is also our statement and promise, and should be granted by people who are willing to contribute to the documentation.

***
# Getting started
## Features

***
# Using Zitadel

## I’m a user
* QuickStart
* Creating a new login
* Login
* Resetting my password
* Getting around most common problems
    * I forgot my password, what can I do?
    * I cannot join my IDP with my Zitadel account, what is wrong?
    * I don’t receive the tokens for 2FA (Two Factor Authentication), where did it get stuck?
## I’m an administrator
* Configuring Zitadel
## I’m a developer
* Find my way around the repository
* Things to keep in mind when cloning the repository
    * Lead developers
    * Contributors
    * Sponsors
* Becoming part of the community
* Pushing versus creating pull requests
* Testing you code

***
# Learning Zitadel
## I’m a user
## I’m an administrator
### How to deploy Zitadel
* Cloud Service
* On Premise
### How to integrate Zitadel
### How to migrate
* Data model
* Configuration
* Files
## I’m a developer
* The data model
* The Zitadel Event Stream
* Creating connectors
* Configuring the dashboard
* Extending the data model
* Using your favourite programming language

***
# Troubleshooting Zitadel

***
# Learning More
```swift
import Foundation
import SwiftNIO
enum ZitadelComError {
    case .OK
    case .ServerError(Error)
    case .NetworkError(Error)
}
```

***
# FAQ

