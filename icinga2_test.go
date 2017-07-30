package main

import (
	"os"
	"reflect"
	"testing"
)

var (
	i Icinga2
)

const (
	HOST_NAME    = "google.com"
	SERVICE_NAME = "apt"
)

func init() {
	i = Icinga2{
		Host:     os.Getenv("HOST_ICINGA"),
		Username: os.Getenv("USERNAME_ICINGA"),
		Password: os.Getenv("PASSWORD_ICINGA"),
	}
}

func TestConstructFilter(t *testing.T) {
	expectedFilterOneHost := Filter{
		Filter: "match(\"*" + HOST_NAME + "*\", host.display_name)",
		Attrs: []string{
			"display_name",
			"name",
			"last_check",
			"state",
			"check_command",
		},
	}
	resultFilterOneHost := i.constructFilter(HOST_NAME, HOSTS, false)
	if !reflect.DeepEqual(expectedFilterOneHost, resultFilterOneHost) {
		t.Error(
			"For 'check one host status'",
			"expected", expectedFilterOneHost,
			"got", resultFilterOneHost,
		)
	}

	expectedFilterOneService := Filter{
		Filter: "match(\"*" + SERVICE_NAME + "*\", service.display_name)",
		Attrs: []string{
			"display_name",
			"name",
			"last_check",
			"state",
			"check_command",
		},
	}
	resultFilterOneService := i.constructFilter(SERVICE_NAME, SERVICES, false)
	if !reflect.DeepEqual(expectedFilterOneService, resultFilterOneService) {
		t.Error(
			"For 'check one service status'",
			"expected", expectedFilterOneService,
			"got", resultFilterOneService,
		)
	}

	expectedFilterAllHosts := Filter{
		Attrs: []string{
			"display_name",
			"name",
			"last_check",
			"state",
			"check_command",
		},
	}
	resultFilterAllHosts := i.constructFilter("", HOSTS, true)
	if !reflect.DeepEqual(expectedFilterAllHosts, resultFilterAllHosts) {
		t.Error(
			"For 'check all hosts status'",
			"expected", expectedFilterAllHosts,
			"got", resultFilterAllHosts,
		)
	}

	expectedFilterAllServices := Filter{
		Attrs: []string{
			"display_name",
			"name",
			"last_check",
			"state",
			"check_command",
		},
	}
	resultFilterAllServices := i.constructFilter("", SERVICES, true)
	if !reflect.DeepEqual(expectedFilterAllServices, resultFilterAllServices) {
		t.Error(
			"For 'check all hosts status'",
			"expected", expectedFilterAllServices,
			"got", resultFilterAllServices,
		)
	}
}

func TestCheck(t *testing.T) {
	_, err := i.check("", HOSTS, true)
	if err != nil {
		t.Error("Cannot get list of hosts")
	}

	_, err = i.check("", SERVICES, true)
	if err != nil {
		t.Error("Cannot get list of services")
	}

	_, err = i.check(SERVICE_NAME, SERVICES, false)
	if err != nil {
		t.Error("Cannot check the service")
	}

	_, err = i.check(HOST_NAME, HOSTS, false)
	if err != nil {
		t.Error("Cannot check the host")
	}
}
