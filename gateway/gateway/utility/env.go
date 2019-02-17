package utility

import (
	"fmt"
	"log"
	"os"
)

func RequireEnv(name string) (envVar string, err error) {
	envVar = os.Getenv(name)
	err = nil
	if len(envVar) == 0 {
		err = fmt.Errorf("RequireEnv: variable %s not set", name)
	}
	return envVar, err
}

func DefaultEnv(name, defaultValue string) string {
	envVar := os.Getenv(name)

	// Default to defaultValue if value not provided
	if len(envVar) == 0 {
		log.Printf("No %s env var set. Defaulting to %s", name, defaultValue)
		envVar = defaultValue
	}
	return envVar
}
