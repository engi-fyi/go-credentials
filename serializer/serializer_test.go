package serializer

import (
	"github.com/engi-fyi/go-credentials/factory"
	"github.com/engi-fyi/go-credentials/global"
	"testing"
)

func TestNew(t *testing.T) {
	assert, _ := global.InitTest(t)
	testFactory, _ := factory.New(global.TEST_VAR_APPLICATION_NAME)
	testSerializer := New(testFactory, global.DEFAULT_PROFILE_NAME)
	assert.Equal(testFactory.CredentialFile, testSerializer.CredentialFile)
	assert.Equal(testFactory.ConfigDirectory+global.DEFAULT_PROFILE_NAME, testSerializer.ConfigFile)
	assert.Equal(global.DEFAULT_PROFILE_NAME, testSerializer.ProfileName)
	assert.True(testSerializer.Initialized)
	assert.Equal(testFactory, testSerializer.Factory)
}

func TestSerialize(t *testing.T) {
	assert, _ := global.InitTest(t)
	testFactory, _ := factory.New(global.TEST_VAR_APPLICATION_NAME)
	testFactory.OutputType = global.OUTPUT_TYPE_INVALID
	testSerializer := New(testFactory, global.DEFAULT_PROFILE_NAME)
	serializeErr := testSerializer.Serialize("", "", make(map[string]map[string]string))
	assert.EqualError(serializeErr, ERR_UNRECOGNIZED_OUTPUT_TYPE)
}

func TestDeserialize(t *testing.T) {
	assert, _ := global.InitTest(t)
	testFactory, _ := factory.New(global.TEST_VAR_APPLICATION_NAME)
	testFactory.OutputType = global.OUTPUT_TYPE_INVALID
	testSerializer := New(testFactory, global.DEFAULT_PROFILE_NAME)
	_, _, _, deserializeErr := testSerializer.Deserialize()
	assert.EqualError(deserializeErr, ERR_UNRECOGNIZED_OUTPUT_TYPE)
}
