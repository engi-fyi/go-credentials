package environment

import (
	"errors"
	"github.com/rs/zerolog/log"
	"os"
	"regexp"
	"strings"
)

func Deploy(applicationName string, username string, password string, alternates map[string]string, attributes map[string]string) ([]string, error) {
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

	log.Trace().Str("alternate_username", alternates["username"]).Msg("Found alternate username.")
	log.Trace().Str("alternate_password", alternates["password"]).Msg("Found alternate password.")

	if _, exists := parsedVariables["username"]; !exists {
		if _, exists := parsedVariables[strings.ToLower(alternates["username"])]; !exists {
			return make(map[string]string), errors.New(ERR_REQUIRED_VARIABLE_USERNAME_NOT_FOUND)
		}
	}

	if _, exists := parsedVariables["password"]; !exists {
		if _, exists := parsedVariables[strings.ToLower(alternates["password"])]; !exists {
			return make(map[string]string), errors.New(ERR_REQUIRED_VARIABLE_PASSWORD_NOT_FOUND)
		}
	}

	return parsedVariables, nil
}

func deployUsername(prefix string, username string, alternates map[string]string) (string, error) {
	usernameKey := prefix + "_USERNAME"
	keyRegex := regexp.MustCompile(REGEX_KEY_NAME)

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
		keyRegex := regexp.MustCompile(REGEX_KEY_NAME)

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

func deployAttributes(prefix string, attributes map[string]string) ([]string, error) {
	keyRegex := regexp.MustCompile(REGEX_KEY_NAME)
	var attributeKeys []string

	for key, value := range attributes {
		fullKey := prefix + "_" + strings.ToUpper(key)

		if !keyRegex.MatchString(fullKey) {
			return attributeKeys, errors.New(ERR_KEY_MUST_MATCH_REGEX)
		}

		log.Trace().Str("key", fullKey).Msg("Setting attribute environment variable.")
		setErr := os.Setenv(fullKey, value)

		if setErr != nil {
			return attributeKeys, setErr
		}

		attributeKeys = append(attributeKeys, fullKey)
	}

	return attributeKeys, nil
}
