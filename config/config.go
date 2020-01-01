// Package config provides a central facility for loading and parsing configuration files, along with schemas
// for various Porter component configuration objects
package config

import (
	"encoding/json"
	"os"
)

// Config contains the full contents of a Porter configuration file
type Config struct {
	Doors []Door     `json:"doors"`
	Keys  []APIKey   `json:"keys"`
	HTTP  HTTPServer `json:"http"`
}

// Door contains the configuration for a door.Door as read from a configuration file
type Door struct {
	Name                 string `json:"name"`
	LiftCtlPin           int    `json:"lift_ctl_pin"`
	LiftCtlInactiveState int    `json:"lift_ctl_inactive_state"`
	LiftCtlTripMs        int    `json:"lift_ctl_trip_time_ms"`
	SensorPin            int    `json:"door_sensor_pin"`
	SensorClosedState    int8   `json:"door_sensor_closed_state"`
}

// APIKey contains an API key definition as read from a configuration file
type APIKey struct {
	Name            string   `json:"name"`
	Secret          string   `json:"secret"`
	MethodWhitelist []string `json:"allow_methods"`
}

// HTTPServer contains the API's HTTP server config as read from a configuration file
type HTTPServer struct {
	IndexFile   string `json:"index_file"`
	ListenAddr  string `json:"listen_addr"`
	TLSCertFile string `json:"tls_cert"`
	TLSKeyFile  string `json:"tls_key"`
}

// Load loads and parses a Porter configuration file at the provided path
func Load(filepath string) (*Config, error) {
	cfgFile, err := os.Open(filepath)
	if err != nil {
		return &Config{}, err
	}

	defer func() {
		err := cfgFile.Close()
		if err != nil {
			panic(err)
		}
	}()

	parsedCfg := &Config{}

	jsonParser := json.NewDecoder(cfgFile)
	if err := jsonParser.Decode(&parsedCfg); err != nil {
		return &Config{}, err
	}

	return parsedCfg, nil
}

// String serializes Porter's configuration as a JSON string
func (c *Config) String() string {
	bytes, err := json.Marshal(c)
	if err != nil {
		return ""
	}

	return string(bytes)
}
