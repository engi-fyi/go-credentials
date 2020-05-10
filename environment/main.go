package environment

import (
	"errors"
	"github.com/engi-fyi/go-credentials/global"
	"github.com/rs/zerolog/log"
	"os"
	"regexp"
	"strings"
)

// Deploy is responsible for setting environment variables based on the values passed in. If no alternates are set, the
// username and password will be set as APP_NAME_USERNAME and APP_NAME_PASSWORD respectively. If alternates are
// set, the key values will be replaced with those names instead. All attributes will then be set in the format of
// APP_NAME_ATTRIBUTE_NAME. The names of each environment variable are returned as a string array. This is to allow
// for cleanup of the environment variables if you want them deleted.
func Deploy(applicationName string, username string, password string, alternates map[string]string, attributes map[string]map[string]string) ([]string, error) {
	var setEnvironmentVariables []string

	if username != "" && applicationName != "" {
		log.Trace().Str("Application Name", applicationName)

		prefix := strings.ToUpper(applicationName)
		log.Trace().Str("prefix", prefix).Msg("Prefix set.")

		usernameKey, usernameErr := deployUsername(prefix, username, alternates)

		if usernameErr != nil {
			return setEnvironmentVariables, usernameErr
		} else {
			setEnvironmentVariables = append(setEnvironmentVariables, usernameKey)
		}

		passwordKey, passwordErr := deployPassword(prefix, password, alternates)

		if passwordErr != nil {
			return setEnvironmentVariables, passwordErr
		} else {
			setEnvironmentVariables = append(setEnvironmentVariables, passwordKey)
		}

		attributeKey, attributeErr := deployAttributes(prefix, attributes)

		if attributeErr != nil {
			return setEnvironmentVariables, attributeErr
		} else {
			setEnvironmentVariables = append(setEnvironmentVariables, attributeKey...)
		}
	} else {
		log.Error().Msg("Sorry we need valid values.")
		return nil, errors.New(ERR_CANT_LOAD_WITH_EMPTY_VALUES)
	}

	return setEnvironmentVariables, nil
}

// CleanEnvironment deletes each environment variable in the variables array.
func CleanEnvironment(variables []string) error {
	for key := range variables {
		if _, ok := os.LookupEnv(variables[key]); ok {
			unsetErr := os.Unsetenv(variables[key])

			if unsetErr != nil {
				return unsetErr
			}
		}
	}

	return nil
}

// Load is responsible for scanning environment variables and retrieves applicable variables that have
// the prefix of applicationName. When it finds an environment variable with this prefix, it will note the key (sans
// prefix) and the value. The resulting map that is returned is combination of the cleansed key and the value that was
// set. This map can be empty.
func Load(applicationName string, alternates map[string]string) (map[string]string, error) {
	envVariables := os.Environ()
	parsedVariables := make(map[string]string)

	for i := range envVariables {
		appPrefix := strings.ToUpper(applicationName) + "_"
		splitIndex := strings.Index(envVariables[i], "=")
		key := strings.ToUpper(envVariables[i][0:splitIndex])
		value := envVariables[i][splitIndex+1 : len(envVariables[i])]

		if len(key) > len(appPrefix) {
			if key[0:len(appPrefix)] == appPrefix {
				parsedKey := strings.ToLower(strings.Replace(key, appPrefix, "", 1))
				log.Trace().Str("Key", parsedKey).Msg("Variable found in environment.")
				parsedVariables[parsedKey] = value
			}
		}
	}

	if _, exists := parsedVariables["username"]; !exists {
		if _, exists := parsedVariables[strings.ToLower(alternates["username"])]; !exists {
			return make(map[string]string), errors.New(ERR_REQUIRED_VARIABLE_USERNAME_NOT_FOUND)
		}

		log.Trace().Str("alternate_username", alternates["username"]).Msg("Found alternate username.")
	}

	if _, exists := parsedVariables["password"]; !exists {
		if _, exists := parsedVariables[strings.ToLower(alternates["password"])]; !exists {
			return make(map[string]string), errors.New(ERR_REQUIRED_VARIABLE_PASSWORD_NOT_FOUND)
		}

		log.Trace().Str("alternate_password", alternates["password"]).Msg("Found alternate password.")
	}

	return parsedVariables, nil
}

func deployUsername(prefix string, username string, alternates map[string]string) (string, error) {
	usernameKey := prefix + "_USERNAME"
	keyRegex := regexp.MustCompile(global.REGEX_KEY_NAME)

	if replaceKey, ok := alternates["username"]; ok {
		usernameKey = prefix + "_" + strings.ToUpper(replaceKey)

		if !keyRegex.MatchString(usernameKey) {
			return "", errors.New(ERR_KEY_MUST_MATCH_REGEX)
		}
	}

	log.Trace().Str("key", usernameKey).Msg("Setting username environment variable.")
	setErr := os.Setenv(usernameKey, username)

	if setErr != nil {
		return usernameKey, setErr
	}

	return usernameKey, nil
}

func deployPassword(prefix string, password string, alternates map[string]string) (string, error) {
	if password != "" {
		passwordKey := prefix + "_PASSWORD"
		keyRegex := regexp.MustCompile(global.REGEX_KEY_NAME)

		if replaceKey, ok := alternates["password"]; ok {
			passwordKey = prefix + "_" + strings.ToUpper(replaceKey)

			if !keyRegex.MatchString(passwordKey) {
				return "", errors.New(ERR_KEY_MUST_MATCH_REGEX)
			}
		}

		log.Trace().Str("key", passwordKey).Msg("Setting password environment variable.")
		setErr := os.Setenv(passwordKey, password)

		if setErr != nil {
			return passwordKey, setErr
		}

		return passwordKey, nil
	}

	return "", nil
}

func deployAttributes(prefix string, attributes map[string]map[string]string) ([]string, error) {
	keyRegex := regexp.MustCompile(global.REGEX_KEY_NAME)
	var attributeKeys []string

	for key, value := range attributes {
		subPrefix := prefix

		if strings.ToUpper(key) != strings.ToUpper(global.DEFAULT_PROFILE_NAME) {
			subPrefix = prefix + "_" + strings.ToUpper(key)
		}

		for subKey, subValue := range value {
			fullKey := subPrefix + "_" + strings.ToUpper(subKey)
			if !keyRegex.MatchString(fullKey) {
				return attributeKeys, errors.New(ERR_KEY_MUST_MATCH_REGEX)
			}

			log.Trace().Str("key", fullKey).Msg("Setting attribute environment variable.")
			setErr := os.Setenv(fullKey, subValue)

			if setErr != nil {
				return attributeKeys, setErr
			}

			attributeKeys = append(attributeKeys, fullKey)
		}
	}

	return attributeKeys, nil
}
