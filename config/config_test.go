package config

import (
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
)

const configurations = `
[slack]
  token="d1524391-c0f5-4355-99a7-8201143fef21"

[icinga2]
  host="localhost"
  username="root"
  password="icinga"
`

func TestRead(t *testing.T) {
	var config Config
	toml.Decode(configurations, &config)

	expectedConfig := Config{}
	expectedConfig.Slack.Token = "d1524391-c0f5-4355-99a7-8201143fef21"
	expectedConfig.Icinga2.Host = "localhost"
	expectedConfig.Icinga2.Username = "root"
	expectedConfig.Icinga2.Password = "icinga"

	if !reflect.DeepEqual(expectedConfig, config) {
		t.Error(
			"expected", expectedConfig,
			"got", config,
		)
	}
}
