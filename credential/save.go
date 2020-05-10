package credential

import (
	"errors"
	"github.com/engi-fyi/go-credentials/factory"
	"github.com/engi-fyi/go-credentials/global"
	"github.com/engi-fyi/go-credentials/profile"
	"github.com/rs/zerolog/log"
)

/*
Save is responsible for saving the credential at ~/.application_name/credentials in the specified output format
that has been set on the Credentials' Factory object.
*/
// TODO(7): Implement json format.
func (thisCredential *Credential) Save() error {
	if !thisCredential.Initialized {
		return errors.New(ERR_CREDENTIAL_NOT_INITIALIZED)
	}

	if thisCredential.Factory.Initialized {
		if !thisCredential.Profile.Initialized {
			return errors.New(profile.ERR_PROFILE_NOT_INITIALIZED)
		}

		saveErr := thisCredential.Profile.Save()

		if saveErr != nil {
			return saveErr
		}


		if thisCredential.Factory.OutputType == global.OUTPUT_TYPE_JSON {
			return thisCredential.saveJson()
		} else if thisCredential.Factory.OutputType == global.OUTPUT_TYPE_INI {
			credErr := thisCredential.saveIni()

			if credErr != nil {
				return credErr
			}
		} else {
			return errors.New(factory.ERR_FACTORY_INCONSISTENT_STATE)
		}
	} else {
		return errors.New(factory.ERR_FACTORY_NOT_INITIALIZED)
	}

	return nil
}

// BUG(4): Respect alternates when saving to file.
// TODO(1000): Should we actually load an existing credential or would the user expect the whole state to be serialized?
func (thisCredential *Credential) saveIni() error {
	log.Trace().Msg("Saving credentials as ini.")
	credentialFile, loadErr := loadCredentialIni(thisCredential.Factory)

	if loadErr != nil {
		return loadErr
	}

	log.Trace().Msg("Existing credentials loaded from file.")

	credentialFile.Section(thisCredential.Profile.Name).Key("username").SetValue(thisCredential.Username)
	credentialFile.Section(thisCredential.Profile.Name).Key("password").SetValue(thisCredential.Password)
	log.Trace().Msg("Username and password set in ini.")
	saveErr := credentialFile.SaveTo(thisCredential.Factory.CredentialFile)

	if saveErr != nil {
		log.Error().Err(saveErr).Msg("Error saving credentials ini file.")
		return saveErr
	}

	log.Trace().Msg("Credential ini file saved successfully.")

	return nil
}

// TODO(7): Implement json format.
func (thisCredential *Credential) saveJson() error {
	return errors.New(ERR_JSON_FUNCTIONALITY_NOT_IMPLEMENTED)
}
