package environment

import (
	"github.com/engi-fyi/go-credentials/global"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"testing"
)

func TestDeploySimple(t *testing.T) {
	assert := global.InitTest(t)

	prefix := strings.ToUpper(global.TEST_VAR_APPLICATION_NAME) + "_"

	log.Info().Msg("Testing deploy with no attributes or alternates.")
	setKeys, deployErr := Deploy(
		global.TEST_VAR_APPLICATION_NAME,
		global.TEST_VAR_USERNAME,
		global.TEST_VAR_PASSWORD,
		nil,
		nil)
	assert.NoError(deployErr)
	assert.Equal(global.TEST_VAR_USERNAME, os.Getenv(prefix+"USERNAME"))
	assert.Equal(global.TEST_VAR_PASSWORD, os.Getenv(prefix+"PASSWORD"))

	CleanEnvironment(setKeys)
}

func TestDeploySimpleNoPassword(t *testing.T) {
	assert := global.InitTest(t)

	prefix := strings.ToUpper(global.TEST_VAR_APPLICATION_NAME) + "_"

	log.Info().Msg("Testing deploy with no attributes or alternates.")
	setKeys, deployErr := Deploy(
		global.TEST_VAR_APPLICATION_NAME,
		global.TEST_VAR_USERNAME,
		"",
		nil,
		nil)
	assert.NoError(deployErr)
	_, passwordExists := os.LookupEnv(prefix + "PASSWORD")
	assert.Equal(global.TEST_VAR_USERNAME, os.Getenv(prefix+"USERNAME"))
	assert.False(passwordExists)

	CleanEnvironment(setKeys)
}

func TestDeployAttributes(t *testing.T) {
	assert := global.InitTest(t)

	attributes := make(map[string]map[string]string)
	attributes[global.DEFAULT_PROFILE_NAME] = make(map[string]string)
	attributes[global.DEFAULT_PROFILE_NAME][global.TEST_VAR_ATTRIBUTE_NAME_LABEL] = global.TEST_VAR_ATTRIBUTE_VALUE

	log.Info().Msg("Testing deploy with attributes and no alternates.")
	setKeys, deployErr := Deploy(
		global.TEST_VAR_APPLICATION_NAME,
		global.TEST_VAR_USERNAME,
		global.TEST_VAR_PASSWORD,
		nil,
		attributes)
	assert.NoError(deployErr)
	assert.Equal(global.TEST_VAR_ATTRIBUTE_VALUE, os.Getenv(global.TEST_VAR_ENVIRONMENT_ATTRIBUTE_NAME_LABEL))

	CleanEnvironment(setKeys)
}

func TestDeployAlternates(t *testing.T) {
	assert := global.InitTest(t)

	alternates := make(map[string]string)
	alternates["username"] = global.TEST_VAR_USERNAME_ALTERNATE_LABEL
	alternates["password"] = global.TEST_VAR_PASSWORD_ALTERNATE_LABEL

	log.Info().Msg("Testing deploy with attributes and alternates.")

	setKeys, deployErr := Deploy(
		global.TEST_VAR_APPLICATION_NAME,
		global.TEST_VAR_USERNAME,
		global.TEST_VAR_PASSWORD,
		alternates,
		nil)
	assert.NoError(deployErr)
	assert.Equal(global.TEST_VAR_USERNAME, os.Getenv(global.TEST_VAR_ENVIRONMENT_USERNAME_ALTERNATE_LABEL))
	assert.Equal(global.TEST_VAR_PASSWORD, os.Getenv(global.TEST_VAR_ENVIRONMENT_PASSWORD_ALTERNATE_LABEL))

	CleanEnvironment(setKeys)
}

func TestDeployAttributesAndAlternates(t *testing.T) {
	assert := global.InitTest(t)

	alternates := make(map[string]string)
	alternates["username"] = global.TEST_VAR_USERNAME_ALTERNATE_LABEL
	alternates["password"] = global.TEST_VAR_PASSWORD_ALTERNATE_LABEL
	attributes := make(map[string]map[string]string)
	attributes[global.DEFAULT_PROFILE_NAME] = make(map[string]string)
	attributes[global.DEFAULT_PROFILE_NAME][global.TEST_VAR_ATTRIBUTE_NAME_LABEL] = global.TEST_VAR_ATTRIBUTE_VALUE

	log.Info().Msg("Testing deploy with attributes and alternates.")

	setKeys, deployErr := Deploy(
		global.TEST_VAR_APPLICATION_NAME,
		global.TEST_VAR_USERNAME,
		global.TEST_VAR_PASSWORD,
		alternates,
		attributes)
	assert.NoError(deployErr)
	assert.Equal(global.TEST_VAR_ATTRIBUTE_VALUE, os.Getenv(global.TEST_VAR_ENVIRONMENT_ATTRIBUTE_NAME_LABEL))
	assert.Equal(global.TEST_VAR_USERNAME, os.Getenv(global.TEST_VAR_ENVIRONMENT_USERNAME_ALTERNATE_LABEL))
	assert.Equal(global.TEST_VAR_PASSWORD, os.Getenv(global.TEST_VAR_ENVIRONMENT_PASSWORD_ALTERNATE_LABEL))

	CleanEnvironment(setKeys)
}

func TestLoad(t *testing.T) {
	assert := global.InitTest(t)

	usernameKey := global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL
	passwordKey := global.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL
	attributeKey := global.TEST_VAR_ENVIRONMENT_ATTRIBUTE_NAME_LABEL
	alternateUsernameKey := global.TEST_VAR_ENVIRONMENT_USERNAME_ALTERNATE_LABEL
	alternatePasswordKey := global.TEST_VAR_ENVIRONMENT_PASSWORD_ALTERNATE_LABEL
	alternates := make(map[string]string)
	alternates[global.TEST_VAR_USERNAME_LABEL] = global.TEST_VAR_USERNAME_ALTERNATE_LABEL
	alternates[global.TEST_VAR_PASSWORD_LABEL] = global.TEST_VAR_PASSWORD_ALTERNATE_LABEL

	log.Trace().Msg("Testing basic load from environment variables.")
	_, loadErr := Load(global.TEST_VAR_APPLICATION_NAME, nil)
	assert.EqualError(loadErr, ERR_REQUIRED_VARIABLE_USERNAME_NOT_FOUND)
	os.Setenv(usernameKey, global.TEST_VAR_USERNAME)
	_, loadErr = Load(global.TEST_VAR_APPLICATION_NAME, nil)
	assert.EqualError(loadErr, ERR_REQUIRED_VARIABLE_PASSWORD_NOT_FOUND)
	os.Setenv(passwordKey, global.TEST_VAR_PASSWORD)
	values, loadErr := Load(global.TEST_VAR_APPLICATION_NAME, nil)
	assert.NoError(loadErr)
	assert.Equal(values["username"], global.TEST_VAR_USERNAME)
	assert.Equal(values["password"], global.TEST_VAR_PASSWORD)
	os.Setenv(attributeKey, global.TEST_VAR_ATTRIBUTE_VALUE)
	values, loadErr = Load(global.TEST_VAR_APPLICATION_NAME, nil)
	assert.NoError(loadErr)
	assert.Equal(global.TEST_VAR_ATTRIBUTE_VALUE, values[strings.ToLower(global.TEST_VAR_ATTRIBUTE_NAME_LABEL)])
	CleanEnvironment([]string{usernameKey, passwordKey, attributeKey})

	log.Trace().Msg("Testing alternate load from environment variables.")
	// Test alternates
	os.Setenv(alternateUsernameKey, global.TEST_VAR_USERNAME)
	os.Setenv(alternatePasswordKey, global.TEST_VAR_PASSWORD)
	values, loadErr = Load(global.TEST_VAR_APPLICATION_NAME, alternates)
	assert.NoError(loadErr)
	assert.Equal(global.TEST_VAR_USERNAME, values[global.TEST_VAR_USERNAME_ALTERNATE_LABEL])
	assert.Equal(global.TEST_VAR_PASSWORD, values[global.TEST_VAR_PASSWORD_ALTERNATE_LABEL])
	CleanEnvironment([]string{alternateUsernameKey, alternatePasswordKey})
}

func TestCleanEnvironment(t *testing.T) {
	assert := global.InitTest(t)

	_, exists := os.LookupEnv(global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL)
	assert.False(exists)
	_, exists = os.LookupEnv(global.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL)
	assert.False(exists)
	_, exists = os.LookupEnv(global.TEST_VAR_ENVIRONMENT_USERNAME_ALTERNATE_LABEL)
	assert.False(exists)
	_, exists = os.LookupEnv(global.TEST_VAR_ENVIRONMENT_PASSWORD_ALTERNATE_LABEL)
	assert.False(exists)
	_, exists = os.LookupEnv(global.TEST_VAR_ENVIRONMENT_ATTRIBUTE_NAME_LABEL)
	assert.False(exists)

	os.Setenv(global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL, global.TEST_VAR_USERNAME)
	os.Setenv(global.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL, global.TEST_VAR_PASSWORD)
	os.Setenv(global.TEST_VAR_ENVIRONMENT_USERNAME_ALTERNATE_LABEL, global.TEST_VAR_USERNAME)
	os.Setenv(global.TEST_VAR_ENVIRONMENT_PASSWORD_ALTERNATE_LABEL, global.TEST_VAR_PASSWORD)
	os.Setenv(global.TEST_VAR_ENVIRONMENT_ATTRIBUTE_NAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	_, exists = os.LookupEnv(global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL)
	assert.True(exists)
	_, exists = os.LookupEnv(global.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL)
	assert.True(exists)
	_, exists = os.LookupEnv(global.TEST_VAR_ENVIRONMENT_USERNAME_ALTERNATE_LABEL)
	assert.True(exists)
	_, exists = os.LookupEnv(global.TEST_VAR_ENVIRONMENT_PASSWORD_ALTERNATE_LABEL)
	assert.True(exists)
	_, exists = os.LookupEnv(global.TEST_VAR_ENVIRONMENT_ATTRIBUTE_NAME_LABEL)
	assert.True(exists)

	CleanEnvironment([]string{
		global.TEST_VAR_ENVIRONMENT_ATTRIBUTE_NAME_LABEL,
		global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL,
		global.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL,
		global.TEST_VAR_ENVIRONMENT_USERNAME_ALTERNATE_LABEL,
		global.TEST_VAR_ENVIRONMENT_PASSWORD_ALTERNATE_LABEL,
	})

	_, exists = os.LookupEnv(global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL)
	assert.False(exists)
	_, exists = os.LookupEnv(global.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL)
	assert.False(exists)
	_, exists = os.LookupEnv(global.TEST_VAR_ENVIRONMENT_USERNAME_ALTERNATE_LABEL)
	assert.False(exists)
	_, exists = os.LookupEnv(global.TEST_VAR_ENVIRONMENT_PASSWORD_ALTERNATE_LABEL)
	assert.False(exists)
	_, exists = os.LookupEnv(global.TEST_VAR_ENVIRONMENT_ATTRIBUTE_NAME_LABEL)
	assert.False(exists)
}
