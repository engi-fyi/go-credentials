package credential

import (
	"github.com/HammoTime/go-credentials/environment"
	"github.com/HammoTime/go-credentials/factory"
	"github.com/HammoTime/go-credentials/test"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
)

func TestCredentialNew(t *testing.T) {
	assert := test.InitTest(t)

	// Test basic info.

	log.Info().Msg("Testing to ensure an initialized factory is required.")
	blankFactory := factory.Factory{}
	_, factoryErr := New(&blankFactory, test.TEST_VAR_USERNAME, test.TEST_VAR_PASSWORD)
	assert.EqualError(factoryErr, ERR_FACTORY_MUST_BE_INITIALIZED)
	_, factoryErr = New(nil, test.TEST_VAR_USERNAME, test.TEST_VAR_PASSWORD)
	assert.EqualError(factoryErr, ERR_FACTORY_MUST_BE_INITIALIZED)

	factory, factoryErr := factory.New(test.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)
	fiErr := factory.Initialize()
	assert.NoError(fiErr)
	assert.True(factory.Initialized)

	log.Info().Msg("Testing to ensure blank username or password can't be used.")
	_, missingErr := New(factory, "", test.TEST_VAR_PASSWORD)
	assert.EqualError(missingErr, ERR_USERNAME_OR_PASSWORD_NOT_SET)
	_, missingErr = New(factory, test.TEST_VAR_USERNAME, "")
	assert.EqualError(missingErr, ERR_USERNAME_OR_PASSWORD_NOT_SET)
	_, missingErr = New(factory, "", "")
	assert.EqualError(missingErr, ERR_USERNAME_OR_PASSWORD_NOT_SET)

	log.Info().Msg("Testing basic creation.")
	testCredential, newErr := New(factory, test.TEST_VAR_USERNAME, test.TEST_VAR_PASSWORD)
	assert.NoError(newErr)
	assert.Equal(test.TEST_VAR_USERNAME, testCredential.Username)
	assert.Equal(test.TEST_VAR_PASSWORD, testCredential.Password)
	assert.True(testCredential.Factory.Initialized)
	assert.True(testCredential.Initialized)
}

func TestCredentialAttributes(t *testing.T) {
	assert := test.InitTest(t)

	// Test basic info.
	factory, factoryErr := factory.New(test.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)
	fiError := factory.Initialize()
	assert.NoError(fiError)
	assert.True(factory.Initialized)
	testCredential, newErr := New(factory, test.TEST_VAR_USERNAME, test.TEST_VAR_PASSWORD)
	assert.NoError(newErr)

	// Username testing
	log.Info().Msg("Testing username accessibility through GetAttribute.")
	username, usernameErr := testCredential.GetAttribute("username")
	assert.Equal(test.TEST_VAR_USERNAME, username)
	assert.NoError(usernameErr)

	// Username testing
	log.Info().Msg("Testing password accessibility through GetAttribute.")
	password, passwordErr := testCredential.GetAttribute("password")
	assert.Equal(test.TEST_VAR_PASSWORD, password)
	assert.NoError(passwordErr)

	// Attribute that hasn't been set
	log.Info().Msg("Testing that correct error returned when an attribute not set and GetAttribute is used.")
	notSet, notSetErr := testCredential.GetAttribute(test.TEST_VAR_ATTRIBUTE_NAME)
	assert.Equal("", notSet)
	assert.EqualError(notSetErr, ERR_ATTRIBUTE_NOT_EXIST)

	// Setting an attribute
	log.Info().Msg("Testing the setting of an attribute.")
	testCredential.SetAttribute(test.TEST_VAR_ATTRIBUTE_NAME, test.TEST_VAR_ATTRIBUTE_VALUE)
	set, setErr := testCredential.GetAttribute(test.TEST_VAR_ATTRIBUTE_NAME)
	assert.Equal(test.TEST_VAR_ATTRIBUTE_VALUE, set)
	assert.NoError(setErr)

	// Testing bad attribute
	log.Info().Msg("Testing a bad attribute.")
	badSetErr := testCredential.SetAttribute(test.TEST_VAR_BAD_ATTRIBUTE_NAME, test.TEST_VAR_ATTRIBUTE_VALUE)
	assert.EqualError(badSetErr, ERR_KEY_MUST_MATCH_REGEX)
	badSetGet, bsgErr := testCredential.GetAttribute(test.TEST_VAR_BAD_ATTRIBUTE_NAME)
	assert.Equal("", badSetGet)
	assert.EqualError(bsgErr, ERR_ATTRIBUTE_NOT_EXIST)

	// Test username and password redirect
	setAppend := "setattr"
	testCredential.SetAttribute("username", test.TEST_VAR_USERNAME+setAppend)
	testCredential.SetAttribute("password", test.TEST_VAR_PASSWORD+setAppend)
	assert.Equal(testCredential.Username, test.TEST_VAR_USERNAME+setAppend)
	assert.Equal(testCredential.Password, test.TEST_VAR_PASSWORD+setAppend)
}

func TestCredentialSave(t *testing.T) {
	assert := test.InitTest(t)

	// Testing a factory that has been messed with
	testFactory, factoryErr := factory.New(test.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)
	testFactory.Initialize()
	testFactory.OutputType = factory.OUTPUT_TYPE_INVALID
	messedWithCredential := Credential{}
	messedWithCredential.Factory = testFactory
	messedSave := messedWithCredential.Save()
	assert.EqualError(messedSave, factory.ERR_FACTORY_INCONSISTENT_STATE)

	// Reinitialise everything
	testFactory, factoryErr = factory.New(test.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)
	testFactory.Initialize()
	testCredential, newErr := New(testFactory, test.TEST_VAR_USERNAME, test.TEST_VAR_PASSWORD)
	assert.NoError(newErr)

	testCredential.SetAttribute(test.TEST_VAR_ATTRIBUTE_NAME, test.TEST_VAR_ATTRIBUTE_VALUE)
	saveErr := testCredential.Save()
	log.Info().Msg("Testing that files save correctly.")
	assert.NoError(saveErr)

	test.TestCleanup(
		testCredential.Factory.ConfigurationDirectory,
		testCredential.Factory.CredentialFile)
}

func TestCredentialLoadFromIni(t *testing.T) {
	assert := test.InitTest(t)

	testFactory, factoryErr := factory.New(test.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)
	fiErr := testFactory.Initialize()
	assert.NoError(fiErr)
	assert.True(testFactory.Initialized)

	log.Info().Msg("Testing load with no file present.")
	_, failLoadErr := LoadFromIniFile(testFactory)
	assert.Error(failLoadErr)

	log.Info().Msg("Testing ini basic load.")
	testCredential, newErr := New(testFactory, test.TEST_VAR_USERNAME, test.TEST_VAR_PASSWORD)
	assert.NoError(newErr)
	assert.Equal(test.TEST_VAR_USERNAME, testCredential.Username)
	assert.Equal(test.TEST_VAR_PASSWORD, testCredential.Password)
	assert.True(testCredential.Factory.Initialized)
	assert.True(testCredential.Initialized)

	log.Info().Msg("Saving temporary credentials")
	saveErr := testCredential.Save()
	assert.NoError(saveErr)
	testCredential = nil

	log.Info().Msg("Reloading temporary credentials")
	loadedCredential, loadErr := LoadFromIniFile(testFactory)
	assert.NoError(loadErr)
	assert.Equal(test.TEST_VAR_USERNAME, loadedCredential.Username)
	assert.Equal(test.TEST_VAR_PASSWORD, loadedCredential.Password)

	log.Info().Msg("Testing re-saving and adding attributes.")
	csError := loadedCredential.SetAttribute(test.TEST_VAR_ATTRIBUTE_NAME, test.TEST_VAR_ATTRIBUTE_VALUE)
	assert.NoError(csError)
	csError = loadedCredential.SetAttribute(test.TEST_VAR_ATTRIBUTE_NAME+"v2", test.TEST_VAR_ATTRIBUTE_VALUE+"v2")
	assert.NoError(csError)
	rsErr := loadedCredential.Save()
	assert.NoError(rsErr)

	loadedCredential, loadErr = LoadFromIniFile(testFactory)
	assert.NoError(loadErr)
	fav, favErr := loadedCredential.GetAttribute(test.TEST_VAR_ATTRIBUTE_NAME)
	sav, savErr := loadedCredential.GetAttribute(test.TEST_VAR_ATTRIBUTE_NAME + "v2")
	assert.NoError(favErr)
	assert.NoError(savErr)
	assert.EqualValues(test.TEST_VAR_ATTRIBUTE_VALUE, fav)
	assert.EqualValues(test.TEST_VAR_ATTRIBUTE_VALUE+"v2", sav)

	test.TestCleanup(
		loadedCredential.Factory.ConfigurationDirectory,
		loadedCredential.Factory.CredentialFile)
}

func TestCredentialLoad(t *testing.T) {
	assert := test.InitTest(t)

	testFactory, factoryErr := factory.New(test.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)
	fiErr := testFactory.Initialize()
	assert.NoError(fiErr)
	assert.True(testFactory.Initialized)

	log.Info().Msg("Testing load with no file or environment variables present.")
	_, failLoadErr := Load(testFactory)
	assert.Error(failLoadErr)

	log.Info().Msg("Testing basic environment load.")
	os.Setenv(test.TEST_VAR_ENVIRONMENT_USERNAME_LABEL, test.TEST_VAR_USERNAME)
	os.Setenv(test.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL, test.TEST_VAR_PASSWORD)
	loadedCredential, envLoadErr := Load(testFactory)
	assert.NoError(envLoadErr)
	assert.Equal(test.TEST_VAR_USERNAME, loadedCredential.Username)
	assert.Equal(test.TEST_VAR_PASSWORD, loadedCredential.Password)

	log.Info().Msg("Testing to ensure environment variables are loaded before file credentials.")
	testCredential, newErr := New(testFactory, test.TEST_VAR_USERNAME_ALTERNATE, test.TEST_VAR_PASSWORD_ALTERNATE)
	assert.NoError(newErr)
	assert.Equal(test.TEST_VAR_USERNAME_ALTERNATE, testCredential.Username)
	assert.Equal(test.TEST_VAR_PASSWORD_ALTERNATE, testCredential.Password)
	assert.True(testCredential.Factory.Initialized)
	assert.True(testCredential.Initialized)

	log.Info().Msg("Saving file credentials")
	saveErr := testCredential.Save()
	assert.NoError(saveErr)

	log.Info().Msg("Testing that environment variables are loaded.")
	loadedCredential, envLoadErr = Load(testFactory)
	assert.NoError(envLoadErr)
	assert.Equal(test.TEST_VAR_USERNAME, loadedCredential.Username)
	assert.Equal(test.TEST_VAR_PASSWORD, loadedCredential.Password)

	environment.CleanEnvironment([]string{
		test.TEST_VAR_ENVIRONMENT_USERNAME_LABEL,
		test.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL})

	log.Info().Msg("Testing that file credentials are loaded when no environment variables present.")
	loadedCredential, envLoadErr = Load(testFactory)
	assert.NoError(envLoadErr)
	assert.Equal(test.TEST_VAR_USERNAME_ALTERNATE, loadedCredential.Username)
	assert.Equal(test.TEST_VAR_PASSWORD_ALTERNATE, loadedCredential.Password)

	test.TestCleanup(
		testFactory.ConfigurationDirectory,
		testFactory.CredentialFile)
}

func TestCredentialDeployEnv(t *testing.T) {
	assert := test.InitTest(t)

	testCredentials, credErr := buildTestCredentials()
	prefix := strings.ToUpper(test.TEST_VAR_APPLICATION_NAME + "_")
	assert.NoError(credErr)

	_, exists := os.LookupEnv(prefix + test.TEST_VAR_USERNAME_ENV_LABEL)
	assert.False(exists)
	_, exists = os.LookupEnv(prefix + test.TEST_VAR_PASSWORD_ENV_LABEL)
	assert.False(exists)
	_, exists = os.LookupEnv(prefix + strings.ToUpper(test.TEST_VAR_ATTRIBUTE_NAME))
	assert.False(exists)

	testCredentials.SetAttribute(test.TEST_VAR_ATTRIBUTE_NAME, test.TEST_VAR_ATTRIBUTE_VALUE)
	deployErr := testCredentials.DeployEnv()
	assert.NoError(deployErr)

	_, exists = os.LookupEnv(prefix + test.TEST_VAR_USERNAME_ENV_LABEL)
	assert.True(exists)
	_, exists = os.LookupEnv(prefix + test.TEST_VAR_PASSWORD_ENV_LABEL)
	assert.False(exists)
	_, exists = os.LookupEnv(prefix + strings.ToUpper(test.TEST_VAR_ATTRIBUTE_NAME))
	assert.True(exists)

	environment.CleanEnvironment(testCredentials.GetEnvironmentVariables())
}

// TODO: Implement tests for LoadFromEnvironment()
func TestCredentialLoadEnv(t *testing.T) {
	assert := test.InitTest(t)
	testFactory, _ := factory.New(test.TEST_VAR_APPLICATION_NAME, false)
	testFactory.Initialize()

	usernameKey := strings.ToUpper(test.TEST_VAR_APPLICATION_NAME + "_" + test.TEST_VAR_USERNAME_ENV_LABEL)
	passwordKey := strings.ToUpper(test.TEST_VAR_APPLICATION_NAME + "_" + test.TEST_VAR_PASSWORD_ENV_LABEL)
	attributeKey := strings.ToUpper(test.TEST_VAR_APPLICATION_NAME + "_" + test.TEST_VAR_ATTRIBUTE_NAME)
	alternateUsernameKey := strings.ToUpper(test.TEST_VAR_APPLICATION_NAME + "_" + test.TEST_VAR_ALTERNATE_USERNAME)
	alternatePasswordKey := strings.ToUpper(test.TEST_VAR_APPLICATION_NAME + "_" + test.TEST_VAR_ALTERNATE_PASSWORD)

	os.Setenv(usernameKey, test.TEST_VAR_USERNAME)
	os.Setenv(passwordKey, test.TEST_VAR_PASSWORD)
	os.Setenv(attributeKey, test.TEST_VAR_ATTRIBUTE_VALUE)

	testCredential, credErr := LoadFromEnvironment(testFactory)
	assert.NoError(credErr)
	assert.Equal(test.TEST_VAR_USERNAME, testCredential.Username)
	assert.Equal(test.TEST_VAR_PASSWORD, testCredential.Password)
	attrValue, _ := testCredential.GetAttribute(test.TEST_VAR_ATTRIBUTE_NAME)
	assert.Equal(test.TEST_VAR_ATTRIBUTE_VALUE, attrValue)
	assert.True(testCredential.Initialized)
	environment.CleanEnvironment([]string{usernameKey, passwordKey, attributeKey})

	testFactory.SetAlternateUsername(test.TEST_VAR_ALTERNATE_USERNAME)
	testFactory.SetAlternatePassword(test.TEST_VAR_ALTERNATE_PASSWORD)
	os.Setenv(alternateUsernameKey, test.TEST_VAR_USERNAME)
	os.Setenv(alternatePasswordKey, test.TEST_VAR_PASSWORD)
	testCredential, credErr = LoadFromEnvironment(testFactory)
	assert.NoError(credErr)
	assert.Equal(test.TEST_VAR_USERNAME, testCredential.Username)
	assert.Equal(test.TEST_VAR_PASSWORD, testCredential.Password)
	assert.True(testCredential.Initialized)
}

func buildTestCredentials() (*Credential, error) {
	factory, factoryErr := factory.New(test.TEST_VAR_APPLICATION_NAME, true)

	if factoryErr != nil {
		return nil, factoryErr
	}

	fiError := factory.Initialize()

	if fiError != nil {
		return nil, fiError
	}

	return New(factory, test.TEST_VAR_USERNAME, test.TEST_VAR_PASSWORD)
}
