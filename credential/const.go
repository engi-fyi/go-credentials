package credential

const ATTRIBUTE_USERNAME = "username"
const ATTRIBUTE_PASSWORD = "password"

const ERR_ATTRIBUTE_NOT_EXIST = "sorry that attribute hasn't been set on this credential"
const ERR_USERNAME_OR_PASSWORD_NOT_SET = "sorry you must set a username or password"
const ERR_FACTORY_MUST_BE_INITIALIZED = "sorry the factory must be created correctly before proceeding"
const ERR_KEY_MUST_MATCH_REGEX = "sorry the key must only include numbers, letters and underscores [0-9A-Za-z_]"
const ERR_FACTORY_PRIVATE_ATTEMPT_DEPLOY = "you have attempted to deploy environment variables but your factory is private"
const ERR_JSON_FUNCTIONALITY_NOT_IMPLEMENTED = "sorry, that feature has not been implemented yet"

const REGEX_KEY_NAME = "(?m)^[0-9A-Za-z_]+$"
