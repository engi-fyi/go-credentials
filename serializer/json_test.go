package serializer

import (
	"github.com/engi-fyi/go-credentials/factory"
	"github.com/engi-fyi/go-credentials/global"
	"os"

	//"os"
	"testing"
)

func TestToJson(t *testing.T) {
	assert, _ := global.InitTest(t)
	testFactory, testSerializer, serializeErr := createTestJson(global.DEFAULT_PROFILE_NAME, false)

	assert.NoError(serializeErr)
	assert.FileExists(testFactory.CredentialFile)
	assert.FileExists(testSerializer.ConfigFile)

	testFactoryTwo, testSerializerTwo, serializeErr := createTestJson(global.TEST_VAR_FIRST_PROFILE_LABEL, true)
	assert.NoError(serializeErr)
	assert.FileExists(testFactory.CredentialFile)
	assert.FileExists(testSerializer.ConfigFile)

	assert.NoError(serializeErr)
	assert.FileExists(testFactoryTwo.CredentialFile)
	assert.FileExists(testSerializerTwo.ConfigFile)

	os.RemoveAll(testFactory.ParentDirectory)
}

func TestFromJson(t *testing.T) {
	assert, _ := global.InitTest(t)
	testFactory, testSerializer, serializeErr := createTestJson(global.DEFAULT_PROFILE_NAME, false)

	assert.NoError(serializeErr)
	assert.FileExists(testFactory.CredentialFile)
	assert.FileExists(testSerializer.ConfigFile)

	username, password, attributes, deserializeErr := testSerializer.Deserialize()
	assert.Equal(global.TEST_VAR_USERNAME, username)
	assert.Equal(global.TEST_VAR_PASSWORD, password)
	assert.Equal(global.TEST_VAR_ATTRIBUTE_VALUE, attributes[global.TEST_VAR_FIRST_SECTION_KEY][global.TEST_VAR_ATTRIBUTE_NAME_LABEL])
	assert.NoError(deserializeErr)

	testFactoryTwo, testSerializerTwo, serializeErr := createTestJson(global.TEST_VAR_FIRST_PROFILE_LABEL, true)

	assert.NoError(serializeErr)
	assert.FileExists(testFactoryTwo.CredentialFile)
	assert.FileExists(testSerializerTwo.ConfigFile)

	username, password, attributes, deserializeErr = testSerializerTwo.Deserialize()
	assert.Equal(global.TEST_VAR_USERNAME_ALTERNATE, username)
	assert.Equal(global.TEST_VAR_PASSWORD_ALTERNATE, password)
	assert.Equal(global.TEST_VAR_ATTRIBUTE_VALUE, attributes[global.TEST_VAR_FIRST_SECTION_KEY][global.TEST_VAR_ATTRIBUTE_NAME_LABEL])
	assert.NoError(deserializeErr)

	os.RemoveAll(testFactory.ParentDirectory)
}

func createTestJson(profileName string, useAlternates bool) (*factory.Factory, *Serializer, error) {
	testFactory, _ := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	testFactory.SetOutputType(global.OUTPUT_TYPE_JSON)
	testSerializer := New(testFactory, profileName)
	attributes := map[string]map[string]string{
		global.TEST_VAR_FIRST_SECTION_KEY: {
			global.TEST_VAR_ATTRIBUTE_NAME_LABEL: global.TEST_VAR_ATTRIBUTE_VALUE,
		},
	}

	if useAlternates {
		return testFactory, testSerializer, testSerializer.Serialize(global.TEST_VAR_USERNAME_ALTERNATE, global.TEST_VAR_PASSWORD_ALTERNATE, attributes)
	}

	return testFactory, testSerializer, testSerializer.Serialize(global.TEST_VAR_USERNAME, global.TEST_VAR_PASSWORD, attributes)
}
