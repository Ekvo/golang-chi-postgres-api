// config  - parse data for run application
package config

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	"github.com/Ekvo/golang-chi-postgres-api/pkg/common"
)

var (
	// ErrConfigFieldEmpty - config field is empty
	ErrConfigFieldEmpty = errors.New("empty")

	// ErrConfigNoNumeric - field - only posotive numerci
	ErrConfigNoNumeric = errors.New("no numeric")
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`

	// DBNameForTest use only tests
	DBNameForTest string `mapstructure:"DB_TEST_NAME"`

	DBSSLMode string `mapstructure:"DB_SSLMODE"`

	// ServerHost - host for http.Server
	ServerHost string `mapstructure:"SRV_ADDR"`
}

// NewConfig - create Config
//
// load env from file if .env exist or from ENV
// set all ENV 'viper.BindEnv'
//
// test = true change DBName to DBNameForTest
func NewConfig(pathToEnv string, test bool) (*Config, error) {
	if err := godotenv.Load(pathToEnv); err != nil {
		log.Printf("config: .env file error - %v", err)
	}
	viper.AutomaticEnv()
	for _, env := range getNameENV() {
		if err := viper.BindEnv(env); err != nil {
			return nil, fmt.Errorf("config: ENV error - %w", err)
		}
	}
	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("config: create cfg error - %w", err)
	}
	if test {
		cfg.DBName = cfg.DBNameForTest
	}
	return cfg, cfg.validConfig()
}

// getNameENV - returns all ENV variable names
func getNameENV() []string {
	return []string{
		`DB_HOST`,
		`DB_PORT`,
		`DB_USER`,
		`DB_PASSWORD`,
		`DB_NAME`,
		`DB_TEST_NAME`,
		`DB_SSLMODE`,
		`SRV_ADDR`,
	}
}

func (cfg *Config) validConfig() error {
	msgErr := common.Message{}
	if cfg.DBHost == "" {
		msgErr["db-host"] = ErrConfigFieldEmpty
	}
	if port, err := strconv.Atoi(cfg.DBPort); err != nil || port < 1 {
		msgErr["db-port"] = ErrConfigNoNumeric
	}
	if cfg.DBUser == "" {
		msgErr["db-user"] = ErrConfigFieldEmpty
	}
	if cfg.DBPassword == "" {
		msgErr["db-password"] = ErrConfigFieldEmpty
	}
	if cfg.DBName == "" {
		msgErr["db-name"] = ErrConfigFieldEmpty
	}
	if cfg.DBSSLMode == "" {
		msgErr["db-ssl"] = ErrConfigFieldEmpty
	}
	if host, err := strconv.Atoi(cfg.ServerHost); err != nil || host < 1 {
		msgErr["server-host"] = ErrConfigNoNumeric
	}
	if len(msgErr) > 0 {
		return fmt.Errorf("config: invalid config - %s", msgErr.String())
	}
	return nil
}
