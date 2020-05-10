package serializer

import (
	"errors"
	"github.com/engi-fyi/go-credentials/factory"
	"github.com/engi-fyi/go-credentials/global"
)

func New(myFactory *factory.Factory, profileName string) *Serializer {
	return &Serializer{
		Factory:        myFactory,
		ProfileName:    profileName,
		CredentialFile: myFactory.CredentialFile,
		ConfigFile: 	myFactory.ConfigDirectory + profileName,
		Initialized:	true,
	}
}

func (thisSerializer *Serializer) Serialize(username string, password string, attributes map[string]map[string]string) error {
	if thisSerializer.Factory.OutputType == global.OUTPUT_TYPE_INI {
		return thisSerializer.ToIni(username, password, attributes)
	} else if thisSerializer.Factory.OutputType == global.OUTPUT_TYPE_ENV {
		return thisSerializer.ToEnv(username, password, attributes)
	} else {
		return errors.New(ERR_UNRECOGNIZED_OUTPUT_TYPE)
	}
}

func (thisSerializer *Serializer) Deserialize() (string, string, map[string]map[string]string, error) {
	if thisSerializer.Factory.OutputType == global.OUTPUT_TYPE_INI {
		return thisSerializer.FromIni()
	} else if thisSerializer.Factory.OutputType == global.OUTPUT_TYPE_ENV {
		return thisSerializer.FromEnv()
	} else {
		return "", "", make(map[string]map[string]string), errors.New(ERR_UNRECOGNIZED_OUTPUT_TYPE)
	}
}
