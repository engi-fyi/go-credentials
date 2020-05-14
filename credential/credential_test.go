package credential

import (
	"github.com/engi-fyi/go-credentials/factory"
	"github.com/engi-fyi/go-credentials/global"
	"github.com/engi-fyi/go-credentials/profile"
	"github.com/engi-fyi/go-credentials/serializer"
	"github.com/rs/zerolog"
	as "github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCredentialNewBadFactory(t *testing.T) {
	assert, log, _ := initTest(t)
	log.Info().Msg("Testing to ensure an initialized factory is required.")
	blankFactory := factory.Factory{}
	_, factoryErr := New(&blankFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.EqualError(factoryErr, ERR_FACTORY_MUST_BE_INITIALIZED)
	_, factoryErr = New(nil, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.EqualError(factoryErr, ERR_FACTORY_MUST_BE_INITIALIZED)
}

func TestCredentialNewBlankUsernameAndPassword(t *testing.T) {
	assert, log, testFactory := initTest(t)
	log.Info().Msg("Testing to ensure blank username or password can't be used.")
	_, missingErr := New(testFactory, "", global.TEST_VAR_PASSWORD)
	assert.EqualError(missingErr, ERR_USERNAME_OR_PASSWORD_NOT_SET)
	_, missingErr = New(testFactory, global.TEST_VAR_USERNAME, "")
	assert.EqualError(missingErr, ERR_USERNAME_OR_PASSWORD_NOT_SET)
	_, missingErr = New(testFactory, "", "")
	assert.EqualError(missingErr, ERR_USERNAME_OR_PASSWORD_NOT_SET)
}

func TestCredentialNew(t *testing.T) {
	assert, log, testFactory := initTest(t)
	log.Info().Msg("Testing basic creation.")
	testCredential, newErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(newErr)
	assert.Equal(global.TEST_VAR_USERNAME, testCredential.Username)
	assert.Equal(global.TEST_VAR_PASSWORD, testCredential.Password)
	assert.True(testCredential.Factory.Initialized)
	assert.True(testCredential.Initialized)
}

func TestCredentialAttributeGetUsernameRedirect(t *testing.T) {
	assert, log, testFactory := initTest(t)
	log.Info().Msg("Testing username accessibility through GetAttribute.")
	testCredential, newErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(newErr)
	username := testCredential.GetAttribute("username")
	assert.Equal(global.TEST_VAR_USERNAME, username)
}

func TestCredentialAttributeSetRedirectUsername(t *testing.T) {
	assert, log, testFactory := initTest(t)
	log.Info().Msg("Testing setting username through attribute redirection.")
	testCredential, newErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(newErr)
	attrErr := testCredential.SetAttribute("username", global.TEST_VAR_USERNAME_ALTERNATE)
	assert.NoError(attrErr)
	assert.Equal(testCredential.Username, global.TEST_VAR_USERNAME_ALTERNATE)
}

func TestCredentialAttributeGetPasswordRedirect(t *testing.T) {
	assert, log, testFactory := initTest(t)
	log.Info().Msg("Testing getting password through attribute redirection.")
	testCredential, newErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(newErr)
	password := testCredential.GetAttribute("password")
	assert.Equal(global.TEST_VAR_PASSWORD, password)
}

func TestCredentialAttributeSetRedirectPassword(t *testing.T) {
	assert, log, testFactory := initTest(t)
	log.Info().Msg("Testing setting password through attribute redirection.")
	testCredential, newErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(newErr)
	attrErr := testCredential.SetAttribute("password", global.TEST_VAR_PASSWORD_ALTERNATE)
	assert.NoError(attrErr)
	assert.Equal(testCredential.Password, global.TEST_VAR_PASSWORD_ALTERNATE)
}

func TestCredentialAttributeGetAttributeNotSet(t *testing.T) {
	assert, log, testFactory := initTest(t)
	log.Info().Msg("Testing getting an attribute when it hasn't been set.")
	testCredential, newErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(newErr)
	log.Info().Msg("Testing that correct error returned when an attribute not set and GetAttribute is used.")
	notSet := testCredential.GetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL)
	assert.Equal("", notSet)
}

func TestCredentialAttributeSet(t *testing.T) {
	assert, log, testFactory := initTest(t)
	log.Info().Msg("Test attribute setting.")
	testCredential, newErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(newErr)
	attrErr := testCredential.SetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	assert.NoError(attrErr)
	set := testCredential.GetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL)
	assert.Equal(global.TEST_VAR_ATTRIBUTE_VALUE, set)
}

func TestCredentialAttributeBadSet(t *testing.T) {
	assert, log, testFactory := initTest(t)
	log.Info().Msg("Testing attribute setting when the key does not match standards")
	testCredential, newErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(newErr)
	badSetErr := testCredential.SetAttribute(global.TEST_VAR_BAD_ATTRIBUTE_NAME, global.TEST_VAR_ATTRIBUTE_VALUE)
	assert.EqualError(badSetErr, ERR_KEY_MUST_MATCH_REGEX)
	badSetGet := testCredential.GetAttribute(global.TEST_VAR_BAD_ATTRIBUTE_NAME)
	assert.Equal("", badSetGet)
}

func TestCredentialSave(t *testing.T) {
	assert, log, testFactory := initTest(t)
	log.Info().Msg("Testing the saving of a credential.")

	testCredential, newErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(newErr)
	attrErr := testCredential.SetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	assert.NoError(attrErr)
	saveErr := testCredential.Save()
	assert.NoError(saveErr)
}

func TestCredentialLoadFromFile(t *testing.T) {
	assert, log, testFactory := initTest(t)
	log.Info().Msg("Testing load when no file present.")
	testCredential, loadErr := Load(testFactory)
	assert.NoError(loadErr)
	assert.Equal(global.TEST_VAR_USERNAME, testCredential.Username)
	assert.Equal(global.TEST_VAR_PASSWORD, testCredential.Password)
	assert.Equal(global.TEST_VAR_ATTRIBUTE_VALUE, testCredential.GetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL))

	parentDirectoryCleanup(t)
}

func TestCredentialNotInitializedSave(t *testing.T) {
	assert, log, testFactory := initTest(t)
	log.Info().Msg("Testing saving a credential that has not been initialized.")
	testCredential := Credential{}
	testCredential.Factory = testFactory
	saveErr := testCredential.Save()
	assert.EqualError(saveErr, ERR_NOT_INITIALIZED)
}

func TestCredentialSaveProfileNotInitialized(t *testing.T) {
	assert, log, testFactory := initTest(t)
	log.Info().Msg("Testing that files save correctly.")
	testCredential, newErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(newErr)
	testCredential.Factory = &factory.Factory{}
	saveErr := testCredential.Save()
	assert.Error(saveErr, factory.ERR_FACTORY_NOT_INITIALIZED)
	parentDirectoryCleanup(t)
}

func TestCredentialSaveFactoryInconsistent(t *testing.T) {
	assert, log := global.InitTest(t)
	log.Info().Msg("Testing a factory that has been messed and had it's output changed.")
	testCredential, newErr := buildTestCredentials()
	assert.NoError(newErr)
	testCredential.Factory.OutputType = global.OUTPUT_TYPE_INVALID
	saveErr := testCredential.Save()
	log.Info().Msg("Testing that files save correctly.")
	assert.Error(saveErr, factory.ERR_FACTORY_NOT_INITIALIZED)
	parentDirectoryCleanup(t)
}

func TestCredentialProfileNotInitialized(t *testing.T) {
	assert, log, testFactory := initTest(t)
	log.Info().Msg("Testing a profile that has not been initialized.")
	testCredential, newErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(newErr)
	testCredential.Profile = &profile.Profile{}
	saveErr := testCredential.Save()
	assert.Error(saveErr, profile.ERR_PROFILE_NOT_INITIALIZED)
	parentDirectoryCleanup(t)
}

func TestCredentialLoadFromFileNoFile(t *testing.T) {
	assert, log := global.InitTest(t)
	log.Info().Msg("Testing load when no file present.")

	for _, fileType := range serializer.GetSupportedFileTypes() {
		log.Info().Msgf("Testing the '%v' file type.", fileType)
		secondTestFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME)
		assert.NoError(factoryErr)
		outputErr := secondTestFactory.SetOutputType(fileType)
		assert.NoError(outputErr)
		_, loadErr := LoadFromProfile(global.DEFAULT_PROFILE_NAME, secondTestFactory)
		assert.Error(loadErr)
	}
}

func TestCredentialLoadFactoryInvalidOutputType(t *testing.T) {
	assert, log := global.InitTest(t)
	log.Info().Msg("Testing an invalid output type against a factory object.")
	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME)
	assert.NoError(factoryErr)
	testFactory.OutputType = global.OUTPUT_TYPE_INVALID
	_, loadErr := Load(testFactory)
	assert.EqualError(loadErr, serializer.ERR_UNRECOGNIZED_OUTPUT_TYPE)
}

func TestCredentialLoadProfileNotInitialized(t *testing.T) {
	assert, log := global.InitTest(t)
	log.Info().Msg("Testing the creation of a credential with a uninitialized factory.")
	secondTestFactory := &factory.Factory{}
	_, loadErr := Load(secondTestFactory)
	assert.EqualError(loadErr, ERR_FACTORY_MUST_BE_INITIALIZED)
}

func TestCredentialLoadEnvNoValues(t *testing.T) {
	assert, log := global.InitTest(t)
	log.Info().Msg("Testing failing to serialize variables from the environment due to not existing.")
	for _, label := range []string{global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL, global.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL, global.TEST_VAR_ENVIRONMENT_ATTRIBUTE_NAME_LABEL} {
		usErr := os.Unsetenv(label)
		assert.NoError(usErr)
	}

	secondTestFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME)
	assert.NoError(factoryErr)
	soErr := secondTestFactory.SetOutputType(global.OUTPUT_TYPE_ENV)
	assert.NoError(soErr)
	_, credErr := Load(secondTestFactory)
	assert.EqualError(credErr, serializer.ERR_REQUIRED_VARIABLE_USERNAME_NOT_FOUND)

	for _, label := range []string{global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL, global.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL, global.TEST_VAR_ENVIRONMENT_ATTRIBUTE_NAME_LABEL} {
		_, exists := os.LookupEnv(label)
		assert.False(exists)
	}
}

func TestCredentialSaveEnv(t *testing.T) {
	assert, log := global.InitTest(t)
	log.Info().Msg("Testing saving variables to environment variables.")
	testCredentials, tcErr := buildTestCredentials()
	testCredentials.Factory.SetOutputType(global.OUTPUT_TYPE_ENV)
	assert.NoError(tcErr)

	soErr := testCredentials.Section(global.TEST_VAR_FIRST_SECTION_KEY).SetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	assert.NoError(soErr)
	deployErr := testCredentials.Save()
	assert.NoError(deployErr)

	for _, label := range []string{global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL, global.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL, global.TEST_VAR_ENVIRONMENT_ATTRIBUTE_NAME_LABEL} {
		_, exists := os.LookupEnv(label)
		assert.True(exists)
	}
}

func TestCredentialLoadEnv(t *testing.T) {
	assert, log, testFactory := initTest(t)
	log.Info().Msg("Testing loading credentials from the environment.")

	soErr := testFactory.SetOutputType(global.OUTPUT_TYPE_ENV)
	assert.NoError(soErr)
	testCredential, credErr := Load(testFactory)
	assert.NoError(credErr)

	assert.Equal(global.TEST_VAR_USERNAME, testCredential.Username)
	assert.Equal(global.TEST_VAR_PASSWORD, testCredential.Password)
	assert.Equal(global.TEST_VAR_ATTRIBUTE_VALUE, testCredential.Section(global.TEST_VAR_FIRST_SECTION_KEY).GetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL))

	for _, label := range []string{global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL, global.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL, global.TEST_VAR_ENVIRONMENT_ATTRIBUTE_NAME_LABEL} {
		usErr := os.Unsetenv(label)
		assert.NoError(usErr)
	}
}

func TestCredentialCreateNoProfile(t *testing.T) {
	assert, log, testFactory := initTest(t)
	log.Info().Msg("Testing the creation of a credential using the default profile.")
	testCredential, tcErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(tcErr)
	setErr := testCredential.SetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	assert.NoError(setErr)
	saveErr := testCredential.Save()
	assert.NoError(saveErr)
}

func TestCredentialCreateFirstProfile(t *testing.T) {
	assert, log, testFactory := initTest(t)
	log.Info().Msg("Testing the creation of a credential adding a profile.")
	testCredential, tcErr := NewProfile(global.TEST_VAR_FIRST_PROFILE_LABEL, testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(tcErr)
	setErr := testCredential.SetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	assert.NoError(setErr)
	saveErr := testCredential.Save()
	assert.NoError(saveErr)
}

func TestCredentialCreateSecondProfile(t *testing.T) {
	assert, log, testFactory := initTest(t)
	log.Info().Msg("Testing the creation of a credential adding a second profile.")
	testCredential, tcErr := NewProfile(global.TEST_VAR_SECOND_PROFILE_LABEL, testFactory, global.TEST_VAR_USERNAME_ALTERNATE, global.TEST_VAR_PASSWORD_ALTERNATE)
	assert.NoError(tcErr)
	setErr := testCredential.SetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	assert.NoError(setErr)
	saveErr := testCredential.Save()
	assert.NoError(saveErr)
}

func TestCredentialLoadNoProfile(t *testing.T) {
	assert, log, testFactory := initTest(t)
	log.Info().Msg("Testing the loading of an existing credential using the default profile.")
	testCredential, tcErr := Load(testFactory)
	assert.NoError(tcErr)
	assert.Equal(global.DEFAULT_PROFILE_NAME, testCredential.Profile.Name)
	attrValue := testCredential.GetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL)
	assert.Equal(global.TEST_VAR_ATTRIBUTE_VALUE, attrValue)
	assert.True(testCredential.Initialized)
	assert.True(testCredential.Profile.Initialized)
}

func TestCredentialLoadFirstProfile(t *testing.T) {
	assert, log, testFactory := initTest(t)
	log.Info().Msg("Testing the loading of an existing credential using the first profile.")
	testCredential, tcErr := LoadFromProfile(global.TEST_VAR_FIRST_PROFILE_LABEL, testFactory)
	assert.NoError(tcErr)
	assert.Equal(global.TEST_VAR_FIRST_PROFILE_LABEL, testCredential.Profile.Name)
	attrValue := testCredential.GetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL)
	assert.Equal(global.TEST_VAR_ATTRIBUTE_VALUE, attrValue)
	assert.Equal(global.TEST_VAR_USERNAME, testCredential.Username)
	assert.Equal(global.TEST_VAR_PASSWORD, testCredential.Password)
	assert.True(testCredential.Initialized)
	assert.True(testCredential.Profile.Initialized)
}

func TestCredentialLoadSecondProfile(t *testing.T) {
	assert, log, testFactory := initTest(t)
	log.Info().Msg("Testing the loading of an existing credential using the second profile.")
	testCredential, tcErr := LoadFromProfile(global.TEST_VAR_SECOND_PROFILE_LABEL, testFactory)
	assert.NoError(tcErr)
	assert.Equal(global.TEST_VAR_SECOND_PROFILE_LABEL, testCredential.Profile.Name)
	attrValue := testCredential.GetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL)
	assert.Equal(global.TEST_VAR_ATTRIBUTE_VALUE, attrValue)
	assert.Equal(global.TEST_VAR_USERNAME_ALTERNATE, testCredential.Username)
	assert.Equal(global.TEST_VAR_PASSWORD_ALTERNATE, testCredential.Password)
	assert.True(testCredential.Initialized)
	assert.True(testCredential.Profile.Initialized)
	parentDirectoryCleanup(t)
}

func TestCredentialDeleteAttributes(t *testing.T) {
	assert, log, testFactory := initTest(t)
	log.Info().Msg("Testing the loading of an existing credential using the second profile.")
	testCredential, tcErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(tcErr)
	setErr := testCredential.SetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	assert.NoError(setErr)
	assert.Equal(global.TEST_VAR_ATTRIBUTE_VALUE, testCredential.GetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL))
	deleteErr := testCredential.DeleteAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL)
	assert.NoError(deleteErr)
	assert.Equal("", testCredential.GetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL))
	deleteErr = testCredential.DeleteAttribute(global.TEST_VAR_USERNAME_LABEL)
	assert.EqualError(deleteErr, ERR_CANNOT_REMOVE_USERNAME)
	deleteErr = testCredential.DeleteAttribute(global.TEST_VAR_PASSWORD_LABEL)
	assert.EqualError(deleteErr, ERR_CANNOT_REMOVE_PASSWORD)
	parentDirectoryCleanup(t)
}

func TestCredentialCannotRedirectUsernameSetAttribute(t *testing.T) {
	assert, log, testFactory := initTest(t)
	log.Info().Msg("Testing attempting to set username when a section has been set.")
	testCredential, tcErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(tcErr)
	setErr := testCredential.Section(global.TEST_VAR_FIRST_SECTION_KEY).SetAttribute(global.TEST_VAR_USERNAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	assert.EqualError(setErr, ERR_CANNOT_SET_USERNAME_WHEN_USING_SECTION)
}

func TestCredentialCannotRedirectPasswordSetAttribute(t *testing.T) {
	assert, log, testFactory := initTest(t)
	log.Info().Msg("Testing attempting to set password when a section has been set.")
	testCredential, tcErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(tcErr)
	setErr := testCredential.Section(global.TEST_VAR_FIRST_SECTION_KEY).SetAttribute(global.TEST_VAR_PASSWORD_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	assert.EqualError(setErr, ERR_CANNOT_SET_PASSWORD_WHEN_USING_SECTION)
}

func TestCredentialDeleteSectionAttributes(t *testing.T) {
	assert, log, testFactory := initTest(t)
	log.Info().Msg("Testing the loading of an existing credential using the second profile.")
	testCredential, tcErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(tcErr)
	setErr := testCredential.Section(global.TEST_VAR_FIRST_SECTION_KEY).SetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL, global.TEST_VAR_ATTRIBUTE_VALUE)
	assert.NoError(setErr)
	assert.Equal(global.TEST_VAR_ATTRIBUTE_VALUE, testCredential.Section(global.TEST_VAR_FIRST_SECTION_KEY).GetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL))
	deleteErr := testCredential.Section(global.TEST_VAR_FIRST_SECTION_KEY).DeleteAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL)
	assert.NoError(deleteErr)
	assert.Equal("", testCredential.Section(global.TEST_VAR_FIRST_SECTION_KEY).GetAttribute(global.TEST_VAR_ATTRIBUTE_NAME_LABEL))
	parentDirectoryCleanup(t)
}

/*
This tests:
 - Whether a value is set correctly in selectedSection of the returned Credential.
 - Whether a cloned section is correctly a shallow copy of the initial Credential.
 */
func TestSectionCredentialWithValue(t *testing.T) {
	assert, log := global.InitTest(t)
	log.Info().Msg("Testing to ensure Section() returns a value as expected")
	testFactory, _ := factory.New(global.TEST_VAR_APPLICATION_NAME)
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
	assert, log := global.InitTest(t)
	log.Info().Msg("Testing to ensure Section() returns a the blank section value as expected")
	testFactory, _ := factory.New(global.TEST_VAR_APPLICATION_NAME)
	testCredential, tcErr := New(testFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
	assert.NoError(tcErr)

	sectionCredential := testCredential.Section("")
	assert.Equal(global.SECTION_NAME_BLANK, sectionCredential.selectedSection)
}

func TestProfile(t *testing.T) {
	assert, _ := global.InitTest(t)

	testFactory, _ := factory.New(global.TEST_VAR_APPLICATION_NAME)
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
	buildFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME)

	if factoryErr != nil {
		return nil, factoryErr
	}

	return New(buildFactory, global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD)
}

func initTest(t* testing.T) (*as.Assertions, zerolog.Logger, *factory.Factory) {
	assert, log := global.InitTest(t)
	testFactory, tfErr := factory.New(global.TEST_VAR_APPLICATION_NAME)
	assert.NoError(tfErr)
	assert.True(testFactory.Initialized)
	return assert, log, testFactory
}

func parentDirectoryCleanup(t *testing.T) {
	assert, _, testFactory := initTest(t)
	rmErr := os.RemoveAll(testFactory.ParentDirectory)
	assert.NoError(rmErr)
	assert.NoDirExists(testFactory.ParentDirectory)
}