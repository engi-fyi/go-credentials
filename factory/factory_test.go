package factory

import (
	"github.com/HammoTime/go-credentials/test"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"testing"
)

func TestFactoryNew(t *testing.T) {
	assert := test.InitTest(t)
	_, factoryErr := New("", false)
	assert.EqualError(factoryErr, ERR_APPLICATION_NAME_BLANK)
	newFactory, factoryErr := New(test.TEST_VAR_APPLICATION_NAME, false)
	assert.Equal(newFactory.ApplicationName, test.TEST_VAR_APPLICATION_NAME)
	assert.False(newFactory.UseEnvironment)
	assert.Equal(newFactory.OutputType, OUTPUT_TYPE_INI)
}

func TestFactoryInitialize(t *testing.T) {
	assert := test.InitTest(t)

	log.Info().Msg("Testing to make sure a blank application name cannot be passed.")
	testFactory, blankErr := New("", false)
	assert.EqualError(blankErr, ERR_APPLICATION_NAME_BLANK)

	log.Info().Msg("Creating and initializing a new Factory.")
	testFactory, newErr := New(test.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(newErr)
	assert.False(testFactory.Initialized)
	testFactory.Initialize()
	assert.True(testFactory.Initialized)

	log.Info().Msg("Checking to make sure that the configuration directory was created.")
	assert.Equal(testGetConfigurationDirectory(testFactory.ApplicationName), testFactory.ConfigurationDirectory)
	assert.DirExists(testFactory.ConfigurationDirectory)

	cleanupErr := test.TestCleanup(testFactory.ConfigurationDirectory, testFactory.CredentialFile)
	assert.NoError(cleanupErr)
}

func TestSetOutputType(t *testing.T) {
	assert := test.InitTest(t)

	newFactory, factoryErr := New(test.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)

	nfiErr := newFactory.Initialize()
	assert.NoError(nfiErr)

	otjErr := newFactory.SetOutputType(OUTPUT_TYPE_JSON)
	assert.NoError(otjErr)

	otiErr := newFactory.SetOutputType(OUTPUT_TYPE_INI)
	assert.NoError(otiErr)

	iotErr := newFactory.SetOutputType(OUTPUT_TYPE_INVALID)
	assert.EqualError(iotErr, ERR_INVALID_OUTPUT_TYPE)
}

func TestSetEnvironmentKeys(t *testing.T) {
	assert := test.InitTest(t)

	newFactory, factoryErr := New(test.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)
	initErr := newFactory.Initialize()
	assert.NoError(initErr)

	altErr := newFactory.SetEnvironmentKeys(test.TEST_VAR_ENVIRONMENT_USERNAME_KEY, test.TEST_VAR_ENVIRONMENT_PASSWORD_KEY)
	assert.NoError(altErr)
	assert.EqualValues(test.TEST_VAR_ENVIRONMENT_USERNAME_KEY, newFactory.Alternates["username"])
	assert.EqualValues(test.TEST_VAR_ENVIRONMENT_PASSWORD_KEY, newFactory.Alternates["password"])

	altErr = newFactory.SetEnvironmentKeys(test.TEST_VAR_ENVIRONMENT_BAD_KEY, test.TEST_VAR_ENVIRONMENT_PASSWORD_KEY)
	assert.EqualError(altErr, ERR_KEY_MUST_MATCH_REGEX)
	assert.EqualValues(test.TEST_VAR_ENVIRONMENT_USERNAME_KEY, newFactory.Alternates["username"])

	altErr = newFactory.SetEnvironmentKeys(test.TEST_VAR_ENVIRONMENT_USERNAME_KEY, test.TEST_VAR_ENVIRONMENT_BAD_KEY)
	assert.EqualError(altErr, ERR_KEY_MUST_MATCH_REGEX)
	assert.EqualValues(test.TEST_VAR_ENVIRONMENT_PASSWORD_KEY, newFactory.Alternates["password"])

	altErr = newFactory.SetEnvironmentKeys(test.TEST_VAR_ENVIRONMENT_BAD_KEY, test.TEST_VAR_ENVIRONMENT_BAD_KEY)
	assert.EqualError(altErr, ERR_KEY_MUST_MATCH_REGEX)
	assert.EqualValues(test.TEST_VAR_ENVIRONMENT_USERNAME_KEY, newFactory.Alternates["username"])
	assert.EqualValues(test.TEST_VAR_ENVIRONMENT_USERNAME_KEY, newFactory.Alternates["username"])
	assert.EqualValues(test.TEST_VAR_ENVIRONMENT_PASSWORD_KEY, newFactory.Alternates["password"])
}

func TestFactoryAlternates(t *testing.T) {
	assert := test.InitTest(t)
	newFactory, factoryErr := New(test.TEST_VAR_APPLICATION_NAME, false)
	newFactory.Initialize()
	assert.NoError(factoryErr)

	assert.Equal("username", newFactory.GetAlternateUsername())
	assert.Equal("password", newFactory.GetAlternatePassword())

	bauError := newFactory.SetAlternateUsername("")
	assert.EqualError(bauError, ERR_ALTERNATE_USERNAME_CANNOT_BE_BLANK)
	bapError := newFactory.SetAlternatePassword("")
	assert.EqualError(bapError, ERR_ALTERNATE_PASSWORD_CANNOT_BE_BLANK)
	newFactory.SetAlternateUsername(test.TEST_VAR_ALTERNATE_USERNAME)
	newFactory.SetAlternatePassword(test.TEST_VAR_ALTERNATE_PASSWORD)
	assert.Equal(strings.ToLower(test.TEST_VAR_ALTERNATE_USERNAME), newFactory.GetAlternateUsername())
	assert.Equal(strings.ToLower(test.TEST_VAR_ALTERNATE_PASSWORD), newFactory.GetAlternatePassword())
}

func testGetConfigurationDirectory(applicationName string) string {
	homeDirectory, _ := os.UserHomeDir()
	return homeDirectory + "/." + applicationName + "/"
}

func testGetCredentialFile(configurationDirectory string) string {
	return configurationDirectory + "credentials"
}
