package factory

import "github.com/rs/zerolog"

// Factory is the object that is used to store all of the application-level global configuration. All the settings
// for saving, setting, searching and finding credentials are in this object.
//
// Application Name: this is the name of the application.
// ParentDirectory: automatically set to ~/.application_name
// CredentialFile: automatically set to configurationDirectory + "/credentials".
// UseEnvironment: can I load variables into the environment. Only set this if you intend on use LoadEnv.
// Initialized: has all of my configuration been initialized correctly?
// Output Type: the file type that the CredentialFile contents should be.
// Alternates: if username or password are set, those names are set
type Factory struct {
	ApplicationName string
	ParentDirectory string
	CredentialFile  string
	ConfigDirectory string
	UseEnvironment  bool
	Initialized     bool
	OutputType      string
	Log             *zerolog.Logger
	alternates      map[string]string
}
