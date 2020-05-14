package serializer

import (
	"github.com/engi-fyi/go-credentials/factory"
	"github.com/engi-fyi/go-credentials/global"
	"os"
	"testing"
)

func TestToEnv(t *testing.T) {
	assert, _ := global.InitTest(t)

	_, _, serializeErr := createTestEnv(global.DEFAULT_PROFILE_NAME, false)
	assert.NoError(serializeErr)

	value, exists := os.LookupEnv(global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL)
	assert.True(exists)
	assert.Equal(global.TEST_VAR_USERNAME, value)

	value, exists = os.LookupEnv(global.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL)
	assert.True(exists)
	assert.Equal(global.TEST_VAR_PASSWORD, value)

	value, exists = os.LookupEnv(global.TEST_VAR_ENVIRONMENT_ATTRIBUTE_NAME_LABEL)
	assert.True(exists)
	assert.Equal(global.TEST_VAR_ATTRIBUTE_VALUE, value)

	os.Unsetenv(global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL)
	os.Unsetenv(global.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL)
	os.Unsetenv(global.TEST_VAR_ENVIRONMENT_ATTRIBUTE_NAME_LABEL)
}

func TestFromEnv(t *testing.T) {
	assert, _ := global.InitTest(t)
	_, testSerializer, serializeErr := createTestEnv(global.DEFAULT_PROFILE_NAME, false)
	assert.NoError(serializeErr)

	username, password, attributes, serializeErr := testSerializer.Deserialize()
	assert.NoError(serializeErr)
	assert.Equal(global.TEST_VAR_USERNAME, username)
	assert.Equal(global.TEST_VAR_PASSWORD, password)
	assert.Equal(global.TEST_VAR_ATTRIBUTE_VALUE, attributes[global.TEST_VAR_FIRST_SECTION_KEY][global.TEST_VAR_ATTRIBUTE_NAME_LABEL])

	os.Unsetenv(global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL)
	os.Unsetenv(global.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL)
	os.Unsetenv(global.TEST_VAR_ENVIRONMENT_ATTRIBUTE_NAME_LABEL)
}

func TestParseEnvironmentVariable(t *testing.T) {
	assert, log := global.InitTest(t)
	log.Info().Msg("Testing environment parsing.")
	assert.True(true)

	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME)
	assert.NoError(factoryErr)
	testSerializer := New(testFactory, global.DEFAULT_PROFILE_NAME)

	log.Info().Msg("Testing username parsing.")
	profileName, fieldName, _, didParse := testSerializer.ParseEnvironmentVariable(global.TEST_VAR_ENVIRONMENT_USERNAME_LABEL)
	assert.Equal(global.DEFAULT_PROFILE_NAME, profileName)
	assert.Equal(global.TEST_VAR_USERNAME_LABEL, fieldName)
	assert.True(didParse)

	log.Info().Msg("Testing password parsing.")
	profileName, fieldName, _, didParse = testSerializer.ParseEnvironmentVariable(global.TEST_VAR_ENVIRONMENT_PASSWORD_LABEL)
	assert.Equal(global.DEFAULT_PROFILE_NAME, profileName)
	assert.Equal(global.TEST_VAR_PASSWORD_LABEL, fieldName)
	assert.True(didParse)

	log.Info().Msg("Testing attribute parsing.")
	profileName, fieldName, sectionName, didParse := testSerializer.ParseEnvironmentVariable(global.TEST_VAR_ENVIRONMENT_ATTRIBUTE_NAME_LABEL)
	assert.Equal(global.DEFAULT_PROFILE_NAME, profileName)
	assert.Equal(global.TEST_VAR_ATTRIBUTE_NAME_LABEL, fieldName)
	assert.Equal(global.TEST_VAR_FIRST_SECTION_KEY, sectionName)
	assert.True(didParse)

	log.Info().Msg("Testing alternate username parsing.")
	profileName, fieldName, _, didParse = testSerializer.ParseEnvironmentVariable(global.TEST_VAR_ENVIRONMENT_USERNAME_ALTERNATE_LABEL)
	assert.Equal(global.DEFAULT_PROFILE_NAME, profileName)
	assert.Equal(global.TEST_VAR_USERNAME_ALTERNATE_LABEL, fieldName)
	assert.True(didParse)

	log.Info().Msg("Testing alternate password parsing.")
	profileName, fieldName, _, didParse = testSerializer.ParseEnvironmentVariable(global.TEST_VAR_ENVIRONMENT_PASSWORD_ALTERNATE_LABEL)
	assert.Equal(global.DEFAULT_PROFILE_NAME, profileName)
	assert.Equal(global.TEST_VAR_PASSWORD_ALTERNATE_LABEL, fieldName)
	assert.True(didParse)
}

func createTestEnv(profileName string, useAlternates bool) (*factory.Factory, *Serializer, error) {
	testFactory, _ := factory.New(global.TEST_VAR_APPLICATION_NAME)
	testFactory.SetOutputType(global.OUTPUT_TYPE_ENV)
	testSerializer := New(testFactory, profileName)
	attributes := map[string]map[string]string{
		global.TEST_VAR_FIRST_SECTION_KEY: {
			global.TEST_VAR_ATTRIBUTE_NAME_LABEL: global.TEST_VAR_ATTRIBUTE_VALUE,
		},
	}

	return testFactory, testSerializer, testSerializer.Serialize(global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD, attributes)
}
