package profile

import (
	"errors"
	"github.com/engi-fyi/go-credentials/factory"
	"github.com/engi-fyi/go-credentials/global"
	"github.com/rs/zerolog/log"
	"gopkg.in/ini.v1"
	"os"
)

func Load(profileName string, credentialFactory *factory.Factory) (*Profile, error) {
	newProfile := Profile{
		Name: profileName,
		ConfigFileLocation: credentialFactory.ConfigDirectory + profileName,
		attributes: make(map[string]map[string]string),
		Initialized: true,
		Factory: credentialFactory,
	}

	if newProfile.Factory.OutputType == global.OUTPUT_TYPE_INI {
		loadErr := newProfile.loadFromIni()

		if loadErr != nil {
			return nil, loadErr
		}

		return &newProfile, nil
	}

	return nil, nil
}

func (thisProfile *Profile) loadFromIni() error {
	if _, cfsErr := os.Stat(thisProfile.ConfigFileLocation); os.IsNotExist(cfsErr) {
		return thisProfile.initNewIni()
	} else {
		return thisProfile.loadExistingIni()
	}
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

func (thisProfile *Profile) initNewIni() error {
	log.Trace().Str("config_file", thisProfile.ConfigFileLocation).Msg("Initializing config file.")
	emptyFile, efErr := os.Create(thisProfile.ConfigFileLocation)

	if efErr != nil {
		log.Error().Err(efErr).Msg("Error creating config file.")
		return nil
	}

	closeErr := emptyFile.Close()

	if closeErr != nil {
		return closeErr
	} else {
		return errors.New(ERR_FILE_DID_NOT_EXIST_BLANK_FILE_INIT)
	}
}