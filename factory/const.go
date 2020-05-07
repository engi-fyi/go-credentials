package factory

const OUTPUT_TYPE_JSON = "json"
const OUTPUT_TYPE_INI = "ini"
const OUTPUT_TYPE_INVALID = "nri"

const ERR_INVALID_OUTPUT_TYPE = "sorry the output type you have set is not valid and not supported"
const ERR_FACTORY_INCONSISTENT_STATE = "the Factory is in an inconsistent state, please only use public methods to modify"
const ERR_APPLICATION_NAME_BLANK = "sorry the application name must not be blank"
const ERR_KEY_MUST_MATCH_REGEX = "sorry the key must only include letters and underscores [0-9A-Za-z_]"
const ERR_FACTORY_NOT_INITIALIZED = "sorry the factory has not been initilalized please run f.Intialize() and try again"
const ERR_ALTERNATE_USERNAME_CANNOT_BE_BLANK = "the alternate username you provided is blank, please provide a string"
const ERR_ALTERNATE_PASSWORD_CANNOT_BE_BLANK = "the alternate password you provided is blank, please provide a string"

const REGEX_KEY_NAME = "(?m)^[0-9A-Za-z_]+$"
