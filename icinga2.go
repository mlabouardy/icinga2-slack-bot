package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Icinga2 struct {
	Host     string
	Username string
	Password string
}

type ObjectType string

const (
	SERVICES = "services"
	HOSTS    = "hosts"
)

type Result struct {
	Results []Object `json:"results"`
}

type Object struct {
	Attrs Attribute `json:"attrs"`
}

type Attribute struct {
	CheckCommand string      `json:"check_command"`
	DisplayName  string      `json:"display_name"`
	Name         string      `json:"name"`
	State        float32     `json:"state"`
	CheckTime    json.Number `json:"last_check,float"`
	HostName     string      `json:"host_name"`
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

func (i *Icinga2) check(name string, objectType ObjectType, all bool) (Result, error) {

	server := fmt.Sprintf("https://%s:5665/v1/objects/%s", i.Host, objectType)
	queryParams := "?attrs=name&attrs=state&attrs=display_name&attrs=check_command&attrs=last_check"

	if !all {
		queryParams = fmt.Sprintf("(%%22%s%%22,%s.display_name)&attrs=name&attrs=state&attrs=display_name&attrs=check_command&attrs=last_check", name, strings.TrimSuffix(string(objectType), "s"))
	}

	url := fmt.Sprintf("%s%s", server, queryParams)
	if objectType == SERVICES {
		url = fmt.Sprintf("%s%s&attrs=host_name", server, queryParams)
	}

	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(i.Username, i.Password)
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

	if !all && len(result.Results) == 0 {
		return result, errors.New(name + " not found")
	}

	return result, nil
}

func (i *Icinga2) HostStatus(name string) (Result, error) {
	return i.check(name, HOSTS, false)
}

func (i *Icinga2) ServiceStatus(name string) (Result, error) {
	return i.check(name, SERVICES, false)
}

func (i *Icinga2) ListServices() (Result, error) {
	return i.check("", SERVICES, true)
}

func (i *Icinga2) ListHosts() (Result, error) {
	return i.check("", HOSTS, true)
}
