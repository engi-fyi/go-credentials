package serializer

import (
	"github.com/rs/zerolog/log"
	"gopkg.in/ini.v1"
	"os"
)

/*
ToIni is responsible for serializing a Credential and Profile to an ini file. Attribute sections are directly
translatable to sections in the Profile. Username and Password will have an appropriate label (either the default or an
alternate set in the Credential's related Factory.
 */
func (thisSerializer *Serializer) ToIni(username string, password string, attributes map[string]map[string]string) error {
	log.Trace().Msg("Serializing credential and profile to ini file.")
	credentialErr := thisSerializer.saveCredentialIni(username, password)

	if credentialErr != nil {
		return credentialErr
	}

	profileErr := thisSerializer.saveProfileIni(attributes)

	if profileErr != nil {
		return profileErr
	}

	return nil
}

func (thisSerializer *Serializer) saveCredentialIni(username string, password string) error {
	log.Trace().Msg("Serializing credential to ini file.")
	credentialIni, credIniError := initIni(thisSerializer.CredentialFile)
	log.Trace().Msg("Getting alternate username and password labels.")
	usernameKey := thisSerializer.Factory.GetAlternateUsername()
	passwordKey := thisSerializer.Factory.GetAlternatePassword()

	if credIniError != nil {
		return credIniError
	}

	log.Trace().Msg("Adding username and password to credentials file.")
	credentialIni.Section(thisSerializer.ProfileName).Key(usernameKey).SetValue(username)
	credentialIni.Section(thisSerializer.ProfileName).Key(passwordKey).SetValue(password)

	log.Trace().Msg("Saving credential ini file.")
	saveErr := credentialIni.SaveTo(thisSerializer.CredentialFile)

	if saveErr != nil {
		log.Error().Str("file", thisSerializer.CredentialFile).Err(saveErr).Msg("Error saving ini file.")
		return saveErr
	}

	log.Trace().Msg("Credential ini file saved successfully.")
	return nil
}

func (thisSerializer *Serializer) saveProfileIni(attributes map[string]map[string]string) error {
	log.Trace().Msg("Serializing profile to ini file.")
	profileIni := ini.Empty()

	log.Trace().Msg("Processing attributes.")
	for key, value := range attributes {
		log.Trace().Str("attribute", key).Msg("Adding section.")
		mySection, sectionErr := profileIni.NewSection(key)

		if sectionErr != nil {
			return sectionErr
		}

		for subKey, subValue := range value {
			log.Trace().Str("attribute", subKey).Msg("Adding attribute.")
			_, keyErr := mySection.NewKey(subKey, subValue)

			if keyErr != nil {
				return keyErr
			}
		}
	}

	log.Trace().Msg("Saving profile ini file.")
	saveErr := profileIni.SaveTo(thisSerializer.ConfigFile)

	if saveErr != nil {
		return saveErr
	}

	log.Trace().Msg("Profile ini file saved successfully.")
	return nil
}

func initIni(fileName string) (*ini.File, error) {
	if _, statErr := os.Stat(fileName); os.IsNotExist(statErr) {
		emptyFile, emptyErr := os.Create(fileName)

		if emptyErr != nil {
			log.Error().Str("file", fileName).Err(emptyErr).Msg("Error creating file.")
			return nil, emptyErr
		}

		closeErr := emptyFile.Close()

		if closeErr != nil {
			return nil, closeErr
		}
	}

	return ini.Load(fileName)
}

/*
FromIni is responsible for deserializing a Credential and Profile from an ini file. Attribute sections are directly
translatable to sections in the Profile. Alternate field labels are restored from the ini, so the same Factory
object will need to be used when deserializing.
*/
func (thisSerializer *Serializer) FromIni() (string, string, map[string]map[string]string, error)  {
	username, password, credentialErr := thisSerializer.loadCredentialIni()

	if credentialErr != nil {
		return "", "", make(map[string]map[string]string), credentialErr
	}

	attributes, attributeErr := thisSerializer.loadProfileIni()

	if attributeErr != nil {
		return "", "", make(map[string]map[string]string), attributeErr
	}

	return username, password, attributes, nil
}

func (thisSerializer *Serializer) loadCredentialIni() (string, string, error) {
	credentialIni, initErr := initIni(thisSerializer.CredentialFile)
	usernameKey := thisSerializer.Factory.GetAlternateUsername()
	passwordKey := thisSerializer.Factory.GetAlternatePassword()

	if initErr != nil {
		return "", "", initErr
	}

	username := credentialIni.Section(thisSerializer.ProfileName).Key(usernameKey).String()
	password := credentialIni.Section(thisSerializer.ProfileName).Key(passwordKey).String()

	return username, password, nil
}

func (thisSerializer *Serializer) loadProfileIni() (map[string]map[string]string, error) {
	profileIni, initErr := initIni(thisSerializer.ConfigFile)
	myAttributes := make(map[string]map[string]string)

	if initErr != nil {
		return make(map[string]map[string]string), initErr
	}

	for _, section := range profileIni.Sections() {
		keys := section.KeysHash()
		myAttributes[section.Name()] = keys
	}

	return myAttributes, nil
}