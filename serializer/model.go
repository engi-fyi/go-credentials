package serializer

import "github.com/engi-fyi/go-credentials/factory"

/*
Serializer represents the basic settings required to (de)serialize a Credential and Profile.
*/
type Serializer struct {
	Factory        *factory.Factory
	ProfileName    string
	CredentialFile string
	ConfigFile     string
	Initialized    bool
}

type credentialSerializer struct {
	Credentials map[string]serializedCredentials `json:"credentials",yaml:"credentials"`
}

type serializedCredentials struct {
	Username string `json:"username",yaml:"username"`
	Password string `json:"password",yaml:"password"`
}

type profileSerializer struct {
	Attributes map[string]map[string]string `json:"attributes",yaml:"attributes"`
}
