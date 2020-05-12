Welcome to `go-credentials`! This project was built out of a need for a library that handles and manages credential files and allows them to be manipulated without having to deal with file-based libraries. Instead of loading multiple libraries for different file types inside of your library, load `go-credentials` and we'll handle the rest. Instead of complicated logic, simply call `myCredential.Save()`!

## Getting Started
The basic API of `go-credentials` is based on two separate object types. `Factory` is responsible for setting application-level settings. `Credential` is for managing your users credentials.

To get started, is all you need to do is create the following file at __~/.gcea/credentials__.
```
[default]
username = test@engi.fyi
password = !my_test_pass==word
```
Then, create a simple Go project, and enter in your `go-credentials` config:
To get started, is all you need to do is create the following files.
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

Then, install the dependancies and run the project!
```
user@console:~/$ go get -d .
user@console:~/$ go run main.go
10:07AM TRC Deserializing credential and profile. output_type=ini
10:07AM TRC No alternate username found.
10:07AM TRC No alternate password found.
10:07AM TRC Loading credentials from default profile profile=default
10:07AM TRC Building credential object.

Username: test@engi.fyi
Password: !my_test_pass==word

10:07AM TRC Setting attribute. key=last_updated
10:07AM TRC Serializing credential and profile. output_type=ini
10:07AM TRC Serializing credential and profile to ini file.
10:07AM TRC Serializing credential to ini file.
10:07AM TRC Getting alternate username and password labels.
10:07AM TRC No alternate username found.
10:07AM TRC No alternate password found.
10:07AM TRC Adding username and password to credentials file.
10:07AM TRC Saving credential ini file.
10:07AM TRC Credential ini file saved successfully.
10:07AM TRC Serializing profile to ini file.
10:07AM TRC Processing attributes.
10:07AM TRC Adding section. attribute=metadata
10:07AM TRC Adding attribute. attribute=last_update
10:07AM TRC Adding attribute. attribute=last_updated
10:07AM TRC Saving profile ini file.
10:07AM TRC Profile ini file saved successfully.
10:07AM TRC Deserializing credential and profile. output_type=ini
10:07AM TRC No alternate username found.
10:07AM TRC No alternate password found.
10:07AM TRC Loading credentials from default profile profile=default
10:07AM TRC Building credential object.

Username: test@engi.fyi
Password: !my_test_pass==word

10:07AM TRC Retrieving attribute. key=last_updated

Last Updated: 12/05/2020 10:07:49

```
## Basic API Overview
The following section gives a very basic overview of the `go-credentials` API. To view the full details, please visit our [Go Doc](https://pkg.go.dev/mod/github.com/engi-fyi/go-credentials) page.
### Factory

- `New`: creates a new Factory object.
- `SetOutputType`: sets the output type. Currently available types are `ini` with plans to implement `json`.

### Credential

 - `New`: Creates a new Credential object.
 - `Save`: Saves the credential object to file (based on Factory.OutputType)
 - `Load`: Loads a credential object from environment variables or files. See also `LoadFromEnvironment` and `LoadFromIniFile` if you would like to selectively load credentials from certain sources.
 - `SetAttribute`: Allows you to set an attribute on the Credential object.
 - `GetAttribute`: Gets an attribute that has been set on the credential object.
 
## Licensing

`go-credentials` is Licensed under the MIT License (the "License");

You may not use this file except in compliance with the License.
You may obtain a copy of the License at https://engi.fyi/mit-license/.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.

See the License for the specific language governing permissions and
limitations under the License.

## Copyright

Copyright &copy;2020 engi.fyi Contributors.

All contributors to engi.fyi can be found online at https://engi.fyi/contributors/.