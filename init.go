package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	wisdomType "github.com/wisdom-oss/commonTypes"

	"github.com/joho/godotenv"

	"github.com/qustavo/dotsql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/wisdom-oss/service-water-rights/config"
	"github.com/wisdom-oss/service-water-rights/globals"
)

// customTypes contains the PostgreSQL names of the custom types that are used
// by the microservice. these custom types are registered during the execution
// of connectDatabase.
var customTypes = []string{
	"water_rights.numeric_keyed_value",
	"water_rights.quantity",
	"water_rights.rate",
	"water_rights.dam_target",
	"water_rights.land_record",
	"water_rights.legal_department",
	"water_rights.injection_limit",
}

var initLogger = log.With().Bool("startup", true).Logger()

// init is executed at every startup of the microservice and is always executed
// before main
func init() {
	// load the variables found in the .env file into the process environment
	err := godotenv.Load()
	if err != nil {
		initLogger.Debug().Msg("no .env files found")
	}

	configureLogger()
	loadServiceConfiguration()
	connectDatabase()
	loadPreparedQueries()

	initLogger.Info().Msg("initialization process finished")

}

// configureLogger handles the configuration of the logger used in the
// microservice. it reads the logging level from the `LOG_LEVEL` environment
// variable and sets it according to the parsed logging level. if an invalid
// value is supplied or no level is supplied, the service defaults to the
// `INFO` level
func configureLogger() {
	// allow the logger to create an error stack for the logs
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// now use the environment variable `LOG_LEVEL` to determine the logging
	// level for the microservice.
	rawLoggingLevel, isSet := os.LookupEnv("LOG_LEVEL")

	// if the value is not set, use the info level as default.
	var loggingLevel zerolog.Level
	if !isSet {
		loggingLevel = zerolog.InfoLevel
	} else {
		// now try to parse the value of the raw logging level to a logging
		// level for the zerolog package
		var err error
		loggingLevel, err = zerolog.ParseLevel(rawLoggingLevel)
		if err != nil {
			// since an error occurred while parsing the logging level, use info
			loggingLevel = zerolog.InfoLevel
			initLogger.Warn().Msg("unable to parse value from environment. using info")
		}
	}
	// since now a logging level is set, configure the logger
	zerolog.SetGlobalLevel(loggingLevel)
}

// loadServiceConfiguration handles loading the `environment.json` file which
// describes which environment variables are needed for the service to function
// and what variables are optional and their default values
func loadServiceConfiguration() {
	initLogger.Info().Msg("loading service configuration from environment")
	// now check if the default location for the environment configuration
	// was changed via the `ENV_CONFIG_LOCATION` variable
	location, locationChanged := os.LookupEnv("ENV_CONFIG_LOCATION")
	if !locationChanged {
		// since the location has not changed, set the default value
		location = config.EnvironmentFilePath
		initLogger.Debug().Msg("location for environment config not changed")
	}
	initLogger.Debug().Str("path", location).Msg("loading environment requirements file")
	var c wisdomType.EnvironmentConfiguration
	err := c.PopulateFromFilePath(location)
	if err != nil {
		initLogger.Fatal().Err(err).Msg("unable to load environment requirements file")
	}
	initLogger.Info().Msg("validating environment variables")
	globals.Environment, err = c.ParseEnvironment()
	if err != nil {
		initLogger.Fatal().Err(err).Msg("environment validation failed")
	}
	initLogger.Info().Msg("loaded service configuration from environment")
}

// connectDatabase uses the previously read environment variables to connect the
// microservice to the PostgreSQL database used as the backend for all WISdoM
// services
func connectDatabase() {
	initLogger.Info().Msg("connecting to the database")

	address := fmt.Sprintf("postgres://%s:%s@%s:%s/wisdom",
		globals.Environment["PG_USER"], globals.Environment["PG_PASS"],
		globals.Environment["PG_HOST"], globals.Environment["PG_PORT"])

	var err error
	pgxConfig, err := pgxpool.ParseConfig(address)
	if err != nil {
		initLogger.Fatal().Err(err).Msg("unable to create base configuration for connection pool")
	}
	pgxConfig.AfterConnect = registerCustomTypes
	globals.Db, err = pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		initLogger.Fatal().Err(err).Msg("unable to create database connection pool")
	}
	err = globals.Db.Ping(context.Background())
	if err != nil {
		initLogger.Fatal().Err(err).Msg("unable to verify the connection to the database")
	}
	initLogger.Info().Msg("database connection established")
}

// loadPreparedQueries loads the prepared SQL queries from a file specified by
// the QUERY_FILE_LOCATION environment variable.
// It initializes the SqlQueries variable with the loaded queries.
// If there is an error loading the queries, it logs a fatal error and the
// program terminates.
// This function is typically called during the startup of the microservice.
func loadPreparedQueries() {
	initLogger.Info().Msg("loading prepared sql queries")
	location, locationChanged := os.LookupEnv("QUERY_FILE_LOCATION")
	if !locationChanged {
		// since the location has not changed, set the default value
		location = config.QueryFilePath
		initLogger.Debug().Msg("location for environment config not changed")
	}
	var err error
	globals.SqlQueries, err = dotsql.LoadFromFile(location)
	if err != nil {
		initLogger.Fatal().Err(err).Msg("failed to load prepared queries")
	}
}

// registerCustomTypes registers custom types in the database connection.
// It takes the context and a connection object as arguments.
// It iterates over the custom types array and loads each custom type
// from the database. Then, it registers both the custom type and the
// array version of the type using the connection's TypeMap.
// If any error occurs during the process, it logs an error and returns
// the error. Otherwise, it returns nil.
func registerCustomTypes(ctx context.Context, conn *pgx.Conn) error {
	for _, customTypeName := range customTypes {
		// create the array version of the type
		customTypeArray := fmt.Sprintf("%s[]", customTypeName)
		// now load the custom type from the database
		customType, err := conn.LoadType(ctx, customTypeName)
		if err != nil {
			initLogger.Error().Err(err).Str("type", customTypeName).Msg("unable to load custom type from the database")
			return err
		}
		conn.TypeMap().RegisterType(customType)

		// now load the array version of the type
		customType, err = conn.LoadType(ctx, customTypeArray)
		if err != nil {
			initLogger.Error().Err(err).Str("type", customTypeArray).Msg("unable to load custom type from the database")
			return err
		}
		conn.TypeMap().RegisterType(customType)
	}
	return nil
}
