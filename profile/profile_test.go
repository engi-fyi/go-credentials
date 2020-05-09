package profile

import (
	"github.com/engi-fyi/go-credentials/factory"
	"os"
	"testing"
	"github.com/engi-fyi/go-credentials/global"
	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	assert := global.InitTest(t)

	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)

	testCreateIniFile(testFactory, assert)

	os.RemoveAll(testFactory.ParentDirectory)
}

func TestLoad(t *testing.T) {
	assert := global.InitTest(t)

	testFactory, factoryErr := factory.New(global.TEST_VAR_APPLICATION_NAME, false)
	assert.NoError(factoryErr)

	testProfile := testCreateIniFile(testFactory, assert)
	assert.Equal(TEST_VAR_NO_SECTION_UNIQUE_KEY_VALUE, testProfile.GetAttribute("", TEST_VAR_NO_SECTION_UNIQUE_KEY_LABEL))
	assert.Equal(TEST_VAR_DUPLICATE_KEY_VALUE, testProfile.GetAttribute("", TEST_VAR_DUPLICATE_KEY_LABEL))
	assert.Equal(TEST_VAR_FIRST_SECTION_UNIQUE_KEY_VALUE, testProfile.GetAttribute(TEST_VAR_FIRST_SECTION_KEY, TEST_VAR_FIRST_SECTION_UNIQUE_KEY_LABEL))
	assert.Equal(TEST_VAR_DUPLICATE_KEY_VALUE, testProfile.GetAttribute(TEST_VAR_FIRST_SECTION_KEY, TEST_VAR_DUPLICATE_KEY_LABEL))
	assert.Equal(TEST_VAR_SECOND_SECTION_UNIQUE_KEY_VALUE, testProfile.GetAttribute(TEST_VAR_SECOND_SECTION_KEY, TEST_VAR_SECOND_SECTION_UNIQUE_KEY_LABEL))
	assert.Equal(TEST_VAR_DUPLICATE_KEY_VALUE, testProfile.GetAttribute(TEST_VAR_SECOND_SECTION_KEY, TEST_VAR_DUPLICATE_KEY_LABEL))
	assert.Empty("", testProfile.GetAttribute(TEST_VAR_BAD_SECTION_KEY, ""))
	assert.Empty("", testProfile.GetAttribute("", TEST_VAR_BAD_KEY_LABEL))

	os.RemoveAll(testFactory.ParentDirectory)
}

func testCreateIniFile(testFactory *factory.Factory, assert *assert.Assertions) *Profile {
	testProfile, profileErr := New(TEST_VAR_FIRST_PROFILE_LABEL, testFactory)
	assert.NoError(profileErr)

	setErr := testProfile.SetAttribute("", TEST_VAR_NO_SECTION_UNIQUE_KEY_LABEL, TEST_VAR_NO_SECTION_UNIQUE_KEY_VALUE)
	assert.NoError(setErr)
	setErr = testProfile.SetAttribute("", TEST_VAR_DUPLICATE_KEY_LABEL, TEST_VAR_DUPLICATE_KEY_VALUE)
	assert.NoError(setErr)
	setErr = testProfile.SetAttribute(TEST_VAR_FIRST_SECTION_KEY, TEST_VAR_FIRST_SECTION_UNIQUE_KEY_LABEL, TEST_VAR_FIRST_SECTION_UNIQUE_KEY_VALUE)
	assert.NoError(setErr)
	setErr = testProfile.SetAttribute(TEST_VAR_FIRST_SECTION_KEY, TEST_VAR_DUPLICATE_KEY_LABEL, TEST_VAR_DUPLICATE_KEY_VALUE)
	assert.NoError(setErr)
	setErr = testProfile.SetAttribute(TEST_VAR_SECOND_SECTION_KEY, TEST_VAR_SECOND_SECTION_UNIQUE_KEY_LABEL, TEST_VAR_SECOND_SECTION_UNIQUE_KEY_VALUE)
	assert.NoError(setErr)
	setErr = testProfile.SetAttribute(TEST_VAR_SECOND_SECTION_KEY, TEST_VAR_DUPLICATE_KEY_LABEL, TEST_VAR_DUPLICATE_KEY_VALUE)
	assert.NoError(setErr)
	setErr = testProfile.SetAttribute(TEST_VAR_BAD_SECTION_KEY, "", "")
	assert.EqualError(setErr, ERR_MUST_MATCH_REGEX)
	setErr = testProfile.SetAttribute("", TEST_VAR_BAD_KEY_LABEL, "")
	assert.EqualError(setErr, ERR_MUST_MATCH_REGEX)

	saveErr := testProfile.Save()
	assert.NoError(saveErr)

	return testProfile
}