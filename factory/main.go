package factory

import (
	"errors"
	"github.com/engi-fyi/go-credentials/global"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"regexp"
	"strings"
)

/*
New creates a very simple Factory object, with defaults based on the ApplicationName.
*/
func New(applicationName string) (*Factory, error) {
	keyRegex := regexp.MustCompile(global.REGEX_KEY_NAME)

	if applicationName == "" {
		return nil, errors.New(ERR_APPLICATION_NAME_BLANK)
	}

	if !keyRegex.MatchString(applicationName) {
		return nil, errors.New(ERR_KEY_MUST_MATCH_REGEX)
	}

	newFactory := Factory{
		ApplicationName: applicationName,
		OutputType:      global.OUTPUT_TYPE_INI,
	}

	initErr := newFactory.Initialize()

	if initErr != nil {
		return nil, initErr
	}

	return &newFactory, nil
}

func (thisFactory *Factory) initLogger() {
	logLevel := "disabled"

	if value, ok := os.LookupEnv(global.LOG_LEVEL_ENVIRONMENT_KEY); ok {
		logLevel = value
	}

	if outputType, _ := os.LookupEnv(global.LOG_OUTPUT_TYPE_ENV_KEY); outputType == "pretty" {
		thisFactory.ModifyLogger(logLevel, true)
	} else {
		thisFactory.ModifyLogger(logLevel, false)
	}
}

/*
ModifyLogger is responsible for reconfiguring the log level and whether or not pretty output is used. Available log
levels are panic, fatal, error, warn, info, debug, trace, and disabled. By default, all logging output is disabled.
*/
func (thisFactory *Factory) ModifyLogger(logLevel string, pretty bool) {
	var logger zerolog.Logger

	if pretty {
		logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(getLogLevel(logLevel))
	} else {
		logger = zerolog.New(os.Stderr).With().Timestamp().Logger().Level(getLogLevel(logLevel))
	}

	thisFactory.Log = &logger
}

func getLogLevel(logLevel string) zerolog.Level {
	switch logLevel {
	case "panic":
		return zerolog.PanicLevel
	case "fatal":
		return zerolog.FatalLevel
	case "error":
		return zerolog.ErrorLevel
	case "warn":
		return zerolog.WarnLevel
	case "info":
		return zerolog.InfoLevel
	case "debug":
		return zerolog.DebugLevel
	case "trace":
		return zerolog.TraceLevel
	case "disabled":
		return zerolog.Disabled
	default:
		return zerolog.InfoLevel
	}
}

/*
Initialize sets computed properties a Factory object. Specifically, it sets the value of ParentDirectory, ConfigDirectory and
CredentialFile. If ParentDirectory does not exist, it will also create it. Alternates is also initialized as
an empty map and the Initialized flag is set to true. The logger for the Factory is also initialized here.
*/
func (thisFactory *Factory) Initialize() error {
	thisFactory.initLogger()
	thisFactory.Log.Trace().Msg("Initializing application credentials.")
	thisFactory.Log.Trace().Msg("Retrieving user home directory.")
	homeDirectory, hdErr := os.UserHomeDir()

	if hdErr != nil {
		return hdErr
	}

	thisFactory.Log.Trace().Str("Home Directory", homeDirectory).Msg("Found users home directory.")
	thisFactory.ParentDirectory = homeDirectory + "/." + strings.ToLower(thisFactory.ApplicationName) + "/"
	thisFactory.ConfigDirectory = thisFactory.ParentDirectory + "config/"
	thisFactory.CredentialFile = thisFactory.ParentDirectory + "credentials"

	if _, pdsErr := os.Stat(thisFactory.ParentDirectory); os.IsNotExist(pdsErr) {
		thisFactory.Log.Trace().Str("parent", thisFactory.ParentDirectory).Msg("Creating parent directory.")
		mkErr := os.Mkdir(thisFactory.ParentDirectory, os.ModeDir)

		if mkErr != nil {
			return mkErr
		}

		//#nosec
		modErr := os.Chmod(thisFactory.ParentDirectory, 0700)

		if modErr != nil {
			return modErr
		}
	} else {
		thisFactory.Log.Trace().Msg("Configuration directory exists, skipping.")
	}

	if _, cdsErr := os.Stat(thisFactory.ConfigDirectory); os.IsNotExist(cdsErr) {
		thisFactory.Log.Trace().Str("config", thisFactory.ConfigDirectory).Msg("Creating config directory.")
		mkErr := os.Mkdir(thisFactory.ConfigDirectory, os.ModeDir)

		if mkErr != nil {
			return mkErr
		}

		//#nosec
		modErr := os.Chmod(thisFactory.ConfigDirectory, 0700)

		if modErr != nil {
			return modErr
		}
	} else {
		thisFactory.Log.Trace().Msg("Config directory exists, skipping.")
	}

	thisFactory.alternates = make(map[string]string)
	thisFactory.Initialized = true
	thisFactory.Log.Trace().Msg("Credential initialization complete.")
	return nil
}

/*
SetAlternateUsername sets a label to be used in lieu of username in environment variables.
*/
func (thisFactory *Factory) SetAlternateUsername(alternateUsername string) error {
	if alternateUsername != "" {
		thisFactory.alternates["username"] = strings.ToLower(alternateUsername)
		thisFactory.Log.Trace().Str("username", thisFactory.alternates["username"]).Msg("Alternate set.")
		return nil
	} else {
		thisFactory.Log.Error().Msg(ERR_ALTERNATE_USERNAME_CANNOT_BE_BLANK)
		return errors.New(ERR_ALTERNATE_USERNAME_CANNOT_BE_BLANK)
	}
}

/*
GetAlternateUsername gets a label to be used in lieu of username in environment variables.
*/
func (thisFactory *Factory) GetAlternateUsername() string {
	if val, exists := thisFactory.alternates["username"]; exists {
		thisFactory.Log.Trace().Str("username", thisFactory.alternates["username"]).Msg("Found alternate.")
		return val
	} else {
		thisFactory.Log.Trace().Msg("No alternate username found.")
		return "username"
	}
}

/*
SetAlternatePassword sets a label to be used in lieu of password in environment variables.
*/
func (thisFactory *Factory) SetAlternatePassword(alternatePassword string) error {
	if alternatePassword != "" {
		thisFactory.alternates["password"] = strings.ToLower(alternatePassword)
		thisFactory.Log.Trace().Str("password", thisFactory.alternates["password"]).Msg("Alternate set.")
		return nil
	} else {
		thisFactory.Log.Error().Msg(ERR_ALTERNATE_PASSWORD_CANNOT_BE_BLANK)
		return errors.New(ERR_ALTERNATE_PASSWORD_CANNOT_BE_BLANK)
	}
}

/*
GetAlternates returns the username and password.
*/
func (thisFactory *Factory) GetAlternates() (string, string) {
	return thisFactory.GetAlternateUsername(), thisFactory.GetAlternatePassword()
}

/*
GetAlternatePassword sets a label to be used in lieu of password in environment variables.
*/
func (thisFactory *Factory) GetAlternatePassword() string {
	if val, exists := thisFactory.alternates["password"]; exists {
		thisFactory.Log.Trace().Str("password", thisFactory.alternates["password"]).Msg("Found alternate.")
		return val
	} else {
		thisFactory.Log.Trace().Msg("No alternate password found.")
		return "password"
	}
}

/*
SetOutputType determines which of the supported file types Credentials should be serialized to file as. The currently
supported file types are ini with plans to implement json.
*/
func (thisFactory *Factory) SetOutputType(outputType string) error {
	if outputType == global.OUTPUT_TYPE_INI ||
		outputType == global.OUTPUT_TYPE_JSON ||
		outputType == global.OUTPUT_TYPE_ENV {
		thisFactory.Log.Trace().Str("output_type", outputType).Msg("Output type set.")
		thisFactory.OutputType = outputType
		return nil
	} else {
		thisFactory.Log.Error().Str("output_type", outputType).Msg("Invalid output type.")
		return errors.New(ERR_INVALID_OUTPUT_TYPE)
	}
}

/*
SetAlternates is a function that sets the alternates for both username and password at the same time. This
is the same as calling SetAlternateUsername, then calling SetAlternatePassword in a separate call. If you enter a blank
value in either parameter, it won't be set and an error will be returned. However, the other variable will be set.
*/
func (thisFactory *Factory) SetAlternates(usernameKey string, passwordKey string) error {
	usernameErr := thisFactory.SetAlternateUsername(usernameKey)
	passwordErr := thisFactory.SetAlternatePassword(passwordKey)

	if usernameErr != nil {
		return usernameErr
	}

	if passwordErr != nil {
		return passwordErr
	}

	return nil
}
