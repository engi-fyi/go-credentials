package serializer

import "github.com/engi-fyi/go-credentials/factory"

type Serializer struct {
	Factory        *factory.Factory
	ProfileName    string
	CredentialFile string
	ConfigFile     string
	Initialized	   bool
}

type CredentialSerializer struct {
	Credentials map[string]SerializedCredentials `json:"credentials",yaml:"credentials"`
}

type SerializedCredentials struct {
	Username string `json:"username",yaml:"username"`
	Password string `json:"password",yaml:"password"`
}
type ProfileSerializer struct {
	Attributes map[string][]SerializedAttribute `json:"attributes",yaml:"attributes"`
}

type SerializedAttribute struct {
	Key   string `json:"key",yaml:"key"`
	Value string `json:"value",yaml:"value"`
}
