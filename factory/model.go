package factory

type Factory struct {
	ApplicationName        string            `json:"application_name"`
	ConfigurationDirectory string            `json:"configuration_directory"`
	CredentialFile         string            `json:"credential_file"`
	UseEnvironment         bool              `json:"use_environment"`
	Initialized            bool              `json:"initialized"`
	OutputType             string            `json:"output_type"`
	Alternates             map[string]string `json:"alternates"`
}
