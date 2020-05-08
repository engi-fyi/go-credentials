Welcome to `go-credentials`! This project was built out of a need for a library that handles and manages credential files and allows them to be manipulated without having to deal with file-based libraries. Instead of loading multiple libraries for different file types inside of your library, load `go-credentials` and we'll handle the rest. Instead of complicated logic, simply call `myCredential.Save()`!

## Getting Started
The basic API of `go-credentials` is based on two separate object types. `Factory` is responsible for setting application-level settings. `Credential` is for managing your users credentials.

To get started, is all you need to do is the following in your shell:
```
export TEST_APP_USERNAME="my.username"
export TEST_APP_PASSWORD="Password1!"
export TEST_APP_ATTRIBUTE="an attribute with a value"
```
Then, create a simple Go project, and enter in your `go-credentials` config:
```
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
## Basic API Overview
The following section gives a very basic overview of the `go-credentials` API. To view the full details, please visit out Go Doc page.
### Factory

- `New`: creates a new Factory object.
- `Initialize`: sets all of the computed properties of that Factory object.
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