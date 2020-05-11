# go-credentials
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/engi-fyi/go-credentials)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/engi-fyi/go-credentials)
![License](https://img.shields.io/github/license/engi-fyi/go-credentials)
<br />
![Build](https://github.com/engi-fyi/go-credentials/workflows/Build/badge.svg)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=engi-fyi_go-credentials&metric=alert_status)](https://sonarcloud.io/dashboard?id=engi-fyi_go-credentials)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=engi-fyi_go-credentials&metric=coverage)](https://sonarcloud.io/dashboard?id=engi-fyi_go-credentials)
[![Go Report Card](https://goreportcard.com/badge/github.com/engi-fyi/go-credentials)](https://goreportcard.com/report/github.com/engi-fyi/go-credentials)

Welcome to `go-credentials`!

This project is being built out of a need for a library to manage credentials files (similar to AWS credentials), their attributes, sessions, and environment variables.

![go-credentials logo](https://github.com/engi-fyi/go-credentials/raw/master/assets/go-credentials-logo.png)

The Credential API is broken down into two pieces, each with their own functionality:
1. `Factory`: responsible for setting variables that are global to your application, and;
    - Set alternate keys for username/password (e.g. ACCESS_TOKEN/SECRET_KEY).
    - Set the output type of the credentials (only `.ini` currently implemented).
2. `Credential`: represents a user's credentials.
    - Username/Password defined on model.
    - Set Attributes.
    - Get Attributes.
    - Read values from Environment variables.
    - Load from .ini File.
    
To get started, is all you need to do is the following:
```shell script
export TEST_APP_USERNAME="my.username"
export TEST_APP_PASSWORD="Password1!"
export TEST_APP_ATTRIBUTE="an attribute with a value"
```
```go
package main

import (
    "github.com/engi-fyi/go-credentials/credential"
    "github.com/engi-fyi/go-credentials/factory"
)

func main() {
    myFact, _ := factory.New("test_app", false)
    myFact.Initialize()
    myCred, _ := credential.LoadFromEnvironment(myFact)
    print(myCred.Username)                  // "my.username"
    print(myCred.Password)                  // "Password1!"
    print(myCred.GetAttribute("attribute")) // "an attribute with a value"
    myCred.Save()                           // credentials file created at
                                            // ~/.test_app/credentials
}
```