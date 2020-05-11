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
	Attributes map[string][]serializedAttribute `json:"attributes",yaml:"attributes"`
}

type serializedAttribute struct {
	Key   string `json:"key",yaml:"key"`
	Value string `json:"value",yaml:"value"`
}
