package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Doors []Door     `json:"doors"`
	Keys  []APIKey   `json:"keys"`
	HTTP  HTTPServer `json:"http"`
}

type Door struct {
	Name                 string `json:"name"`
	LiftCtlPin           int    `json:"lift_ctl_pin"`
	LiftCtlInactiveState int    `json:"lift_ctl_inactive_state"`
	LiftCtlTripMs        int    `json:"lift_ctl_trip_time_ms"`
	SensorPin            int    `json:"door_sensor_pin"`
	SensorClosedState    int    `json:"door_sensor_closed_state"`
}

type APIKey struct {
	Name            string   `json:"name"`
	Secret          string   `json:"secret"`
	MethodWhitelist []string `json:"allow_methods"`
}

type HTTPServer struct {
	IndexFile   string `json:"index_file"`
	ListenAddr  string `json:"listen_addr"`
	TLSCertFile string `json:"tls_cert"`
	TLSKeyFile  string `json:"tls_key"`
}

func Load(filepath string) (*Config, error) {
	cfgFile, err := os.Open(filepath)
	if err != nil {
		return &Config{}, err
	}

	defer cfgFile.Close()

	parsedCfg := &Config{}

	jsonParser := json.NewDecoder(cfgFile)
	if err := jsonParser.Decode(&parsedCfg); err != nil {
		return &Config{}, err
	}

	return parsedCfg, nil
}

// String enables Porter's configuration to be serialized as a JSON string
func (c *Config) String() string {
	bytes, err := json.Marshal(c)
	if err != nil {
		return ""
	}

	return string(bytes)
}

// GetDoorList returns all configured doors in a JSON string.
func (c *Config) GetDoorList() string {
	bytes, err := json.Marshal(c.Doors)
	if err != nil {
		return ""
	}

	return string(bytes)
}
