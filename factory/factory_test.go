package factory

import (
	"github.com/engi-fyi/go-credentials/global"
	"os"
	"strings"
	"testing"
)

func TestFactoryNew(t *testing.T) {
	assert, _ := global.InitTest(t)
	_, factoryErr := New("", false)
	assert.EqualError(factoryErr, ERR_APPLICATION_NAME_BLANK)
	newFactory, factoryErr := New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)
	assert.Equal(newFactory.ApplicationName, global.TEST_VAR_APPLICATION_NAME)
	assert.False(newFactory.UseEnvironment)
	assert.Equal(newFactory.OutputType, global.OUTPUT_TYPE_INI)
	assert.True(newFactory.Initialized)
}

func TestFactoryInitialize(t *testing.T) {
	assert, log := global.InitTest(t)

	log.Info().Msg("Testing to make sure a blank application name cannot be passed.")
	testFactory, blankErr := New("", false)
	assert.EqualError(blankErr, ERR_APPLICATION_NAME_BLANK)

	log.Info().Msg("Creating and initializing a new Factory.")
	testFactory, newErr := New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(newErr)
	assert.True(testFactory.Initialized)

	log.Info().Msg("Checking to make sure that the parent directory was created.")
	assert.Equal(testGetParentDirectory(testFactory.ApplicationName), testFactory.ParentDirectory)
	assert.DirExists(testFactory.ParentDirectory)

	log.Info().Msg("Checking to make sure that the configuration directory was created.")
	assert.Equal(testGetConfigDirectory(testFactory.ApplicationName), testFactory.ConfigDirectory)
	assert.DirExists(testFactory.ConfigDirectory)

	os.RemoveAll(testFactory.ParentDirectory)
}

func TestSetOutputType(t *testing.T) {
	assert, _ := global.InitTest(t)

	newFactory, factoryErr := New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)

	nfiErr := newFactory.Initialize()
	assert.NoError(nfiErr)

	otjErr := newFactory.SetOutputType(global.OUTPUT_TYPE_JSON)
	assert.NoError(otjErr)

	otiErr := newFactory.SetOutputType(global.OUTPUT_TYPE_INI)
	assert.NoError(otiErr)

	iotErr := newFactory.SetOutputType(global.OUTPUT_TYPE_INVALID)
	assert.EqualError(iotErr, ERR_INVALID_OUTPUT_TYPE)
}

func TestSetEnvironmentKeys(t *testing.T) {
	assert, _ := global.InitTest(t)

	newFactory, factoryErr := New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)
	initErr := newFactory.Initialize()
	assert.NoError(initErr)

	altErr := newFactory.SetEnvironmentKeys(global.TEST_VAR_USERNAME_ALTERNATE_LABEL, global.TEST_VAR_PASSWORD_ALTERNATE_LABEL)
	assert.NoError(altErr)
	assert.EqualValues(global.TEST_VAR_USERNAME_ALTERNATE_LABEL, newFactory.GetAlternateUsername())
	assert.EqualValues(global.TEST_VAR_PASSWORD_ALTERNATE_LABEL, newFactory.GetAlternatePassword())

	altErr = newFactory.SetEnvironmentKeys(global.TEST_VAR_ENVIRONMENT_BAD_LABEL, global.TEST_VAR_USERNAME_ALTERNATE_LABEL)
	assert.EqualError(altErr, ERR_KEY_MUST_MATCH_REGEX)
	assert.EqualValues(global.TEST_VAR_USERNAME_ALTERNATE_LABEL, newFactory.GetAlternateUsername())

	altErr = newFactory.SetEnvironmentKeys(global.TEST_VAR_USERNAME_ALTERNATE_LABEL, global.TEST_VAR_ENVIRONMENT_BAD_LABEL)
	assert.EqualError(altErr, ERR_KEY_MUST_MATCH_REGEX)
	assert.EqualValues(global.TEST_VAR_PASSWORD_ALTERNATE_LABEL, newFactory.GetAlternatePassword())

	altErr = newFactory.SetEnvironmentKeys(global.TEST_VAR_ENVIRONMENT_BAD_LABEL, global.TEST_VAR_ENVIRONMENT_BAD_LABEL)
	assert.EqualError(altErr, ERR_KEY_MUST_MATCH_REGEX)
	assert.EqualValues(global.TEST_VAR_USERNAME_ALTERNATE_LABEL, newFactory.GetAlternateUsername())
	assert.EqualValues(global.TEST_VAR_PASSWORD_ALTERNATE_LABEL, newFactory.GetAlternatePassword())
}

func TestFactoryAlternates(t *testing.T) {
	assert, _ := global.InitTest(t)
	newFactory, factoryErr := New(global.TEST_VAR_APPLICATION_NAME, false)
	newFactory.Initialize()
	assert.NoError(factoryErr)

	assert.Equal("username", newFactory.GetAlternateUsername())
	assert.Equal("password", newFactory.GetAlternatePassword())

	bauError := newFactory.SetAlternateUsername("")
	assert.EqualError(bauError, ERR_ALTERNATE_USERNAME_CANNOT_BE_BLANK)
	bapError := newFactory.SetAlternatePassword("")
	assert.EqualError(bapError, ERR_ALTERNATE_PASSWORD_CANNOT_BE_BLANK)
	newFactory.SetAlternateUsername(global.TEST_VAR_USERNAME_ALTERNATE_LABEL)
	newFactory.SetAlternatePassword(global.TEST_VAR_PASSWORD_ALTERNATE_LABEL)
	assert.Equal(strings.ToLower(global.TEST_VAR_USERNAME_ALTERNATE_LABEL), newFactory.GetAlternateUsername())
	assert.Equal(strings.ToLower(global.TEST_VAR_PASSWORD_ALTERNATE_LABEL), newFactory.GetAlternatePassword())
}

func TestFactoryLogging(t *testing.T) {
	assert, _ := global.InitTest(t)
	os.Unsetenv(global.LOG_LEVEL_ENVIRONMENT_KEY)
	newFactory, factoryErr := New(global.TEST_VAR_APPLICATION_NAME, false)
	os.Setenv(global.LOG_LEVEL_ENVIRONMENT_KEY, "trace")
	assert.NoError(factoryErr)
	assert.Equal("info", newFactory.Log.GetLevel().String())

	newFactory, factoryErr = New(global.TEST_VAR_APPLICATION_NAME, false)
	os.Unsetenv(global.LOG_LEVEL_ENVIRONMENT_KEY)
	assert.NoError(factoryErr)
	assert.Equal("trace", newFactory.Log.GetLevel().String())

	newFactory, factoryErr = New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)
	newFactory.ModifyLogger("fatal", true)
	assert.Equal("fatal", newFactory.Log.GetLevel().String())

	// Setting this again so that the test logging can continue on correctly.
	os.Setenv(global.LOG_LEVEL_ENVIRONMENT_KEY, "trace")
}

func testGetParentDirectory(applicationName string) string {
	homeDirectory, _ := os.UserHomeDir()
	return homeDirectory + "/." + applicationName + "/"
}

func testGetConfigDirectory(applicationName string) string {
	return testGetParentDirectory(applicationName) + "config/"
}
