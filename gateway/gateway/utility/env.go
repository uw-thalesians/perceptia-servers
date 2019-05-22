package utility

import (
	"errors"
	"os"
)

var ErrEnvVarNotSet = errors.New("environment variable not set")

// RequireEnv will retrieve the named environment variable or return an error if it was not set.
//
// Parameters
//		name: the name of the environment variable to check for
//
// Outputs
//		envValue: the value of the environment variable if set, or an empty string if not set
//		err: any errors that occur
func RequireEnv(name string) (envValue string, err error) {
	envValueLE, envSet := os.LookupEnv(name)
	if !envSet {
		return "", ErrEnvVarNotSet
	}
	return envValueLE, nil
}

// DefaultEnv checks for the named environment variable, if not set will use the provided default value.
//
// Parameters
//		name: the name of the environment variable to check for
//		defaultValue: the default value to use if the variable was not set
//
// Outputs
//		envValue: either the value set by the environment variable or the provided default value
//		envVarSet: if the named environment variable was set will be true, otherwise false
func DefaultEnv(name, defaultValue string) (envValue string, envVarSet bool) {
	envValueLE, envSet := os.LookupEnv(name)

	// Default to defaultValue if value not provided
	if !envSet {
		return defaultValue, false
	}
	return envValueLE, true
}
