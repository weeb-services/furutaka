package main

import (
	"github.com/teris-io/shortid"
	"net/http"
	"time"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"bytes"
)

type Registrator struct {
	id                                      string
	name, environment, url, token, endpoint string
	port                                    uint32
	active                                  bool
}

type ConsulRegistration struct {
	Id   string   `json:"ID"`
	Name string   `json:"Name"`
	Port uint32   `json:"Port"`
	Tags []string `json:"Tags"`
}

func NewRegistrator(name string, environment string, port uint32, endpoint string, token string) Registrator {
	id, _ := shortid.Generate()
	return Registrator{id: id, name: name, environment: environment, port: port, endpoint: endpoint, token: token, active: true}
}

func (r Registrator) Register() error {
	client := &http.Client{Timeout: time.Second * 10}
	tags := make([]string, 1)
	tags[0] = r.environment
	cr := &ConsulRegistration{Id: fmt.Sprintf("%v-%v", r.name, r.id), Name: r.name, Port: r.port, Tags: tags}
	jsonRegistration, _ := json.Marshal(cr)
	req, err := http.NewRequest("PUT", fmt.Sprintf("http://%v/v1/agent/service/register", r.endpoint), bytes.NewBuffer(jsonRegistration))
	req.Header.Set("Content-Type", "application/json")
	if r.token != "" {
		req.Header.Add("X-Consul-Token", r.token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	return err
}

func (r Registrator) Unregister() error {
	client := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest("PUT", fmt.Sprintf("http://%v/v1/agent/service/deregister/%v-%v", r.endpoint, r.name, r.id), nil)
	if r.token != "" {
		req.Header.Add("X-Consul-Token", r.token)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return err
}
