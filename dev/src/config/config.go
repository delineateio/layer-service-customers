package config

import (
	"fmt"
	"os"

	"delineate.io/customers/src/logging"
	retry "github.com/avast/retry-go/v3"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

const AddressEnvName = "CONSUL_HTTP_ADDR"
const DefaultAddress = "localhost:8500"

func getConsulEndpoint() string {
	value, present := os.LookupEnv(AddressEnvName)
	if present {
		return value
	}
	return DefaultAddress
}

// Viper watching remote config backend Consul
// https://gist.github.com/andrewmeissner/f9524709e0a15b336c174ddbaf82a2d3
func Initialize() {
	// retry added to stop issues on startup when if the Consul config
	// is immediately uploaded
	err := retry.Do(
		func() error {
			err := viper.AddRemoteProvider("consul", getConsulEndpoint(), "services/customers")
			if err != nil {
				return err
			}
			viper.SetConfigType("yaml")
			err = viper.ReadRemoteConfig()
			if err != nil {
				return err
			}
			go watch()
			logging.Info("successfully initialized config")
			return nil
		},
	)
	if err != nil {
		logging.Err(err)
		panic(err)
	}
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetStringOrDefault(key, value string) string {
	logging.Debug(fmt.Sprintf("requested '%s' config key", key))

	if viper.IsSet(key) {
		value = viper.GetString(key)
		logging.Debug(fmt.Sprintf("found '%s' value for config key '%s'", value, key))
		return value
	}
	logging.Warn(fmt.Sprintf("used '%s' as default value for config key '%s'", value, key))
	return value
}

func GetSection(key string) map[string]string {
	logging.Debug(fmt.Sprintf("requested '%s' config key", key))
	return viper.GetStringMapString(key)
}

func GetSlice(key string) []string {
	return viper.GetStringSlice(key)
}
