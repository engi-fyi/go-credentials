package environment

import (
	"github.com/HammoTime/go-credentials/test"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"testing"
)

func TestDeploySimple(t *testing.T) {
	assert := test.InitTest(t)

	prefix := strings.ToUpper(test.TEST_VAR_APPLICATION_NAME) + "_"

	log.Info().Msg("Testing deploy with no attributes or alternates.")
	setKeys, deployErr := Deploy(
		test.TEST_VAR_APPLICATION_NAME,
		test.TEST_VAR_USERNAME,
		test.TEST_VAR_PASSWORD,
		nil,
		nil)
	assert.NoError(deployErr)
	assert.Equal(test.TEST_VAR_USERNAME, os.Getenv(prefix+"USERNAME"))
	assert.Equal(test.TEST_VAR_PASSWORD, os.Getenv(prefix+"PASSWORD"))

	CleanEnvironment(setKeys)
}

func TestDeploySimpleNoPassword(t *testing.T) {
	assert := test.InitTest(t)

	prefix := strings.ToUpper(test.TEST_VAR_APPLICATION_NAME) + "_"

	log.Info().Msg("Testing deploy with no attributes or alternates.")
	setKeys, deployErr := Deploy(
		test.TEST_VAR_APPLICATION_NAME,
		test.TEST_VAR_USERNAME,
		"",
		nil,
		nil)
	assert.NoError(deployErr)
	_, passwordExists := os.LookupEnv(prefix + "PASSWORD")
	assert.Equal(test.TEST_VAR_USERNAME, os.Getenv(prefix+"USERNAME"))
	assert.False(passwordExists)

	CleanEnvironment(setKeys)
}

func TestDeployAttributes(t *testing.T) {
	assert := test.InitTest(t)

	prefix := strings.ToUpper(test.TEST_VAR_APPLICATION_NAME) + "_"
	attributes := make(map[string]string)
	attributes[test.TEST_VAR_ATTRIBUTE_NAME] = test.TEST_VAR_ATTRIBUTE_VALUE

	log.Info().Msg("Testing deploy with attributes and no alternates.")
	setKeys, deployErr := Deploy(
		test.TEST_VAR_APPLICATION_NAME,
		test.TEST_VAR_USERNAME,
		test.TEST_VAR_PASSWORD,
		nil,
		attributes)
	assert.NoError(deployErr)
	assert.Equal(test.TEST_VAR_ATTRIBUTE_VALUE, os.Getenv(prefix+strings.ToUpper(test.TEST_VAR_ATTRIBUTE_NAME)))

	CleanEnvironment(setKeys)
}

func TestDeployAlternates(t *testing.T) {
	assert := test.InitTest(t)

	prefix := strings.ToUpper(test.TEST_VAR_APPLICATION_NAME) + "_"
	alternates := make(map[string]string)
	alternates["username"] = test.TEST_VAR_ALTERNATE_USERNAME
	alternates["password"] = test.TEST_VAR_ALTERNATE_PASSWORD

	log.Info().Msg("Testing deploy with attributes and alternates.")

	setKeys, deployErr := Deploy(
		test.TEST_VAR_APPLICATION_NAME,
		test.TEST_VAR_USERNAME,
		test.TEST_VAR_PASSWORD,
		alternates,
		nil)
	assert.NoError(deployErr)
	assert.Equal(test.TEST_VAR_USERNAME, os.Getenv(prefix+strings.ToUpper(test.TEST_VAR_ALTERNATE_USERNAME)))
	assert.Equal(test.TEST_VAR_PASSWORD, os.Getenv(prefix+strings.ToUpper(test.TEST_VAR_ALTERNATE_PASSWORD)))

	CleanEnvironment(setKeys)
}

func TestDeployAttributesAndAlternates(t *testing.T) {
	assert := test.InitTest(t)

	prefix := strings.ToUpper(test.TEST_VAR_APPLICATION_NAME) + "_"
	alternates := make(map[string]string)
	alternates["username"] = test.TEST_VAR_ALTERNATE_USERNAME
	alternates["password"] = test.TEST_VAR_ALTERNATE_PASSWORD
	attributes := make(map[string]string)
	attributes[test.TEST_VAR_ATTRIBUTE_NAME] = test.TEST_VAR_ATTRIBUTE_VALUE

	log.Info().Msg("Testing deploy with attributes and alternates.")

	setKeys, deployErr := Deploy(
		test.TEST_VAR_APPLICATION_NAME,
		test.TEST_VAR_USERNAME,
		test.TEST_VAR_PASSWORD,
		alternates,
		attributes)
	assert.NoError(deployErr)
	assert.Equal(test.TEST_VAR_ATTRIBUTE_VALUE, os.Getenv(prefix+strings.ToUpper(test.TEST_VAR_ATTRIBUTE_NAME)))
	assert.Equal(test.TEST_VAR_USERNAME, os.Getenv(prefix+strings.ToUpper(test.TEST_VAR_ALTERNATE_USERNAME)))
	assert.Equal(test.TEST_VAR_PASSWORD, os.Getenv(prefix+strings.ToUpper(test.TEST_VAR_ALTERNATE_PASSWORD)))

	CleanEnvironment(setKeys)
}

func TestLoad(t *testing.T) {
	assert := test.InitTest(t)

	usernameKey := strings.ToUpper(test.TEST_VAR_APPLICATION_NAME + "_" + test.TEST_VAR_USERNAME_ENV_LABEL)
	passwordKey := strings.ToUpper(test.TEST_VAR_APPLICATION_NAME + "_" + test.TEST_VAR_PASSWORD_ENV_LABEL)
	attributeKey := strings.ToUpper(test.TEST_VAR_APPLICATION_NAME + "_" + test.TEST_VAR_ATTRIBUTE_NAME)
	alternateUsernameKey := strings.ToUpper(test.TEST_VAR_APPLICATION_NAME + "_" + test.TEST_VAR_ALTERNATE_USERNAME)
	alternatePasswordKey := strings.ToUpper(test.TEST_VAR_APPLICATION_NAME + "_" + test.TEST_VAR_ALTERNATE_PASSWORD)
	alternates := make(map[string]string)
	alternates[strings.ToLower(test.TEST_VAR_USERNAME_ENV_LABEL)] = test.TEST_VAR_ALTERNATE_USERNAME
	alternates[strings.ToLower(test.TEST_VAR_PASSWORD_ENV_LABEL)] = test.TEST_VAR_ALTERNATE_PASSWORD

	log.Trace().Msg("Testing basic load from environment variables.")
	_, loadErr := Load(test.TEST_VAR_APPLICATION_NAME, nil)
	assert.EqualError(loadErr, ERR_REQUIRED_VARIABLE_USERNAME_NOT_FOUND)
	os.Setenv(usernameKey, test.TEST_VAR_USERNAME)
	_, loadErr = Load(test.TEST_VAR_APPLICATION_NAME, nil)
	assert.EqualError(loadErr, ERR_REQUIRED_VARIABLE_PASSWORD_NOT_FOUND)
	os.Setenv(passwordKey, test.TEST_VAR_PASSWORD)
	values, loadErr := Load(test.TEST_VAR_APPLICATION_NAME, nil)
	assert.NoError(loadErr)
	assert.Equal(values["username"], test.TEST_VAR_USERNAME)
	assert.Equal(values["password"], test.TEST_VAR_PASSWORD)
	os.Setenv(attributeKey, test.TEST_VAR_ATTRIBUTE_VALUE)
	values, loadErr = Load(test.TEST_VAR_APPLICATION_NAME, nil)
	assert.NoError(loadErr)
	assert.Equal(test.TEST_VAR_ATTRIBUTE_VALUE, values[strings.ToLower(test.TEST_VAR_ATTRIBUTE_NAME)])
	CleanEnvironment([]string{usernameKey, passwordKey})

	log.Trace().Msg("Testing alternate load from environment variables.")
	// Test alternates
	os.Setenv(alternateUsernameKey, test.TEST_VAR_USERNAME)
	os.Setenv(alternatePasswordKey, test.TEST_VAR_PASSWORD)
	values, loadErr = Load(test.TEST_VAR_APPLICATION_NAME, alternates)
	assert.NoError(loadErr)
	assert.Equal(test.TEST_VAR_USERNAME, values[strings.ToLower(test.TEST_VAR_ALTERNATE_USERNAME)])
	assert.Equal(test.TEST_VAR_PASSWORD, values[strings.ToLower(test.TEST_VAR_ALTERNATE_PASSWORD)])
	CleanEnvironment([]string{usernameKey, passwordKey})
}

func TestCleanEnvironment(t *testing.T) {
	assert := test.InitTest(t)

	_, exists := os.LookupEnv(test.TEST_VAR_ALTERNATE_USERNAME)
	assert.False(exists)
	_, exists = os.LookupEnv(test.TEST_VAR_ALTERNATE_PASSWORD)
	assert.False(exists)
	_, exists = os.LookupEnv(test.TEST_VAR_ATTRIBUTE_NAME)
	assert.False(exists)

	os.Setenv(test.TEST_VAR_ALTERNATE_USERNAME, test.TEST_VAR_USERNAME)
	os.Setenv(test.TEST_VAR_ALTERNATE_PASSWORD, test.TEST_VAR_PASSWORD)
	os.Setenv(test.TEST_VAR_ATTRIBUTE_NAME, test.TEST_VAR_ATTRIBUTE_VALUE)
	_, exists = os.LookupEnv(test.TEST_VAR_ALTERNATE_USERNAME)
	assert.True(exists)
	_, exists = os.LookupEnv(test.TEST_VAR_ALTERNATE_PASSWORD)
	assert.True(exists)
	_, exists = os.LookupEnv(test.TEST_VAR_ATTRIBUTE_NAME)
	assert.True(exists)

	CleanEnvironment([]string{test.TEST_VAR_ALTERNATE_USERNAME, test.TEST_VAR_ALTERNATE_PASSWORD, test.TEST_VAR_ATTRIBUTE_NAME})
	_, exists = os.LookupEnv(test.TEST_VAR_ALTERNATE_USERNAME)
	assert.False(exists)
	_, exists = os.LookupEnv(test.TEST_VAR_ALTERNATE_PASSWORD)
	assert.False(exists)
	_, exists = os.LookupEnv(test.TEST_VAR_ATTRIBUTE_NAME)
	assert.False(exists)
}
