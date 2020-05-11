package serializer

import (
	"github.com/engi-fyi/go-credentials/factory"
	"github.com/engi-fyi/go-credentials/global"
	"os"
	"testing"
)

func TestToIni(t *testing.T) {
	assert := global.InitTest(t)
	testFactory, testSerializer, serializeErr := createTestIni(global.DEFAULT_PROFILE_NAME, false)
	assert.NoError(serializeErr)

	testFactory, testSerializerTwo, serializeErr := createTestIni(global.TEST_VAR_FIRST_PROFILE_LABEL, true)
	assert.NoError(serializeErr)

	assert.FileExists(testSerializer.CredentialFile)
	assert.FileExists(testSerializer.ConfigFile)
	assert.FileExists(testSerializerTwo.ConfigFile)

	os.RemoveAll(testFactory.ParentDirectory)
}

func TestFromIni(t *testing.T) {
	assert := global.InitTest(t)
	testFactory, _, serializeErr := createTestIni(global.DEFAULT_PROFILE_NAME, false)
	assert.NoError(serializeErr)

	testSerializerTwo := New(testFactory, global.DEFAULT_PROFILE_NAME)
	username, password, attributes, serializeErr := testSerializerTwo.Deserialize()
	assert.NoError(serializeErr)
	assert.Equal(global.TEST_VAR_USERNAME, username)
	assert.Equal(global.TEST_VAR_PASSWORD, password)
	assert.Equal(global.TEST_VAR_ATTRIBUTE_VALUE, attributes[global.TEST_VAR_FIRST_SECTION_KEY][global.TEST_VAR_ATTRIBUTE_NAME_LABEL])

	os.RemoveAll(testFactory.ParentDirectory)
}

func createTestIni(profileName string, useAlternates bool) (*factory.Factory, *Serializer, error) {
	testFactory, _ := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
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
