package factory

import (
	"errors"
	"github.com/engi-fyi/go-credentials/global"
	"os"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

// New creates a very simple Factory object, with defaults based on the ApplicationName.
func New(applicationName string, useEnvironment bool) (*Factory, error) {
	keyRegex := regexp.MustCompile(global.REGEX_KEY_NAME)

	if applicationName == "" {
		return nil, errors.New(ERR_APPLICATION_NAME_BLANK)
	}

	if !keyRegex.MatchString(applicationName) {
		return nil, errors.New(ERR_KEY_MUST_MATCH_REGEX)
	}

	newFactory := Factory{
		ApplicationName: applicationName,
		UseEnvironment:  useEnvironment,
		OutputType:      global.OUTPUT_TYPE_INI,
	}

	initErr := newFactory.Initialize()

	if initErr != nil {
		return nil, initErr
	}

	return &newFactory, nil
}

// Initialize sets computed properties a Factory object. Specifically, it sets the value of ParentDirectory, ConfigDirectory and
// CredentialFile. If ParentDirectory does not exist, it will also create it. Alternates is also initialized as
// an empty map and the Initialized flag is set to true.
func (factory *Factory) Initialize() error {
	log.Trace().Msg("Initializing application credentials.")
	log.Trace().Msg("Retrieving user home directory.")
	homeDirectory, hdErr := os.UserHomeDir()

	if hdErr != nil {
		return hdErr
	}

	log.Trace().Str("Home Directory", homeDirectory).Msg("Found users home directory.")
	factory.ParentDirectory = homeDirectory + "/." + strings.ToLower(factory.ApplicationName) + "/"
	factory.ConfigDirectory = factory.ParentDirectory + "config/"
	factory.CredentialFile = factory.ParentDirectory + "credentials"

	if _, pdsErr := os.Stat(factory.ParentDirectory); os.IsNotExist(pdsErr) {
		log.Trace().Str("parent", factory.ParentDirectory).Msg("Creating parent directory.")
		mkErr := os.Mkdir(factory.ParentDirectory, os.ModeDir)

		if mkErr != nil {
			return mkErr
		}

		//#nosec
		modErr := os.Chmod(factory.ParentDirectory, 0700)

		if modErr != nil {
			return modErr
		}
	} else {
		log.Trace().Msg("Configuration directory exists, skipping.")
	}

	if _, cdsErr := os.Stat(factory.ConfigDirectory); os.IsNotExist(cdsErr) {
		log.Trace().Str("config", factory.ConfigDirectory).Msg("Creating config directory.")
		mkErr := os.Mkdir(factory.ConfigDirectory, os.ModeDir)

		if mkErr != nil {
			return mkErr
		}

		//#nosec
		modErr := os.Chmod(factory.ConfigDirectory, 0700)

		if modErr != nil {
			return modErr
		}
	} else {
		log.Trace().Msg("Config directory exists, skipping.")
	}

	factory.alternates = make(map[string]string)
	factory.Initialized = true
	log.Trace().Msg("Credential initialization complete.")
	return nil
}

// SetAlternateUsername sets a label to be used in lieu of username in environment variables.
func (factory *Factory) SetAlternateUsername(alternateUsername string) error {
	if alternateUsername != "" {
		factory.alternates["username"] = strings.ToLower(alternateUsername)
		return nil
	} else {
		return errors.New(ERR_ALTERNATE_USERNAME_CANNOT_BE_BLANK)
	}
}

// GetAlternateUsername gets a label to be used in lieu of username in environment variables.
func (factory *Factory) GetAlternateUsername() string {
	if val, exists := factory.alternates["username"]; exists {
		return val
	} else {
		return "username"
	}
}

// SetAlternatePassword sets a label to be used in lieu of password in environment variables.
func (factory *Factory) SetAlternatePassword(alternatePassword string) error {
	if alternatePassword != "" {
		factory.alternates["password"] = strings.ToLower(alternatePassword)
		return nil
	} else {
		return errors.New(ERR_ALTERNATE_PASSWORD_CANNOT_BE_BLANK)
	}
}

// GetAlternates returns both the alternates map
func (factory *Factory) GetAlternates() map[string]string {
	return factory.alternates
}

// GetAlternatePassword sets a label to be used in lieu of password in environment variables.
func (factory *Factory) GetAlternatePassword() string {
	if val, exists := factory.alternates["password"]; exists {
		return val
	} else {
		return "password"
	}
}

// SetOutputType determines which of the supported file types Credentials should be serialized to file as. The currently
// supported file types are ini with plans to implement json.
//
// TODO(7): Implement json format.
func (factory *Factory) SetOutputType(outputType string) error {
	if outputType == global.OUTPUT_TYPE_INI || outputType == global.OUTPUT_TYPE_JSON {
		factory.OutputType = outputType
		return nil
	} else {
		return errors.New(ERR_INVALID_OUTPUT_TYPE)
	}
}

// Set environment keys is a function that sets the alternates for both username and password at the same time. This
// is the same as calling SetAlternateUsername, then calling SetAlternatePassword in a seperate call.
func (factory *Factory) SetEnvironmentKeys(usernameKey string, passwordKey string) error {
	keyRegex := regexp.MustCompile(global.REGEX_KEY_NAME)

	if keyRegex.MatchString(usernameKey) {
		log.Trace().Str("key", usernameKey).Msg("Alternate environment key for username registered.")
		factory.alternates["username"] = usernameKey
	} else {
		log.Error().Msg("Sorry the key for username does not match the requirements.")
		return errors.New(ERR_KEY_MUST_MATCH_REGEX)
	}

	if keyRegex.MatchString(passwordKey) {
		log.Trace().Str("key", passwordKey).Msg("Alternate environment key for password registered.")
		factory.alternates["password"] = passwordKey
	} else {
		log.Error().Msg("Sorry the key for password does not match the requirements.")
		return errors.New(ERR_KEY_MUST_MATCH_REGEX)
	}

	return nil
}
