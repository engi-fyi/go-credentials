package serializer

import (
	"encoding/json"
	"github.com/engi-fyi/go-credentials/global"
	"io/ioutil"
	"os"
)

/*
ToJson is responsible for serializing a Credential and Profile to an json file. Attribute sections are directly
translatable to parent keys in the config file. Username and Password will have an appropriate label (either the default
or an alternate set in the Credential's related Factory.
*/
func (thisSerializer *Serializer) ToJson(username string, password string, attributes map[string]map[string]string) error {
	thisSerializer.Factory.Log.Info().Msg("Serializing credential and profile to json file.")
	credentialErr := thisSerializer.saveCredentialJson(username, password)

	if credentialErr != nil {
		return credentialErr
	}

	profileErr := thisSerializer.saveProfileJson(attributes)

	if profileErr != nil {
		return profileErr
	}

	return nil
}

func (thisSerializer *Serializer) saveCredentialJson(username string, password string) error {
	thisSerializer.Factory.Log.Trace().Msg("Serializing credential to json file.")
	existingCredential, initErr := initJsonCredential(thisSerializer.CredentialFile)

	if initErr != nil {
		return initErr
	}

	existingCredential.Credentials[thisSerializer.ProfileName] = serializedCredentials{
		Username: username,
		Password: password,
	}

	outJson, marshalErr := json.MarshalIndent(existingCredential, "", global.INDENT_JSON)

	if marshalErr != nil {
		return marshalErr
	}

	writeErr := ioutil.WriteFile(thisSerializer.CredentialFile, outJson, 0600)

	if writeErr != nil {
		return writeErr
	}

	thisSerializer.Factory.Log.Info().Msg("Credential json file saved successfully.")
	return nil
}

func (thisSerializer *Serializer) saveProfileJson(attributes map[string]map[string]string) error {
	thisSerializer.Factory.Log.Trace().Msg("Serializing profile to json file.")
	existingProfile, initErr := initJsonProfile(thisSerializer.CredentialFile)

	if initErr != nil {
		return initErr
	}

	existingProfile.Attributes = attributes
	outJson, marshalErr := json.MarshalIndent(existingProfile, "", global.INDENT_JSON)

	if marshalErr != nil {
		return marshalErr
	}

	writeErr := ioutil.WriteFile(thisSerializer.ConfigFile, outJson, 0600)

	if writeErr != nil {
		return writeErr
	}

	thisSerializer.Factory.Log.Info().Msg("Profile json file saved successfully.")
	return nil
}

func initJson(fileName string) ([]byte, error) {
	if _, statErr := os.Stat(fileName); os.IsNotExist(statErr) {
		emptyFile, emptyErr := os.Create(fileName)

		if emptyErr != nil {
			return []byte{}, emptyErr
		}

		_, writeErr := emptyFile.WriteString("{}")

		if writeErr != nil {
			return []byte{}, writeErr
		}

		closeErr := emptyFile.Close()

		if closeErr != nil {
			return []byte{}, closeErr
		}
	}

	//#nosec
	inJson, readErr := ioutil.ReadFile(fileName)

	if readErr != nil {
		return []byte{}, readErr
	}

	return inJson, nil
}

func initJsonCredential(filename string) (*credentialSerializer, error) {
	credentialContents, initErr := initJson(filename)

	if initErr != nil {
		return nil, initErr
	}

	var existingCredential credentialSerializer
	unmarshalErr := json.Unmarshal(credentialContents, &existingCredential)

	if existingCredential.Credentials == nil {
		existingCredential.Credentials = make(map[string]serializedCredentials)
	}

	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return &existingCredential, nil
}

func initJsonProfile(filename string) (*profileSerializer, error) {
	profileContents, initErr := initJson(filename)

	if initErr != nil {
		return nil, initErr
	}

	var existingProfile profileSerializer
	unmarshalErr := json.Unmarshal(profileContents, &existingProfile)

	if existingProfile.Attributes == nil {
		existingProfile.Attributes = make(map[string]map[string]string)
	}

	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return &existingProfile, nil
}

/*
FromJson is responsible for deserializing a Credential and Profile from an json file. Attribute sections are directly
translatable to parent keys in the config file. Alternate field labels are restored from the json, so the same Factory
object will need to be used when deserializing.
*/
func (thisSerializer *Serializer) FromJson() (string, string, map[string]map[string]string, error) {
	thisSerializer.Factory.Log.Info().Msg("Deserializing credentials and profile from json file.")
	username, password, credentialErr := thisSerializer.loadCredentialJson()

	if credentialErr != nil {
		return "", "", make(map[string]map[string]string), credentialErr
	}

	attributes, attributeErr := thisSerializer.loadProfileJson()

	if attributeErr != nil {
		return "", "", make(map[string]map[string]string), attributeErr
	}

	return username, password, attributes, nil
}

func (thisSerializer *Serializer) loadCredentialJson() (string, string, error) {
	existingCredential, initErr := initJsonCredential(thisSerializer.CredentialFile)

	if initErr != nil {
		return "", "", initErr
	}

	return existingCredential.Credentials[thisSerializer.ProfileName].Username,
		existingCredential.Credentials[thisSerializer.ProfileName].Password,
		nil
}

func (thisSerializer *Serializer) loadProfileJson() (map[string]map[string]string, error) {
	existingProfile, initErr := initJsonProfile(thisSerializer.ConfigFile)

	if initErr != nil {
		return make(map[string]map[string]string), initErr
	}

	return existingProfile.Attributes, nil
}
