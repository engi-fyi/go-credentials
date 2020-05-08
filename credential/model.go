package credential

import "github.com/engi-fyi/go-credentials/factory"

// Credential is the main object used by the go-credential library to manage user credentials. Username and Password are
// exported by default. The attributes of the Credential are not exported, as they should only be accessed via
// SetAttribute or GetAttribute.
type Credential struct {
	Username             string            `json:"username"`
	Password             string            `json:"password"`
	Initialized          bool              `json:"initialized"`
	attributes           map[string]string `json:"attributes"`
	Factory              *factory.Factory
	environmentVariables []string `json:"linked_environment_variables"`
}
