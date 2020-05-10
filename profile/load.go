package profile

import (
	"errors"
	"github.com/engi-fyi/go-credentials/factory"
	"github.com/engi-fyi/go-credentials/global"
	"github.com/rs/zerolog/log"
	"gopkg.in/ini.v1"
	"os"
)

/*
Load is responsible for loading a profile from file. If the file does not exist, an error is returned.
 */
func Load(profileName string, credentialFactory *factory.Factory) (*Profile, error) {
	log.Trace().Msg("Loading existing profile.")
	newProfile := Profile{
		Name:               profileName,
		ConfigFileLocation: credentialFactory.ConfigDirectory + profileName,
		attributes:         make(map[string]map[string]string),
		Initialized:        true,
		Factory:            credentialFactory,
	}

	log.Trace().Msg("Loading profile with details from file.")
	if newProfile.Factory.OutputType == global.OUTPUT_TYPE_INI {
		loadErr := newProfile.loadFromIni()

		if loadErr != nil {
			log.Trace().Err(loadErr).Msg("Error loading profile.")
			return nil, loadErr
		}

		return &newProfile, nil
	}

	return nil, errors.New(ERR_NOT_YET_IMPLEMENTED)
}

func (thisProfile *Profile) loadFromIni() error {
	if _, cfsErr := os.Stat(thisProfile.ConfigFileLocation); os.IsNotExist(cfsErr) {
		return errors.New(ERR_PROFILE_DID_NOT_EXIST)
	}

	return thisProfile.loadExistingIni()
}

func (thisProfile *Profile) loadExistingIni() error {
	configFile, loadErr := ini.Load(thisProfile.ConfigFileLocation)

	if loadErr != nil {
		return loadErr
	}

	for i := 0; i < len(configFile.Sections()); i++ {
		currSection := configFile.Sections()[i]
		keys := currSection.KeysHash()
		thisProfile.attributes[currSection.Name()] = keys
	}

	return nil
}
