package profile

import "github.com/engi-fyi/go-credentials/factory"

// Profile is used to hold information about the user settings or metadata attached to a Credential. Each profile stores
// its credential in the main credentials file, then the other information is held under config/profile_name
type Profile struct {
	Name                 string
	ConfigFileLocation	 string
	attributes           map[string]map[string]string
	Initialized          bool
	Factory				 *factory.Factory
}