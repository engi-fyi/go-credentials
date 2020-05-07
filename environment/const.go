package environment

const ERR_CANT_LOAD_WITH_EMPTY_VALUES = "sorry values required for application name and username, cannot deploy"
const ERR_KEY_MUST_MATCH_REGEX = "sorry the key must only include letters and underscores [0-9A-Za-z_]"
const ERR_REQUIRED_VARIABLE_USERNAME_NOT_FOUND = "username has not been set via the environment and is required, load failed"
const ERR_REQUIRED_VARIABLE_PASSWORD_NOT_FOUND = "password has not been set via the environment and is required, load failed"

const REGEX_KEY_NAME = "(?m)^[0-9A-Za-z_]+$"
