package serializer

import (
	"errors"
	"github.com/engi-fyi/go-credentials/global"
	"os"
	"strings"
)

/*
ToEnv is responsible for serializing a Credential/Profile combination into the environment. Both of the credential
values are serialized as so:

	APPLICATION_NAME::PROFILE_NAME::USERNAME (or alternate)
	APPLICATION_NAME::PROFILE_NAME::PASSWORD (or alternate)

Profile attributes are also added to the profile in the following format:

	APPLICATION_NAME::PROFILE_NAME::ATTRIBUTE::SECTION_NAME::KEY_VALUE

SECTION_NAME can be blank.
*/
func (thisSerializer *Serializer) ToEnv(username string, password string, attributes map[string]map[string]string) error {
	thisSerializer.Factory.Log.Info().Msg("Serializing credential and profile to environment.")
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
	thisSerializer.Factory.Log.Trace().Str("key", usernameKey).Msg("Setting username environment variable.")
	setErr := os.Setenv(usernameKey, username)

	if setErr != nil {
		thisSerializer.Factory.Log.Error().Err(setErr).Msg("Error setting username in environment.")
		return setErr
	}

	thisSerializer.Factory.Log.Info().Msg("Username set.")

	passwordKey := strings.ToUpper(prefix + thisSerializer.Factory.GetAlternatePassword())
	thisSerializer.Factory.Log.Trace().Str("key", passwordKey).Msg("Setting password environment variable.")
	setErr = os.Setenv(passwordKey, password)

	if setErr != nil {
		thisSerializer.Factory.Log.Error().Err(setErr).Msg("Error setting password in environment.")
		return setErr
	}

	thisSerializer.Factory.Log.Info().Msg("Password set.")

	return nil
}

func (thisSerializer *Serializer) saveProfileEnv(attributes map[string]map[string]string) error {
	prefix := thisSerializer.getEnvPrefix()

	for key, value := range attributes {
		for subKey, subValue := range value {
			fullKey := prefix + "ATTRIBUTE::" + strings.ToUpper(key) + "::" + strings.ToUpper(subKey)
			thisSerializer.Factory.Log.Trace().Str("key", fullKey).Msg("Setting attribute environment variable.")
			setErr := os.Setenv(fullKey, subValue)

			if setErr != nil {
				return setErr
			}
		}
	}

	thisSerializer.Factory.Log.Trace().Msg("Attribute environment variables set successfully.")

	return nil
}

/*
FromEnv is responsible for scanning environment variables and retrieves applicable variables that have
the prefix of applicationName and that have at least two "::" in them (which is the separator). The format for an
environment variable managed by serializer is:

	APPLICATION_NAME::PROFILE_NAME::FIELD_TYPE::SECTION_NAME::KEY_VALUE

Important to note is that if SECTION_NAME is blank then the default profile will be filled, and FIELD_TYPE can be one
of three values which are USERNAME, PASSWORD, or ATTRIBUTE. If the FIELD_TYPE is ATTRIBUTE, then KEY_VALUE is mandatory.
*/
func (thisSerializer *Serializer) FromEnv() (string, string, map[string]map[string]string, error) {
	thisSerializer.Factory.Log.Info().Msg("Deserializing credential and profile from environment.")
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

	thisSerializer.Factory.Log.Debug().Str("username_label", thisSerializer.Factory.GetAlternateUsername()).Msg("Found username label.")

	if _, exists := parsedVariables[thisSerializer.Factory.GetAlternatePassword()]; !exists {
		return "", "", errors.New(ERR_REQUIRED_VARIABLE_PASSWORD_NOT_FOUND)
	}

	thisSerializer.Factory.Log.Debug().Str("password_label", thisSerializer.Factory.GetAlternatePassword()).Msg("Found password label.")

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

/*
ParseEnvironmentVariable is able to process an environment variable and see if matches the expected format for our
environment variables.

The format for an environment variable managed by serializer is:

	APPLICATION_NAME::PROFILE_NAME::FIELD_TYPE::SECTION_NAME::KEY_VALUE

Important to note is that if SECTION_NAME is blank then the default profile will be filled, and FIELD_TYPE can be one
of three values which are USERNAME, PASSWORD, or ATTRIBUTE. If the FIELD_TYPE is ATTRIBUTE, then KEY_VALUE is mandatory.
*/
func (thisSerializer *Serializer) ParseEnvironmentVariable(environmentVariable string) (string, string, string, bool) {
	if strings.Count(environmentVariable, "::") < 2 {
		return "", "", "", false
	}

	if len(environmentVariable) < len(thisSerializer.Factory.ApplicationName+"::") {
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
