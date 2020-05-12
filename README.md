# go-credentials
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/engi-fyi/go-credentials)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/engi-fyi/go-credentials)
[![GoDoc](https://img.shields.io/badge/go--doc-go--credentials-blue)](https://pkg.go.dev/mod/github.com/engi-fyi/go-credentials)
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
    - Set the output type of the credentials (environment, ini, and json supported).
    - Responsible for logging.
2. `Credential`: represents a user's credentials.
    - Username/Password defined on model.
    - Can have a profile.
    - Save and Load Credentials (and Profiles).
3. `Profile`: represents a profile, containing variables specific to a profile.
    - Username/Password defined on model.
    - Set Attributes (including sections).
    - Get Attributes (including sections).
    
To get started, is all you need to do is create the following files.

__~/.gcea/credentials__
```
[default]
username = test@engi.fyi
password = !my_test_pass==word
```
__main.go__
```go
package main

import (
    "github.com/engi-fyi/go-credentials/credential"
    "github.com/engi-fyi/go-credentials/factory"
    "fmt" 
    "time"
)

func main() {
	myFact, _ := factory.New("gcea", false) // go-credentials-example-application
	myFact.ModifyLogger("trace", true) // let's see under the hood and make it pretty.
	myCredential, _ := credential.Load(myFact)

	fmt.Printf("Username: %s\n", myCredential.Username)
	fmt.Printf("Password: %s\n\n", myCredential.Password)

	myCredential.SetSectionAttribute("metadata", "last_updated", time.Now().Format("02/01/2006 15:04:05"))
	myCredential.Save()
	myCredential = nil

	yourCredential, _ := credential.Load(myFact)

	fmt.Printf("Username: %s\n", yourCredential.Username)
	fmt.Printf("Password: %s\n", yourCredential.Password)

	lastUpdated, _ := yourCredential.GetSectionAttribute("metadata", "last_updated")
	fmt.Printf("Last Updated: %s\n", lastUpdated)
}
```