package credential

import (
	"github.com/engi-fyi/go-credentials/factory"
	"github.com/engi-fyi/go-credentials/profile"
)

// Credential is the main object used by the go-credential library to manage user credentials. Username and Password are
// exported by default. The attributes of the Credential are not exported, as they should only be accessed via
// SetAttribute or GetAttribute.
type Credential struct {
	Username             string
	Password             string
	Initialized          bool
	Factory              *factory.Factory
	Profile              *profile.Profile
	environmentVariables []string
	selectedSection      string
}
