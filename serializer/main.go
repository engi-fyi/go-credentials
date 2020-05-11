package serializer

import (
	"errors"
	"github.com/engi-fyi/go-credentials/factory"
	"github.com/engi-fyi/go-credentials/global"
)

/*
New returns a Serializer object with all the defaults required to de(serialize) objects based on the Factory settings.
*/
func New(myFactory *factory.Factory, profileName string) *Serializer {
	return &Serializer{
		Factory:        myFactory,
		ProfileName:    profileName,
		CredentialFile: myFactory.CredentialFile,
		ConfigFile:     myFactory.ConfigDirectory + profileName,
		Initialized:    true,
	}
}

/*
Serialize is responsible for serializing an Credential and Profile, determining the file output type based on the value
of thisSerializer.Factory.OutputType. It is possible to serialize into multiple formats by initiating new factories, but
there is only one version of config with no extension. Every time a Serialize call is made, the file contents are
overwritten with the new values. Two formats cannot exist together.

The one exception to these rules is Environment, which doesn't save settings to file, although won't persists between
sessions.
*/
func (thisSerializer *Serializer) Serialize(username string, password string, attributes map[string]map[string]string) error {
	if thisSerializer.Factory.OutputType == global.OUTPUT_TYPE_INI {
		return thisSerializer.ToIni(username, password, attributes)
	} else if thisSerializer.Factory.OutputType == global.OUTPUT_TYPE_ENV {
		return thisSerializer.ToEnv(username, password, attributes)
	} else {
		return errors.New(ERR_UNRECOGNIZED_OUTPUT_TYPE)
	}
}

/*
Deserialize is responsible for deserializing an Credential and Profile, determining the file input type based on the
value of thisSerializer.Factory.OutputType.

For the format expected of each file, please see the appropriate From<Type> function.
*/
func (thisSerializer *Serializer) Deserialize() (string, string, map[string]map[string]string, error) {
	if thisSerializer.Factory.OutputType == global.OUTPUT_TYPE_INI {
		return thisSerializer.FromIni()
	} else if thisSerializer.Factory.OutputType == global.OUTPUT_TYPE_ENV {
		return thisSerializer.FromEnv()
	} else {
		return "", "", make(map[string]map[string]string), errors.New(ERR_UNRECOGNIZED_OUTPUT_TYPE)
	}
}
