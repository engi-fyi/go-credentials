package factory

import (
	"github.com/engi-fyi/go-credentials/global"
	"github.com/rs/zerolog/log"
	"os"
	"testing"
)

func TestFactoryNew(t *testing.T) {
	assert, _ := global.InitTest(t)
	log.Info().Msg("Testing the creation of a new factory.")
	testFactory, factoryErr := New(global.TEST_VAR_APPLICATION_NAME)
	assert.NoError(factoryErr)
	assert.Equal(testFactory.ApplicationName, global.TEST_VAR_APPLICATION_NAME)
	assert.False(testFactory.UseEnvironment)
	assert.Equal(testFactory.OutputType, global.OUTPUT_TYPE_INI)
	assert.True(testFactory.Initialized)
	assert.Equal(testGetParentDirectory(testFactory.ApplicationName), testFactory.ParentDirectory)
	assert.DirExists(testFactory.ParentDirectory)
	os.RemoveAll(testFactory.ParentDirectory)
}

func TestFactoryNoApplicationName(t *testing.T) {
	assert, log := global.InitTest(t)
	log.Info().Msg("Testing to make sure a blank application name cannot be passed.")
	_, factoryErr := New("")
	assert.EqualError(factoryErr, ERR_APPLICATION_NAME_BLANK)
}

func TestSetOutputType(t *testing.T) {
	assert, log := global.InitTest(t)
	log.Info().Msg("Testing valid output types.")
	testFactory, factoryErr := New(global.TEST_VAR_APPLICATION_NAME)
	assert.NoError(factoryErr)
	assert.Equal(global.OUTPUT_TYPE_INI, testFactory.OutputType)

	for _, value := range []string{global.OUTPUT_TYPE_JSON, global.OUTPUT_TYPE_INI} {
		outputErr := testFactory.SetOutputType(value)
		assert.NoError(outputErr)
	}
}

func TestSetInvalidOutputType(t *testing.T) {
	assert, _ := global.InitTest(t)
	log.Info().Msg("Testing invalid output type.")
	testFactory, factoryErr := New(global.TEST_VAR_APPLICATION_NAME)
	assert.NoError(factoryErr)
	iotErr := testFactory.SetOutputType(global.OUTPUT_TYPE_INVALID)
	assert.EqualError(iotErr, ERR_INVALID_OUTPUT_TYPE)
}

func TestAlternateUsernameEmpty(t *testing.T) {
	assert, log := global.InitTest(t)
	log.Info().Msg("Testing blank alternate username.")
	testFactory, factoryErr := New(global.TEST_VAR_APPLICATION_NAME)
	assert.NoError(factoryErr)
	setErr := testFactory.SetAlternateUsername("")
	assert.EqualError(setErr, ERR_ALTERNATE_USERNAME_CANNOT_BE_BLANK)
	assert.Equal(global.TEST_VAR_USERNAME_LABEL, testFactory.GetAlternateUsername())
}

func TestAlternateUsername(t *testing.T) {
	assert, log := global.InitTest(t)
	log.Info().Msg("Testing setting an alternate username.")
	testFactory, factoryErr := New(global.TEST_VAR_APPLICATION_NAME)
	assert.NoError(factoryErr)
	testFactory.SetAlternateUsername(global.TEST_VAR_USERNAME_ALTERNATE_LABEL)
	assert.Equal(global.TEST_VAR_USERNAME_ALTERNATE_LABEL, testFactory.GetAlternateUsername())
	assert.Equal(global.TEST_VAR_PASSWORD_LABEL, testFactory.GetAlternatePassword())
	os.RemoveAll(testFactory.ParentDirectory)
}

func TestAlternatePasswordEmpty(t *testing.T) {
	assert, log := global.InitTest(t)
	log.Info().Msg("Testing blank alternate username.")
	testFactory, factoryErr := New(global.TEST_VAR_APPLICATION_NAME)
	assert.NoError(factoryErr)
	setErr := testFactory.SetAlternatePassword("")
	assert.EqualError(setErr, ERR_ALTERNATE_PASSWORD_CANNOT_BE_BLANK)
	assert.Equal(global.TEST_VAR_PASSWORD_LABEL, testFactory.GetAlternatePassword())
}

func TestAlternatePassword(t *testing.T) {
	assert, log := global.InitTest(t)
	log.Info().Msg("Testing setting an alternate password.")
	testFactory, factoryErr := New(global.TEST_VAR_APPLICATION_NAME)
	assert.NoError(factoryErr)
	testFactory.SetAlternatePassword(global.TEST_VAR_PASSWORD_ALTERNATE_LABEL)
	assert.Equal(global.TEST_VAR_PASSWORD_ALTERNATE_LABEL, testFactory.GetAlternatePassword())
	assert.Equal(global.TEST_VAR_USERNAME_LABEL, testFactory.GetAlternateUsername())
	os.RemoveAll(testFactory.ParentDirectory)
}

func TestAlternatesEmpty(t *testing.T) {
	assert, log := global.InitTest(t)
	log.Info().Msg("Testing blank alternates.")
	testFactory, factoryErr := New(global.TEST_VAR_APPLICATION_NAME)
	assert.NoError(factoryErr)

	setErr := testFactory.SetAlternates("", global.TEST_VAR_PASSWORD_ALTERNATE_LABEL)
	assert.EqualError(setErr, ERR_ALTERNATE_USERNAME_CANNOT_BE_BLANK)
	assert.Equal(global.TEST_VAR_USERNAME_LABEL, testFactory.GetAlternateUsername())
	assert.Equal(global.TEST_VAR_PASSWORD_ALTERNATE_LABEL, testFactory.GetAlternatePassword())
	setErr = testFactory.SetAlternates(global.TEST_VAR_USERNAME_ALTERNATE_LABEL, "")
	assert.EqualError(setErr, ERR_ALTERNATE_PASSWORD_CANNOT_BE_BLANK)
	assert.Equal(global.TEST_VAR_USERNAME_ALTERNATE_LABEL, testFactory.GetAlternateUsername())
	assert.Equal(global.TEST_VAR_PASSWORD_ALTERNATE_LABEL, testFactory.GetAlternatePassword())

	alternateUsername, alternatePassword := testFactory.GetAlternates()
	assert.Equal(global.TEST_VAR_USERNAME_ALTERNATE_LABEL, alternateUsername)
	assert.Equal(global.TEST_VAR_PASSWORD_ALTERNATE_LABEL, alternatePassword)
}

func TestAlternates(t *testing.T) {
	assert, log := global.InitTest(t)
	log.Info().Msg("Testing alternates.")
	testFactory, factoryErr := New(global.TEST_VAR_APPLICATION_NAME)
	assert.NoError(factoryErr)

	setErr := testFactory.SetAlternates(global.TEST_VAR_USERNAME_ALTERNATE_LABEL, global.TEST_VAR_PASSWORD_ALTERNATE_LABEL)
	assert.NoError(setErr)
	assert.Equal(global.TEST_VAR_USERNAME_ALTERNATE_LABEL, testFactory.GetAlternateUsername())
	assert.Equal(global.TEST_VAR_PASSWORD_ALTERNATE_LABEL, testFactory.GetAlternatePassword())

	alternateUsername, alternatePassword := testFactory.GetAlternates()
	assert.Equal(global.TEST_VAR_USERNAME_ALTERNATE_LABEL, alternateUsername)
	assert.Equal(global.TEST_VAR_PASSWORD_ALTERNATE_LABEL, alternatePassword)
	os.RemoveAll(testFactory.ParentDirectory)
}

func TestFactoryLogging(t *testing.T) {
	assert, log := global.InitTest(t)
	log.Info().Msg("Testing factory logging methods.")
	os.Unsetenv(global.LOG_LEVEL_ENVIRONMENT_KEY)
	testFactory, factoryErr := New(global.TEST_VAR_APPLICATION_NAME)
	os.Setenv(global.LOG_LEVEL_ENVIRONMENT_KEY, "trace")
	assert.NoError(factoryErr)
	assert.Equal("", testFactory.Log.GetLevel().String())

	testFactory, factoryErr = New(global.TEST_VAR_APPLICATION_NAME)
	os.Unsetenv(global.LOG_LEVEL_ENVIRONMENT_KEY)
	assert.NoError(factoryErr)
	assert.Equal("trace", testFactory.Log.GetLevel().String())

	testLogLevels := []string{
		"panic", "fatal", "error", "warn", "info", "debug", "trace",
	}

	for i := range testLogLevels {
		testFactory, factoryErr = New(global.TEST_VAR_APPLICATION_NAME)
		assert.NoError(factoryErr)
		testFactory.ModifyLogger(testLogLevels[i], true)
		assert.Equal(testLogLevels[i], testFactory.Log.GetLevel().String())
	}

	// Setting this again so that the test logging can continue on correctly.
	os.Setenv(global.LOG_LEVEL_ENVIRONMENT_KEY, "trace")
	os.RemoveAll(testFactory.ParentDirectory)
}

func testGetParentDirectory(applicationName string) string {
	homeDirectory, _ := os.UserHomeDir()
	return homeDirectory + "/." + applicationName + "/"
}
