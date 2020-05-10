package serializer

import (
	"errors"
	"github.com/engi-fyi/go-credentials/global"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

// BUG(4): Respect alternates when saving to file.
// TODO(1000): Should we actually load an existing credential or would the user expect the whole state to be serialized?
func (thisSerializer *Serializer) ToEnv(username string, password string, attributes map[string]map[string]string) error {
	log.Trace().Msg("Serializing credential and profile to environment.")
	credentialErr := thisSerializer.saveCredentialEnv(username, password)

	if credentialErr != nil {
		return credentialErr
	}

	profileErr := thisSerializer.saveProfileEnv(attributes)

	if profileErr != nil {
		return profileErr
	}

	return nil
}

func (thisSerializer *Serializer) saveCredentialEnv(username string, password string) error {
	prefix := thisSerializer.getEnvPrefix()

	usernameKey := strings.ToUpper(prefix + thisSerializer.Factory.GetAlternateUsername())
	log.Trace().Str("key", usernameKey).Msg("Setting username environment variable.")
	setErr := os.Setenv(usernameKey, username)

	if setErr != nil {
		log.Error().Err(setErr).Msg("Error setting username in environment.")
		return setErr
	}

	log.Trace().Msg("Username set.")

	passwordKey := strings.ToUpper(prefix + thisSerializer.Factory.GetAlternatePassword())
	log.Trace().Str("key", passwordKey).Msg("Setting password environment variable.")
	setErr = os.Setenv(passwordKey, password)

	if setErr != nil {
		log.Error().Err(setErr).Msg("Error setting password in environment.")
		return setErr
	}

	log.Trace().Msg("Password set.")

	return nil
}

func (thisSerializer *Serializer) saveProfileEnv(attributes map[string]map[string]string) error {
	prefix := thisSerializer.getEnvPrefix()

	for key, value := range attributes {
		for subKey, subValue := range value {
			fullKey := prefix + "ATTRIBUTE::" + strings.ToUpper(key) + "::" + strings.ToUpper(subKey)
			log.Trace().Str("key", fullKey).Msg("Setting attribute environment variable.")
			setErr := os.Setenv(fullKey, subValue)

			if setErr != nil {
				return setErr
			}
		}
	}

	return nil
}

func (thisSerializer *Serializer) FromEnv() (string, string, map[string]map[string]string, error)  {
	parsedVariables, parseErr := thisSerializer.loadVariablesEnv()

	if parseErr != nil {
		return "", "", make(map[string]map[string]string), parseErr
	}

	username, password, credErr := thisSerializer.loadCredentialEnv(parsedVariables[global.NO_SECTION_KEY])

	if credErr != nil {
		return "", "", make(map[string]map[string]string), credErr
	}

	attributes, attrErr := thisSerializer.loadProfileEnv(parsedVariables)

	if attrErr != nil {
		return "", "", make(map[string]map[string]string), attrErr
	}

	return username, password, attributes, nil
}

func (thisSerializer *Serializer) loadCredentialEnv(parsedVariables map[string]string) (string, string, error) {
	if _, exists := parsedVariables[thisSerializer.Factory.GetAlternateUsername()]; !exists {
		return "", "", errors.New(ERR_REQUIRED_VARIABLE_USERNAME_NOT_FOUND)
	}

	log.Trace().Str("username_label", thisSerializer.Factory.GetAlternateUsername()).Msg("Found username label.")

	if _, exists := parsedVariables[thisSerializer.Factory.GetAlternatePassword()]; !exists {
		return "", "", errors.New(ERR_REQUIRED_VARIABLE_PASSWORD_NOT_FOUND)
	}

	log.Trace().Str("password_label", thisSerializer.Factory.GetAlternatePassword()).Msg("Found password label.")

	return parsedVariables[thisSerializer.Factory.GetAlternateUsername()],
	parsedVariables[thisSerializer.Factory.GetAlternatePassword()],
	nil
}

// currently only supports default profile, maybe change key name at some point in future
// like maybe APP_NAME_ATTRIBUTE::SECTION_NAME::VARIABLE_NAME
func (thisSerializer *Serializer) loadProfileEnv(parsedVariables map[string]map[string]string) (map[string]map[string]string, error) {
	delete(parsedVariables[global.NO_SECTION_KEY], thisSerializer.Factory.GetAlternateUsername())
	delete(parsedVariables[global.NO_SECTION_KEY], thisSerializer.Factory.GetAlternatePassword())

	return parsedVariables, nil
}

func (thisSerializer *Serializer) loadVariablesEnv() (map[string]map[string]string, error) {
	envVariables := os.Environ()
	parsedVariables := make(map[string]map[string]string)
	parsedVariables[global.NO_SECTION_KEY] = make(map[string]string)

	for i := range envVariables {
		splitIndex := strings.Index(envVariables[i], "=")
		key := strings.ToUpper(envVariables[i][0:splitIndex])
		value := envVariables[i][splitIndex+1 : len(envVariables[i])]
		profileName, fieldName, sectionName, didParse := thisSerializer.ParseEnvironmentVariable(key)

		if didParse {
			if profileName == thisSerializer.ProfileName {
				if sectionName == "" {
					parsedVariables[global.NO_SECTION_KEY][fieldName] = value
				} else {
					if _, ok := parsedVariables[sectionName]; !ok {
						parsedVariables[sectionName] = make(map[string]string)
					}

					parsedVariables[sectionName][fieldName] = value
				}
			}
		}
	}

	return parsedVariables, nil
}

func (thisSerializer *Serializer) getEnvPrefix() string {
	return strings.ToUpper(thisSerializer.Factory.ApplicationName) + "::" + strings.ToUpper(thisSerializer.ProfileName) + "::"
}

func (thisSerializer *Serializer) ParseEnvironmentVariable(environmentVariable string) (string, string, string, bool) {
	if strings.Count(environmentVariable, "::") < 2 {
		return "", "", "", false
	}

	if len(environmentVariable) < len(thisSerializer.Factory.ApplicationName + "::") {
		return "", "", "", false
	}

	environmentVariable = environmentVariable[strings.Index(environmentVariable, "::")+2:]

	profile := strings.ToLower(environmentVariable[0:strings.Index(environmentVariable, "::")])
	environmentVariable = environmentVariable[strings.Index(environmentVariable, "::")+2:]
	fieldName := strings.ToLower(environmentVariable)

	if strings.Count(environmentVariable, "::") > 1 {
		fieldName = strings.ToLower(environmentVariable[0:strings.Index(environmentVariable, "::")])
	}

	if fieldName != "attribute" {
		return profile, fieldName, "", true
	}

	environmentVariable = environmentVariable[strings.Index(environmentVariable, "::")+2:]
	sectionName := strings.ToLower(environmentVariable[0:strings.Index(environmentVariable, "::")])

	fieldName = strings.ToLower(environmentVariable[strings.Index(environmentVariable, "::")+2:])

	return profile, fieldName, sectionName, true
}