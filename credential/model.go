package credential

import "github.com/HammoTime/go-credentials/factory"

type Credential struct {
	Username             string            `json:"username"`
	Password             string            `json:"password"`
	Initialized          bool              `json:"initialized"`
	attributes           map[string]string `json:"attributes"`
	Factory              *factory.Factory
	environmentVariables []string `json:"linked_environment_variables"`
}
