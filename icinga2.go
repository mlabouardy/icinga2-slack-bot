package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Icinga2 represents an icinga2 server
type Icinga2 struct {
	Host     string
	Username string
	Password string
}

type ObjectType string

// Type of element to monitor in icinga
const (
	SERVICES = "services"
	HOSTS    = "hosts"
)

// Result will encapsulate the icinga2 api json response
type Result struct {
	Results []Object `json:"results"`
}

// Object represents a single host/service
type Object struct {
	Attrs Attribute `json:"attrs"`
}

// Attribute represents properties of host/service
type Attribute struct {
	CheckCommand string      `json:"check_command"`
	DisplayName  string      `json:"display_name"`
	Name         string      `json:"name"`
	State        float32     `json:"state"`
	CheckTime    interface{} `json:"last_check,string"`
	HostName     string      `json:"host_name"`
}

// Filter represents filters to pass to icinga2 api request
type Filter struct {
	Filter string   `json:"filter,omitempty"`
	Attrs  []string `json:"attrs"`
}

var (
	client *http.Client
)

func init() {
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func (i *Icinga2) constructFilter(name string, objectType ObjectType, checkAll bool) Filter {
	filter := Filter{
		Attrs: []string{
			"display_name",
			"name",
			"last_check",
			"state",
			"check_command",
		},
	}

	if !checkAll {
		if objectType == HOSTS {
			filter.Filter = "match(\"*" + name + "*\", host.display_name)"
		} else {
			filter.Filter = "match(\"*" + name + "*\", service.display_name)"
		}
	}

	return filter
}

func (i *Icinga2) check(name string, objectType ObjectType, checkAll bool) (Result, error) {
	url := fmt.Sprintf("https://%s:5665/v1/objects/%s", i.Host, objectType)

	filter := i.constructFilter(name, objectType, checkAll)

	b := new(bytes.Buffer)

	json.NewEncoder(b).Encode(filter)

	req, err := http.NewRequest("POST", url, b)
	req.SetBasicAuth(i.Username, i.Password)
	req.Header.Set("X-HTTP-Method-Override", "GET")

	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var result Result

	json.Unmarshal(body, &result)

	if !checkAll && len(result.Results) == 0 {
		return result, errors.New(name + " not found")
	}

	return result, nil
}

// HostStatus calls check method to get the status of a single host
func (i *Icinga2) HostStatus(name string) (Result, error) {
	return i.check(name, HOSTS, false)
}

// ServiceStatus calls check method to get the status of a single service
func (i *Icinga2) ServiceStatus(name string) (Result, error) {
	return i.check(name, SERVICES, false)
}

// ListServices calls check to get the status of all the services
func (i *Icinga2) ListServices() (Result, error) {
	return i.check("", SERVICES, true)
}

// ListHosts calls check to get the status of all the hosts
func (i *Icinga2) ListHosts() (Result, error) {
	return i.check("", HOSTS, true)
}
