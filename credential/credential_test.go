package credential

import (
	"github.com/HammoTime/go-credentials/environment"
	"github.com/HammoTime/go-credentials/factory"
	"github.com/HammoTime/go-credentials/global"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
)

func TestCredentialNew(t *testing.T) {
	assert := global.InitTest(t)

	// Test basic info.

	log.Info().Msg("Testing to ensure an initialized factory is required.")
	blankFactory := factory.Factory{}
	_, factoryErr := New(&blankFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.EqualError(factoryErr, ERR_FACTORY_MUST_BE_INITIALIZED)
	_, factoryErr = New(nil, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.EqualError(factoryErr, ERR_FACTORY_MUST_BE_INITIALIZED)

	factory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)
	fiErr := factory.Initialize()
	assert.NoError(fiErr)
	assert.True(factory.Initialized)

	log.Info().Msg("Testing to ensure blank username or password can't be used.")
	_, missingErr := New(factory, "", global.TEST_VAR_PASSWORD)
	assert.EqualError(missingErr, ERR_USERNAME_OR_PASSWORD_NOT_SET)
	_, missingErr = New(factory, global.TEST_VAR_USERNAME, "")
	assert.EqualError(missingErr, ERR_USERNAME_OR_PASSWORD_NOT_SET)
	_, missingErr = New(factory, "", "")
	assert.EqualError(missingErr, ERR_USERNAME_OR_PASSWORD_NOT_SET)

	log.Info().Msg("Testing basic creation.")
	testCredential, newErr := New(factory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(newErr)
	assert.Equal(global.TEST_VAR_USERNAME, testCredential.Username)
	assert.Equal(global.TEST_VAR_PASSWORD, testCredential.Password)
	assert.True(testCredential.Factory.Initialized)
	assert.True(testCredential.Initialized)
}

func TestCredentialAttributes(t *testing.T) {
	assert := global.InitTest(t)

	// Test basic info.
	factory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)
	fiError := factory.Initialize()
	assert.NoError(fiError)
	assert.True(factory.Initialized)
	testCredential, newErr := New(factory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(newErr)

	// Username testing
	log.Info().Msg("Testing username accessibility through GetAttribute.")
	username, usernameErr := testCredential.GetAttribute("username")
	assert.Equal(global.TEST_VAR_USERNAME, username)
	assert.NoError(usernameErr)

	// Username testing
	log.Info().Msg("Testing password accessibility through GetAttribute.")
	password, passwordErr := testCredential.GetAttribute("password")
	assert.Equal(global.TEST_VAR_PASSWORD, password)
	assert.NoError(passwordErr)

	// Attribute that hasn't been set
	log.Info().Msg("Testing that correct error returned when an attribute not set and GetAttribute is used.")
	notSet, notSetErr := testCredential.GetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL)
	assert.Equal("", notSet)
	assert.EqualError(notSetErr, ERR_ATTRIBUTE_NOT_EXIST)

	// Setting an attribute
	log.Info().Msg("Testing the setting of an attribute.")
	testCredential.SetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	set, setErr := testCredential.GetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL)
	assert.Equal(global.TEST_VAR_ATTRIBUTE_VALUE, set)
	assert.NoError(setErr)

	// Testing bad attribute
	log.Info().Msg("Testing a bad attribute.")
	badSetErr := testCredential.SetAttribute(global.TEST_VAR_BAD_ATTRIBUTE_NAME, global.TEST_VAR_ATTRIBUTE_VALUE)
	assert.EqualError(badSetErr, ERR_KEY_MUST_MATCH_REGEX)
	badSetGet, bsgErr := testCredential.GetAttribute(global.TEST_VAR_BAD_ATTRIBUTE_NAME)
	assert.Equal("", badSetGet)
	assert.EqualError(bsgErr, ERR_ATTRIBUTE_NOT_EXIST)

	// Test username and password redirect
	setAppend := "setattr"
	testCredential.SetAttribute("username", global.TEST_VAR_USERNAME+setAppend)
	testCredential.SetAttribute("password", global.TEST_VAR_PASSWORD+setAppend)
	assert.Equal(testCredential.Username, global.TEST_VAR_USERNAME+setAppend)
	assert.Equal(testCredential.Password, global.TEST_VAR_PASSWORD+setAppend)
}

func TestCredentialSave(t *testing.T) {
	assert := global.InitTest(t)

	// Testing a factory that has been messed with
	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)
	testFactory.Initialize()
	testFactory.OutputType = global.OUTPUT_TYPE_INVALID
	messedWithCredential := Credential{}
	messedWithCredential.Factory = testFactory
	messedSave := messedWithCredential.Save()
	assert.EqualError(messedSave, factory.ERR_FACTORY_INCONSISTENT_STATE)

	// Reinitialise everything
	testFactory, factoryErr = factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)
	testFactory.Initialize()
	testCredential, newErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(newErr)

	testCredential.SetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	saveErr := testCredential.Save()
	log.Info().Msg("Testing that files save correctly.")
	assert.NoError(saveErr)

	global.TestCleanup(
		testCredential.Factory.ConfigurationDirectory,
		testCredential.Factory.CredentialFile)
}

func TestCredentialLoadFromIni(t *testing.T) {
	assert := global.InitTest(t)

	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)
	fiErr := testFactory.Initialize()
	assert.NoError(fiErr)
	assert.True(testFactory.Initialized)

	log.Info().Msg("Testing load with no file present.")
	_, failLoadErr := LoadFromIniFile(testFactory)
	assert.Error(failLoadErr)

	log.Info().Msg("Testing ini basic load.")
	testCredential, newErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(newErr)
	assert.Equal(global.TEST_VAR_USERNAME, testCredential.Username)
	assert.Equal(global.TEST_VAR_PASSWORD, testCredential.Password)
	assert.True(testCredential.Factory.Initialized)
	assert.True(testCredential.Initialized)

	log.Info().Msg("Saving temporary credentials")
	saveErr := testCredential.Save()
	assert.NoError(saveErr)
	testCredential = nil

	log.Info().Msg("Reloading temporary credentials")
	loadedCredential, loadErr := LoadFromIniFile(testFactory)
	assert.NoError(loadErr)
	assert.Equal(global.TEST_VAR_USERNAME, loadedCredential.Username)
	assert.Equal(global.TEST_VAR_PASSWORD, loadedCredential.Password)

	log.Info().Msg("Testing re-saving and adding attributes.")
	csError := loadedCredential.SetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	assert.NoError(csError)
	csError = loadedCredential.SetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL+"v2", global.TEST_VAR_ATTRIBUTE_VALUE+"v2")
	assert.NoError(csError)
	rsErr := loadedCredential.Save()
	assert.NoError(rsErr)

	loadedCredential, loadErr = LoadFromIniFile(testFactory)
	assert.NoError(loadErr)
	fav, favErr := loadedCredential.GetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL)
	sav, savErr := loadedCredential.GetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL + "v2")
	assert.NoError(favErr)
	assert.NoError(savErr)
	assert.EqualValues(global.TEST_VAR_ATTRIBUTE_VALUE, fav)
	assert.EqualValues(global.TEST_VAR_ATTRIBUTE_VALUE+"v2", sav)

	global.TestCleanup(
		loadedCredential.Factory.ConfigurationDirectory,
		loadedCredential.Factory.CredentialFile)
}

func TestCredentialLoad(t *testing.T) {
	assert := global.InitTest(t)

	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)
	fiErr := testFactory.Initialize()
	assert.NoError(fiErr)
	assert.True(testFactory.Initialized)

	log.Info().Msg("Testing load with no file or environment variables present.")
	_, failLoadErr := Load(testFactory)
	assert.Error(failLoadErr)

	log.Info().Msg("Testing basic environment load.")
	os.Setenv(global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL, global.TEST_VAR_USERNAME)
	os.Setenv(global.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL, global.TEST_VAR_PASSWORD)
	loadedCredential, envLoadErr := Load(testFactory)
	assert.NoError(envLoadErr)
	assert.Equal(global.TEST_VAR_USERNAME, loadedCredential.Username)
	assert.Equal(global.TEST_VAR_PASSWORD, loadedCredential.Password)

	log.Info().Msg("Testing to ensure environment variables are loaded before file credentials.")
	testCredential, newErr := New(testFactory, global.TEST_VAR_USERNAME_ALTERNATE, global.TEST_VAR_PASSWORD_ALTERNATE)
	assert.NoError(newErr)
	assert.Equal(global.TEST_VAR_USERNAME_ALTERNATE, testCredential.Username)
	assert.Equal(global.TEST_VAR_PASSWORD_ALTERNATE, testCredential.Password)
	assert.True(testCredential.Factory.Initialized)
	assert.True(testCredential.Initialized)

	log.Info().Msg("Saving file credentials")
	saveErr := testCredential.Save()
	assert.NoError(saveErr)

	log.Info().Msg("Testing that environment variables are loaded.")
	loadedCredential, envLoadErr = Load(testFactory)
	assert.NoError(envLoadErr)
	assert.Equal(global.TEST_VAR_USERNAME, loadedCredential.Username)
	assert.Equal(global.TEST_VAR_PASSWORD, loadedCredential.Password)

	environment.CleanEnvironment([]string{
		global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL,
		global.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL})

	log.Info().Msg("Testing that file credentials are loaded when no environment variables present.")
	loadedCredential, envLoadErr = Load(testFactory)
	assert.NoError(envLoadErr)
	assert.Equal(global.TEST_VAR_USERNAME_ALTERNATE, loadedCredential.Username)
	assert.Equal(global.TEST_VAR_PASSWORD_ALTERNATE, loadedCredential.Password)

	global.TestCleanup(
		testFactory.ConfigurationDirectory,
		testFactory.CredentialFile)
}

func TestCredentialDeployEnv(t *testing.T) {
	assert := global.InitTest(t)

	testCredentials, credErr := buildTestCredentials()
	assert.NoError(credErr)

	_, exists := os.LookupEnv(global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL)
	assert.False(exists)
	_, exists = os.LookupEnv(global.TEST_VAR_ENVIRONMENT_PASSWORD_ALTERNATE_LABEL)
	assert.False(exists)
	_, exists = os.LookupEnv(global.TEST_VAR_ENVIRONMENT_ATTRIBUTE_NAME_LABEL)
	assert.False(exists)

	testCredentials.SetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	deployErr := testCredentials.DeployEnv()
	assert.NoError(deployErr)

	_, exists = os.LookupEnv(global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL)
	assert.True(exists)
	_, exists = os.LookupEnv(global.TEST_VAR_ENVIRONMENT_PASSWORD_ALTERNATE_LABEL)
	assert.False(exists)
	_, exists = os.LookupEnv(global.TEST_VAR_ENVIRONMENT_ATTRIBUTE_NAME_LABEL)
	assert.True(exists)

	environment.CleanEnvironment(testCredentials.GetEnvironmentVariables())
}

// TODO: Implement tests for LoadFromEnvironment()
func TestCredentialLoadEnv(t *testing.T) {
	assert := global.InitTest(t)
	testFactory, _ := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	testFactory.Initialize()

	usernameKey := strings.ToUpper(global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL)
	passwordKey := strings.ToUpper(global.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL)
	attributeKey := strings.ToUpper(global.TEST_VAR_ENVIRONMENT_ATTRIBUTE_NAME_LABEL)
	alternateUsernameKey := strings.ToUpper(global.TEST_VAR_ENVIRONMENT_USERNAME_ALTERNATE_LABEL)
	alternatePasswordKey := strings.ToUpper(global.TEST_VAR_ENVIRONMENT_PASSWORD_ALTERNATE_LABEL)

	os.Setenv(usernameKey, global.TEST_VAR_USERNAME)
	os.Setenv(passwordKey, global.TEST_VAR_PASSWORD)
	os.Setenv(attributeKey, global.TEST_VAR_ATTRIBUTE_VALUE)

	testCredential, credErr := LoadFromEnvironment(testFactory)
	assert.NoError(credErr)
	assert.Equal(global.TEST_VAR_USERNAME, testCredential.Username)
	assert.Equal(global.TEST_VAR_PASSWORD, testCredential.Password)
	attrValue, _ := testCredential.GetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL)
	assert.Equal(global.TEST_VAR_ATTRIBUTE_VALUE, attrValue)
	assert.True(testCredential.Initialized)
	environment.CleanEnvironment([]string{usernameKey, passwordKey, attributeKey})

	testFactory.SetAlternateUsername(global.TEST_VAR_USERNAME_ALTERNATE_LABEL)
	testFactory.SetAlternatePassword(global.TEST_VAR_PASSWORD_ALTERNATE_LABEL)
	os.Setenv(alternateUsernameKey, global.TEST_VAR_USERNAME)
	os.Setenv(alternatePasswordKey, global.TEST_VAR_PASSWORD)
	testCredential, credErr = LoadFromEnvironment(testFactory)
	assert.NoError(credErr)
	assert.Equal(global.TEST_VAR_USERNAME, testCredential.Username)
	assert.Equal(global.TEST_VAR_PASSWORD, testCredential.Password)
	assert.True(testCredential.Initialized)
}

func buildTestCredentials() (*Credential, error) {
	factory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME, true)

	if factoryErr != nil {
		return nil, factoryErr
	}

	fiError := factory.Initialize()

	if fiError != nil {
		return nil, fiError
	}

	return New(factory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
}
