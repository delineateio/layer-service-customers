package config

import (
	"fmt"
	"reflect"

	"delineate.io/customers/src/logging"
	"github.com/mitchellh/consulstructure"
	"github.com/spf13/viper"
)

// Create a configuration struct that'll be filled by Consul.
type Config struct {
	Addr     string
	DataPath string `consul:"data_path"`
}

// Create our decoder
var updateCh = make(chan interface{})
var errCh = make(chan error)
var decoder = &consulstructure.Decoder{
	Target:   &Config{},
	Prefix:   "services/customers",
	UpdateCh: updateCh,
	ErrCh:    errCh,
}

func watch() {
	go decoder.Run()
	for {
		select {
		case <-updateCh:
			if serverChanged, err := readConfig(); err != nil {
				logging.Err(err)
			} else {
				if serverChanged {
					logging.Warn("server config can't be changed dynamically")
				} else {
					logging.SetLevel()
				}
			}
		case err := <-errCh:
			fmt.Printf("Error: %s\n", err)
		}
	}
}

func readConfig() (bool, error) {
	var serverChanged bool
	server := GetSection("server")
	err := viper.ReadRemoteConfig()
	if err != nil {
		logging.Err(err)
	} else {
		logging.Info("config successfully loaded from consul")
		serverChanged = !reflect.DeepEqual(server, GetSection("server"))
	}
	return serverChanged, err
}
