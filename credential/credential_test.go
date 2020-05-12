package credential

import (
	"github.com/engi-fyi/go-credentials/factory"
	"github.com/engi-fyi/go-credentials/global"
	"github.com/engi-fyi/go-credentials/profile"
	"github.com/engi-fyi/go-credentials/serializer"
	"os"
	"testing"
)

func TestCredentialNew(t *testing.T) {
	assert, log := global.InitTest(t)

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
	assert, log := global.InitTest(t)

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
	username := testCredential.GetAttribute("username")
	assert.Equal(global.TEST_VAR_USERNAME, username)

	// Username testing
	log.Info().Msg("Testing password accessibility through GetAttribute.")
	password := testCredential.GetAttribute("password")
	assert.Equal(global.TEST_VAR_PASSWORD, password)

	// Attribute that hasn't been set
	log.Info().Msg("Testing that correct error returned when an attribute not set and GetAttribute is used.")
	notSet := testCredential.GetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL)
	assert.Equal("", notSet)

	// Setting an attribute
	log.Info().Msg("Testing the setting of an attribute.")
	testCredential.SetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	set := testCredential.GetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL)
	assert.Equal(global.TEST_VAR_ATTRIBUTE_VALUE, set)

	// Testing bad attribute
	log.Info().Msg("Testing a bad attribute.")
	badSetErr := testCredential.SetAttribute(global.TEST_VAR_BAD_ATTRIBUTE_NAME, global.TEST_VAR_ATTRIBUTE_VALUE)
	assert.EqualError(badSetErr, ERR_KEY_MUST_MATCH_REGEX)
	badSetGet := testCredential.GetAttribute(global.TEST_VAR_BAD_ATTRIBUTE_NAME)
	assert.Equal("", badSetGet)

	// Test username and password redirect
	setAppend := "setattr"
	testCredential.SetAttribute("username", global.TEST_VAR_USERNAME+setAppend)
	testCredential.SetAttribute("password", global.TEST_VAR_PASSWORD+setAppend)
	assert.Equal(testCredential.Username, global.TEST_VAR_USERNAME+setAppend)
	assert.Equal(testCredential.Password, global.TEST_VAR_PASSWORD+setAppend)
}

func TestCredentialSave(t *testing.T) {
	assert, log := global.InitTest(t)

	// Reinitialise everything
	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)
	testFactory.Initialize()
	testCredential, newErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(newErr)

	testCredential.SetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	saveErr := testCredential.Save()
	log.Info().Msg("Testing that files save correctly.")
	assert.NoError(saveErr)

	os.RemoveAll(testCredential.Factory.ParentDirectory)
}

func TestCredentialInterferedWithSave(t *testing.T) {
	assert, _ := global.InitTest(t)

	// Testing a factory that has been messed with
	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)
	testFactory.Initialize()
	testFactory.OutputType = global.OUTPUT_TYPE_INVALID
	messedWithCredential := Credential{}
	messedWithCredential.Factory = testFactory
	messedSave := messedWithCredential.Save()
	assert.EqualError(messedSave, ERR_NOT_INITIALIZED)
}

func TestCredentialSaveProfileNotInitialized(t *testing.T) {
	assert, log := global.InitTest(t)

	// Reinitialise everything
	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)
	testFactory.Initialize()
	testCredential, newErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(newErr)
	testCredential.Factory = &factory.Factory{}
	saveErr := testCredential.Save()
	log.Info().Msg("Testing that files save correctly.")
	assert.Error(saveErr, factory.ERR_FACTORY_NOT_INITIALIZED)

	os.RemoveAll(testCredential.Factory.ParentDirectory)
}

func TestCredentialSaveFactoryInconsistent(t *testing.T) {
	assert, log := global.InitTest(t)

	// Reinitialise everything
	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)
	testFactory.Initialize()
	testCredential, newErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(newErr)
	testCredential.Factory.OutputType = global.OUTPUT_TYPE_INVALID
	saveErr := testCredential.Save()
	log.Info().Msg("Testing that files save correctly.")
	assert.Error(saveErr, factory.ERR_FACTORY_NOT_INITIALIZED)

	os.RemoveAll(testCredential.Factory.ParentDirectory)
}

func TestCredentialProfileNotInitialized(t *testing.T) {
	assert, log := global.InitTest(t)

	// Reinitialise everything
	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)
	testFactory.Initialize()
	testCredential, newErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(newErr)
	testCredential.Profile = &profile.Profile{}
	saveErr := testCredential.Save()
	log.Info().Msg("Testing that files save correctly.")
	assert.Error(saveErr, profile.ERR_PROFILE_NOT_INITIALIZED)

	os.RemoveAll(testCredential.Factory.ParentDirectory)
}

func TestCredentialLoadFromIni(t *testing.T) {
	assert, log := global.InitTest(t)

	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)
	fiErr := testFactory.Initialize()
	assert.NoError(fiErr)
	assert.True(testFactory.Initialized)

	log.Info().Msg("Testing load with no file present.")
	_, failLoadErr := LoadFromProfile(global.DEFAULT_PROFILE_NAME, testFactory)
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
	loadedCredential, loadErr := LoadFromProfile(global.DEFAULT_PROFILE_NAME, testFactory)
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

	loadedCredential, loadErr = LoadFromProfile(global.DEFAULT_PROFILE_NAME, testFactory)
	assert.NoError(loadErr)
	fav := loadedCredential.GetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL)
	sav := loadedCredential.GetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL + "v2")
	assert.EqualValues(global.TEST_VAR_ATTRIBUTE_VALUE, fav)
	assert.EqualValues(global.TEST_VAR_ATTRIBUTE_VALUE+"v2", sav)

	os.RemoveAll(loadedCredential.Factory.ParentDirectory)
}

func TestCredentialLoadFactoryNotInitialized(t *testing.T) {
	assert, _ := global.InitTest(t)

	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)
	testFactory.OutputType = global.OUTPUT_TYPE_INVALID
	_, loadErr := Load(testFactory)
	assert.EqualError(loadErr, serializer.ERR_UNRECOGNIZED_OUTPUT_TYPE)
}

func TestCredentialLoadNoIniFile(t *testing.T) {
	assert, _ := global.InitTest(t)

	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)
	os.RemoveAll(testFactory.ParentDirectory)
	_, loadErr := Load(testFactory)
	assert.Error(loadErr)
}

func TestCredentialLoadProfileNotInitialized(t *testing.T) {
	assert, _ := global.InitTest(t)

	testFactory := &factory.Factory{}
	_, loadErr := Load(testFactory)
	assert.EqualError(loadErr, ERR_FACTORY_MUST_BE_INITIALIZED)
}

func TestCredentialLoad(t *testing.T) {
	assert, log := global.InitTest(t)

	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)
	fiErr := testFactory.Initialize()
	assert.NoError(fiErr)
	assert.True(testFactory.Initialized)

	log.Info().Msg("Testing load with no file or environment variables present.")
	_, failLoadErr := Load(testFactory)
	assert.Error(failLoadErr)

	os.Unsetenv(global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL)
	os.Unsetenv(global.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL)
	log.Info().Msg("Testing basic environment load.")
	os.Setenv(global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL, global.TEST_VAR_USERNAME)
	os.Setenv(global.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL, global.TEST_VAR_PASSWORD)
	testFactory.SetOutputType(global.OUTPUT_TYPE_ENV)
	loadedCredential, envLoadErr := Load(testFactory)
	assert.NoError(envLoadErr)
	assert.Equal(global.TEST_VAR_USERNAME, loadedCredential.Username)
	assert.Equal(global.TEST_VAR_PASSWORD, loadedCredential.Password)
	os.Unsetenv(global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL)
	os.Unsetenv(global.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL)

	log.Info().Msg("Saving credential to test load.")
	testCredential, newErr := New(testFactory, global.TEST_VAR_USERNAME_ALTERNATE, global.TEST_VAR_PASSWORD_ALTERNATE)
	assert.NoError(newErr)
	assert.Equal(global.TEST_VAR_USERNAME_ALTERNATE, testCredential.Username)
	assert.Equal(global.TEST_VAR_PASSWORD_ALTERNATE, testCredential.Password)
	assert.True(testCredential.Factory.Initialized)
	assert.True(testCredential.Initialized)

	log.Info().Msg("Saving file credentials")
	saveErr := testCredential.Save()
	assert.NoError(saveErr)

	log.Info().Msg("Testing that file credentials are loaded when no environment variables present.")
	loadedCredential, envLoadErr = Load(testFactory)
	assert.NoError(envLoadErr)
	assert.Equal(global.TEST_VAR_USERNAME_ALTERNATE, loadedCredential.Username)
	assert.Equal(global.TEST_VAR_PASSWORD_ALTERNATE, loadedCredential.Password)

	os.RemoveAll(testCredential.Factory.ParentDirectory)
}

func TestCredentialSaveEnv(t *testing.T) {
	assert, _ := global.InitTest(t)
	os.Unsetenv(global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL)
	os.Unsetenv(global.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL)
	os.Unsetenv(global.TEST_VAR_ENVIRONMENT_ATTRIBUTE_NAME_LABEL)

	testCredentials, credErr := buildTestCredentials()
	testCredentials.Factory.SetOutputType(global.OUTPUT_TYPE_ENV)
	assert.NoError(credErr)

	_, exists := os.LookupEnv(global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL)
	assert.False(exists)
	_, exists = os.LookupEnv(global.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL)
	assert.False(exists)
	_, exists = os.LookupEnv(global.TEST_VAR_ENVIRONMENT_ATTRIBUTE_NAME_LABEL)
	assert.False(exists)

	//testCredentials.SetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	testCredentials.Section(global.TEST_VAR_FIRST_SECTION_KEY).SetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	deployErr := testCredentials.Save()
	assert.NoError(deployErr)

	_, exists = os.LookupEnv(global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL)
	assert.True(exists)
	_, exists = os.LookupEnv(global.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL)
	assert.True(exists)
	_, exists = os.LookupEnv(global.TEST_VAR_ENVIRONMENT_ATTRIBUTE_NAME_LABEL)
	assert.True(exists)

	os.Unsetenv(global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL)
	os.Unsetenv(global.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL)
	os.Unsetenv(global.TEST_VAR_ENVIRONMENT_ATTRIBUTE_NAME_LABEL)
}

func TestCredentialProfiles(t *testing.T) {
	assert, _ := global.InitTest(t)
	testFactory, _ := factory.New(global.TEST_VAR_APPLICATION_NAME, false)

	testCredential, tcErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(tcErr)
	testCredential.SetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	testCredential.Save()

	credentialOne, coErr := NewProfile(global.TEST_VAR_FIRST_PROFILE_LABEL, testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(coErr)
	credentialOne.SetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	credentialOne.Save()

	credentialTwo, ctErr := NewProfile(global.TEST_VAR_SECOND_PROFILE_LABEL, testFactory, global.TEST_VAR_USERNAME_ALTERNATE, global.TEST_VAR_PASSWORD_ALTERNATE)
	assert.NoError(ctErr)
	credentialTwo.SetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	credentialTwo.Save()

	testCredential, tcErr = Load(testFactory)
	assert.NoError(tcErr)
	assert.Equal(global.DEFAULT_PROFILE_NAME, testCredential.Profile.Name)
	attrValue := testCredential.GetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL)
	assert.Equal(global.TEST_VAR_ATTRIBUTE_VALUE, attrValue)
	assert.True(testCredential.Initialized)
	assert.True(testCredential.Profile.Initialized)

	credentialOne, coErr = LoadFromProfile(global.TEST_VAR_FIRST_PROFILE_LABEL, testFactory)
	assert.NoError(coErr)
	assert.Equal(global.TEST_VAR_FIRST_PROFILE_LABEL, credentialOne.Profile.Name)
	attrValue = credentialOne.GetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL)
	assert.Equal(global.TEST_VAR_ATTRIBUTE_VALUE, attrValue)
	assert.Equal(global.TEST_VAR_USERNAME, credentialOne.Username)
	assert.Equal(global.TEST_VAR_PASSWORD, credentialOne.Password)
	assert.True(credentialOne.Initialized)
	assert.True(credentialOne.Profile.Initialized)

	credentialTwo, ctErr = LoadFromProfile(global.TEST_VAR_SECOND_PROFILE_LABEL, testFactory)
	assert.NoError(ctErr)
	assert.Equal(global.TEST_VAR_SECOND_PROFILE_LABEL, credentialTwo.Profile.Name)
	attrValue = credentialTwo.GetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL)
	assert.Equal(global.TEST_VAR_ATTRIBUTE_VALUE, attrValue)
	assert.Equal(global.TEST_VAR_USERNAME_ALTERNATE, credentialTwo.Username)
	assert.Equal(global.TEST_VAR_PASSWORD_ALTERNATE, credentialTwo.Password)
	assert.True(credentialTwo.Initialized)
	assert.True(credentialTwo.Profile.Initialized)

	os.RemoveAll(testFactory.ParentDirectory)
}

/*
This tests:
 - Whether a value is set correctly in selectedSection of the returned Credential.
 - Whether a cloned section is correctly a shallow copy of the initial Credential.
 */
func TestSectionCredentialWithValue(t *testing.T) {
	assert, _ := global.InitTest(t)
	testFactory, _ := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	testCredential, tcErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(tcErr)

	sectionCredential := testCredential.Section(global.TEST_VAR_FIRST_SECTION_KEY)
	assert.NotEqual(testCredential.selectedSection, sectionCredential.selectedSection)
	assert.NotEqual("", sectionCredential.selectedSection)
}

/*
This tests:
 - Whether a blank section name is handled correctly.
*/
func TestSectionBlank(t *testing.T) {
	assert, _ := global.InitTest(t)
	testFactory, _ := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	testCredential, tcErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(tcErr)

	sectionCredential := testCredential.Section("")
	assert.Equal(global.SECTION_NAME_BLANK, sectionCredential.selectedSection)
}

func TestProfile(t *testing.T) {
	assert, _ := global.InitTest(t)

	testFactory, _ := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	_, credentialErr := NewProfile("", testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.EqualError(credentialErr, profile.ERR_PROFILE_NAME_MUST_MATCH_REGEX)

	testCredentials, credentialErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(credentialErr)
	attrErr := testCredentials.Section(global.TEST_VAR_FIRST_SECTION_KEY).SetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	assert.NoError(attrErr)

	attrValue := testCredentials.Section("").GetAttribute("")
	assert.Equal("", attrValue)
	attrValue = testCredentials.Section(global.TEST_VAR_FIRST_SECTION_KEY).GetAttribute("")
	assert.Equal("", attrValue)
	attrValue = testCredentials.Section(global.TEST_VAR_FIRST_SECTION_KEY).GetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL)
	assert.Equal(global.TEST_VAR_ATTRIBUTE_VALUE, attrValue)
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
