// Package internal implements utilities for accessing and loading the
// configuration as well as connecting to the databases
//
// The internal package may only be used with the given functions.
// Manipulating values from outside of the internal package may result in
// unexpected behavior inside the application.
package internal

import (
	"errors"
	"strings"

	_ "embed"

	"github.com/spf13/viper"
)

// Configuration contains the parsed configuration file with the.
var configuration *viper.Viper

// Generic errors returned if something happened during the loading of the
// configuration.
var (
	ErrNoConfigurationFile        = errors.New("no configuration file found")
	ErrConfigurationUnreadable    = errors.New("configuration file not readable")
	ErrConfigurationNotCreateable = errors.New("unable to create configuration file")
)

// Default values used in the configuration.
const (
	defaultPortHttp              = 8000
	defaultPortPostgres          = 5432
	defaultPostgresSSLMode       = "disable"
	defaultPostgresDatabase      = "wisdom"
	defaultOIDCAuthority         = "http://backend/api/auth/"
	defaultAuthorizationRequired = true
)

// Keys for common configuration entries.
const (
	ConfigKey_Postgres_User         = "postgres.user"
	ConfigKey_Postgres_Password     = "postgres.password"
	ConfigKey_Postgres_Host         = "postgres.host"
	ConfigKey_Postgres_Port         = "postgres.port"
	ConfigKey_Postgres_SSLMode      = "postgres.sslmode"
	ConfigKey_Postgres_Database     = "postgres.database"
	ConfigKey_Http_Host             = "http.host"
	ConfigKey_Http_Port             = "http.port"
	ConfigKey_Oidc_Authority        = "oidc.authority"
	ConfigKey_Require_Authorization = "authorization.required"
)

// envAliases contains all allowed environment variable names that are used to
// gather the values for the given keys.
var envAliases = map[string][]string{
	ConfigKey_Postgres_User: {"PGUSER", "PG_USER", "POSTGRES_USER", "DB_USER"},
	ConfigKey_Postgres_Password: {"PGPASSWORD", "PG_PASSWORD", "PGPASS", "PG_PASS", "POSTGRES_PASS",
		"POSTGRES_PASSWORD", "DB_PASS", "DB_PASSWORD"},
	ConfigKey_Postgres_Host:         {"PGHOST", "PG_HOST", "POSTGRES_HOST", "DB_HOST"},
	ConfigKey_Postgres_Port:         {"PGPORT", "PG_PORT", "POSTGRES_PORT", "DB_PORT"},
	ConfigKey_Postgres_Database:     {"PGDATABASE", "PG_DATABASE", "POSTGRES_DATABASE", "DB_DATABASE"},
	ConfigKey_Postgres_SSLMode:      {"PGSSLMODE", "PG_SSLMODE", "POSTGRES_SSLMODE", "DB_SSLMODE"},
	ConfigKey_Oidc_Authority:        {"OIDC_AUTHORITY"},
	ConfigKey_Require_Authorization: {"AUTH_REQUIRED"},
}

// ParseConfiguration initializes the [Configuration] variable and reads the
// configuration file.
// If no configuration file is found, one is written into the current working
// directory containing an examplatory configuration file (in the TOML format).
// That configuration file needs to be edited to contain the correct values.
func ParseConfiguration() error {
	configuration = initializeViperInstance()
	setDefaults(configuration)
	bindEnvironmentVariables(configuration)
	_ = configuration.ReadInConfig()
	return nil
}

// Configuration returns the unexported [*viper.Viper] variable to allow using
// the parsed configuration.
func Configuration() *viper.Viper {
	return configuration
}

// initializeViperInstance creates a new [*viper.Viper] instance.
func initializeViperInstance() *viper.Viper {
	instance := viper.New()
	instance.SetConfigName("config")
	instance.AddConfigPath("/etc/award/")
	instance.AddConfigPath("/run/secrets/")
	instance.AddConfigPath("$HOME/.config/award/")
	instance.AddConfigPath(".")
	instance.SetEnvPrefix("")
	instance.SetEnvKeyReplacer(strings.NewReplacer(".", "__", "-", "_"))
	instance.AutomaticEnv()
	return instance
}

// setDefaults sets up some required default values
//
// the defaults set here are used to ensure the application as a chance at
// starting without requiring a more enhanced configuration.
// this enabled users of the backend to not worry about the http port and
// default database ports.
func setDefaults(instance *viper.Viper) {
	// setup defaults for the http server part
	instance.SetDefault(ConfigKey_Http_Host, "")
	instance.SetDefault(ConfigKey_Http_Port, defaultPortHttp)

	// setup defaults for the database communication
	instance.SetDefault(ConfigKey_Postgres_Port, defaultPortPostgres)
	instance.SetDefault(ConfigKey_Postgres_SSLMode, defaultPostgresSSLMode)
	instance.SetDefault(ConfigKey_Postgres_Database, defaultPostgresDatabase)

	// setup some defaults for communicating with the user management if needed
	instance.SetDefault(ConfigKey_Oidc_Authority, defaultOIDCAuthority)

	// setup the authorization to be enabled by default (however only in release
	// routers this will work)
	instance.SetDefault(ConfigKey_Require_Authorization, defaultAuthorizationRequired)

}

// bindEnvironmentVariables binds commonly used environment varialbes to
// the configuration values.
func bindEnvironmentVariables(instance *viper.Viper) {
	for key, envVars := range envAliases {
		args := make([]string, len(envVars)+1)
		args[0] = key
		for idx, envVar := range envVars {
			args[idx+1] = envVar
		}
		_ = instance.BindEnv(args...)
	}
}
