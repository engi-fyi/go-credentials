package factory

// Factory is the object that is used to store all of the application-level global configuration. All the settings
// for saving, setting, searching and finding credentials are in this object.
//
// Application Name: this is the name of the application.
// ConfigurationDirectory: automatically set to ~/.application_name
// CredentialFile: automatically set to configurationDirectory + "/credentials".
// UseEnvironment: can I load variables into the environment. Only set this if you intend on use LoadEnv.
// Initialized: has all of my configuration been initialized correctly?
// Output Type: the file type that the CredentialFile contents should be.
// Alternates: if username or password are set, those names are set
//
// TODO(6): Refactor alternates so that the property is not exported.
type Factory struct {
	ApplicationName        string            `json:"application_name"`
	ConfigurationDirectory string            `json:"configuration_directory"`
	CredentialFile         string            `json:"credential_file"`
	UseEnvironment         bool              `json:"use_environment"`
	Initialized            bool              `json:"initialized"`
	OutputType             string            `json:"output_type"`
	Alternates             map[string]string `json:"alternates"`
}
