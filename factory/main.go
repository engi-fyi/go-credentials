package factory

import (
	"errors"
	"github.com/HammoTime/go-credentials/global"
	"os"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

// TODO: Run Initialize() at the end of this.
func New(applicationName string, useEnvironment bool) (*Factory, error) {
	keyRegex := regexp.MustCompile(global.REGEX_KEY_NAME)

	if applicationName == "" {
		return nil, errors.New(ERR_APPLICATION_NAME_BLANK)
	}

	if !keyRegex.MatchString(applicationName) {
		return nil, errors.New(ERR_KEY_MUST_MATCH_REGEX)
	}

	return &Factory{
		ApplicationName: applicationName,
		UseEnvironment:  useEnvironment,
		OutputType:      global.OUTPUT_TYPE_INI,
	}, nil
}

// TODO: Include in New
func (factory *Factory) Initialize() error {
	log.Trace().Msg("Initializing application credentials.")
	log.Trace().Msg("Retrieving user home directory.")
	homeDirectory, hdErr := os.UserHomeDir()

	if hdErr != nil {
		return hdErr
	}

	log.Trace().Str("Home Directory", homeDirectory).Msg("Found users home directory.")
	factory.ConfigurationDirectory = homeDirectory + "/." + strings.ToLower(factory.ApplicationName) + "/"
	factory.CredentialFile = factory.ConfigurationDirectory + "credentials"

	if _, cdsErr := os.Stat(factory.ConfigurationDirectory); os.IsNotExist(cdsErr) {
		log.Trace().Str("Configuration Directory", factory.ConfigurationDirectory).Msg("Creating configuration directory.")
		mkErr := os.Mkdir(factory.ConfigurationDirectory, os.ModeDir)

		if mkErr != nil {
			return mkErr
		}

		//#nosec
		modErr := os.Chmod(factory.ConfigurationDirectory, 0700)

		if modErr != nil {
			return modErr
		}
	} else {
		log.Trace().Msg("Configuration directory exists, skipping.")
	}

	factory.Alternates = make(map[string]string)
	factory.Initialized = true
	log.Trace().Msg("Credential initialization complete.")
	return nil
}

func (factory *Factory) SetAlternateUsername(alternateUsername string) error {
	if alternateUsername != "" {
		factory.Alternates["username"] = strings.ToLower(alternateUsername)
		return nil
	} else {
		return errors.New(ERR_ALTERNATE_USERNAME_CANNOT_BE_BLANK)
	}
}

func (factory *Factory) GetAlternateUsername() string {
	if val, exists := factory.Alternates["username"]; exists {
		return val
	} else {
		return "username"
	}
}

func (factory *Factory) SetAlternatePassword(alternatePassword string) error {
	if alternatePassword != "" {
		factory.Alternates["password"] = strings.ToLower(alternatePassword)
		return nil
	} else {
		return errors.New(ERR_ALTERNATE_PASSWORD_CANNOT_BE_BLANK)
	}
}

func (factory *Factory) GetAlternatePassword() string {
	if val, exists := factory.Alternates["password"]; exists {
		return val
	} else {
		return "password"
	}
}

func (factory *Factory) SetOutputType(outputType string) error {
	if outputType == global.OUTPUT_TYPE_INI || outputType == global.OUTPUT_TYPE_JSON {
		factory.OutputType = outputType
		return nil
	} else {
		return errors.New(ERR_INVALID_OUTPUT_TYPE)
	}
}

func (factory *Factory) SetEnvironmentKeys(usernameKey string, passwordKey string) error {
	keyRegex := regexp.MustCompile(global.REGEX_KEY_NAME)

	if keyRegex.MatchString(usernameKey) {
		log.Trace().Str("key", usernameKey).Msg("Alternate environment key for username registered.")
		factory.Alternates["username"] = usernameKey
	} else {
		log.Error().Msg("Sorry the key for username does not match the requirements.")
		return errors.New(ERR_KEY_MUST_MATCH_REGEX)
	}

	if keyRegex.MatchString(passwordKey) {
		log.Trace().Str("key", passwordKey).Msg("Alternate environment key for password registered.")
		factory.Alternates["password"] = passwordKey
	} else {
		log.Error().Msg("Sorry the key for password does not match the requirements.")
		return errors.New(ERR_KEY_MUST_MATCH_REGEX)
	}

	return nil
}
