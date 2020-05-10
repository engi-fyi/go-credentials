package credential

import (
	"errors"
	"fmt"
	"github.com/engi-fyi/go-credentials/factory"
	"github.com/engi-fyi/go-credentials/global"
	"github.com/rs/zerolog/log"
	"gopkg.in/ini.v1"
	"os"
)

/*
Load is the default method of loading existing credentials. First it will attempt to load credentials from the
environment, as they take precedence over file-based credentials. If that is not successful, it will attempt to
load credentials from the credentials file using the appropriate format set in the sourceFactory. If the file doesn't
exist or is in an unexpected format, an error will be thrown.
*/
// TODO(7): Implement json format.
func Load(sourceFactory *factory.Factory) (*Credential, error) {
	var fileErr error
	loadedCredential, envErr := LoadFromEnvironment(sourceFactory)

	if envErr != nil || !loadedCredential.Initialized {
		if sourceFactory.OutputType == global.OUTPUT_TYPE_JSON {
			loadedCredential, fileErr = loadJson(sourceFactory)
		} else if sourceFactory.OutputType == global.OUTPUT_TYPE_INI {
			loadedCredential, fileErr = LoadFromIniFile(global.DEFAULT_PROFILE_NAME, sourceFactory)
		} else {
			return nil, errors.New(factory.ERR_FACTORY_INCONSISTENT_STATE)
		}
	}

	if envErr != nil && fileErr != nil {
		return nil, fmt.Errorf("%v. %v", envErr, fileErr)
	}

	return loadedCredential, nil
}

/*
LoadFromProfile specifically loads a Credential object from file. It will first load the relevant Credentials listed
under the profileName in the credential file, then it will load the variables from the profile config file into Profile.
 */
func LoadFromProfile(profileName string, sourceFactory *factory.Factory) (*Credential, error) {
	var fileErr error
	loadedCredential, envErr := LoadFromEnvironment(sourceFactory)

	if envErr != nil || !loadedCredential.Initialized {
		if sourceFactory.OutputType == global.OUTPUT_TYPE_JSON {
			loadedCredential, fileErr = loadJson(sourceFactory)
		} else if sourceFactory.OutputType == global.OUTPUT_TYPE_INI {
			loadedCredential, fileErr = LoadFromIniFile(profileName, sourceFactory)
		} else {
			return nil, errors.New(factory.ERR_FACTORY_INCONSISTENT_STATE)
		}
	}

	profileErr := loadedCredential.SetProfile(profileName)

	if envErr != nil && fileErr != nil {
		return nil, fmt.Errorf("%v. %v", envErr, fileErr)
	}

	if profileErr != nil {
		return nil, profileErr
	}

	return loadedCredential, nil
}

/*
LoadFromIniFile is responsible for loading a Credential object as in the ini format at ~/.app_name/credentials.
username and password are stored under the default heading, and all attributes are stored under the attributes
heading

Example Credential File
		[default]
		username=my_username
		password=my_password01

		[attributes]
		an_attribute=the value of an attribute

BUG(4): Respect alternates when loading from file.
*/
func LoadFromIniFile(profileName string, fromFactory *factory.Factory) (*Credential, error) {
	if fromFactory.Initialized {
		credentialFile, loadErr := loadCredentialIni(fromFactory)

		if loadErr != nil {
			log.Error().Err(loadErr).Msg("Error loading credentials.")

			return nil, loadErr
		}

		loadedCredential, credErr := getCredentialFromIni(profileName, fromFactory, credentialFile)

		if credErr != nil {
			return nil, credErr
		}

		profileErr := loadedCredential.SetProfile(global.DEFAULT_PROFILE_NAME)

		if profileErr != nil {
			return nil, profileErr
		}

		return loadedCredential, nil
	} else {
		return nil, errors.New(factory.ERR_FACTORY_NOT_INITIALIZED)
	}
}

func loadCredentialIni(fromFactory *factory.Factory) (*ini.File, error) {
	credentialFile, loadErr := ini.Load(fromFactory.CredentialFile)

	if _, cfsErr := os.Stat(fromFactory.CredentialFile); os.IsNotExist(cfsErr) {
		log.Trace().Msg("File doesn't exist, creating then exiting.")
		if _, cfsErr := os.Stat(fromFactory.CredentialFile); os.IsNotExist(cfsErr) {
			log.Trace().Str("Credential File", fromFactory.CredentialFile).Msg("Initializing credential file.")
			emptyFile, efErr := os.Create(fromFactory.CredentialFile)

			if efErr != nil {
				log.Error().Err(efErr).Msg("Error creating credential file.")
				return nil, efErr
			}

			closeErr := emptyFile.Close()

			if closeErr != nil {
				return nil, closeErr
			}

			credentialFile, loadErr = ini.Load(fromFactory.CredentialFile)

			if loadErr != nil {
				return nil, loadErr
			}
		}
	}

	if loadErr != nil {
		return nil, loadErr
	} else {
		return credentialFile, nil
	}
}

func getCredentialFromIni(profileName string, fromFactory *factory.Factory, credentialFile *ini.File) (*Credential, error) {
	log.Trace().Msg("Creating credential from ini.")
	loadedCredential, newErr := New(fromFactory,
		credentialFile.Section(profileName).Key("username").String(),
		credentialFile.Section(profileName).Key("password").String())

	if newErr != nil {
		return nil, newErr
	}

	return loadedCredential, nil
}

/*
func addAttributesFromIni(credentialFile *ini.File, loadedCredential *Credential) error {
	attributes, attributeErr := credentialFile.GetSection("attributes")

	if attributeErr != nil {
		if attributeErr.Error() != "section \"attributes\" does not exist" {
			return attributeErr
		} else {
			log.Trace().Msg("attributes section doesn't exist, skipping.")
		}
	} else {
		attributeKeys := attributes.KeyStrings()

		for i := 0; i < len(attributeKeys); i++ {
			attributeKey := attributeKeys[i]
			attributeValue := credentialFile.Section("attributes").Key(attributeKey).String()
			log.Trace().Str("Attribute Key", attributeKey).Msg("Loading value for attribute.")
			loadedCredential.attributes[attributeKey] = attributeValue
		}
	}

	return nil
}
*/
