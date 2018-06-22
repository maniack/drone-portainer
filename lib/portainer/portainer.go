package portainer

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/goware/urlx"
)

type Env struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Stack struct {
	Id          int    `json:"Id"`
	Name        string `json:"Name"`
	Type        int    `json:"Type"`
	EndpointID  int    `json:"EndpointID"`
	EntryPoint  string `json:"EntryPoint"`
	SwarmID     string `json:"SwarmID"`
	ProjectPath string `json:"ProjectPath"`
	Env         []*Env `json:"Env"`
}

type Endpoint struct {
	Id        int    `json:"Id"`
	Name      string `json:"Name"`
	Type      int    `json:"Type"`
	URL       string `json:"URL"`
	GroupId   int    `json:"GroupId"`
	PublicURL string `json:"PublicURL"`
	SwarmID   string `json:"SwarmID,omitempty"`
}

type Portainer struct {
	client  *http.Client
	address string
	auth    string
}

func NewPortainer(address string, insecure bool) (*Portainer, error) {
	tlsconfig := &tls.Config{InsecureSkipVerify: insecure}
	transport := &http.Transport{TLSClientConfig: tlsconfig}
	client := &http.Client{Transport: transport}

	url, err := urlx.Parse(address)
	if err != nil {
		return nil, fmt.Errorf("Address parsing error : %s", err)
	}

	if url.Scheme == "" {
		switch url.Port() {
		case "443":
			url.Scheme = "https://"
			break
		default:
			url.Scheme = "https://"
		}
	}

	address, err = urlx.Normalize(url)
	if err != nil {
		return nil, fmt.Errorf("Address normalizing error : %s", err)
	}

	return &Portainer{
		client:  client,
		address: address,
	}, nil
}

func (self *Portainer) Connect() error {
	req, err := http.NewRequest("HEAD", fmt.Sprintf("%s/", self.address), nil)
	if err != nil {
		return err
	}

	rsp, err := self.client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode >= 400 {
		return fmt.Errorf("Portainer API error: %s %s %s", req.Method, req.URL.String(), rsp.Status)
	}

	return nil
}

func (self *Portainer) Auth(user, pass string) error {
	args, err := json.Marshal(&struct {
		Username string
		Password string
	}{
		Username: user,
		Password: pass,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/auth", self.address), bytes.NewBuffer(args))
	if err != nil {
		return err
	}

	rsp, err := self.client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode > 200 {
		return fmt.Errorf("Portainer API error: %s %s %s", req.Method, req.URL.String(), rsp.Status)
	}

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	var auth struct {
		JWT string `json:"jwt"`
	}

	err = json.Unmarshal(data, &auth)
	if err != nil {
		return err
	}

	self.auth = auth.JWT

	return nil
}

func (self *Portainer) GetEndpointByName(endpoint string) (*Endpoint, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/endpoints", self.address), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", self.auth))

	rsp, err := self.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode > 200 {
		return nil, fmt.Errorf("Portainer API error: %s %s %s", req.Method, req.URL.String(), rsp.Status)
	}

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	var endpoints []*Endpoint

	err = json.Unmarshal(data, &endpoints)
	if err != nil {
		return nil, err
	}

	for _, e := range endpoints {
		if e.Name == endpoint {
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/endpoints/%d/docker/swarm", self.address, e.Id), nil)
			if err != nil {
				return nil, err
			}
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", self.auth))

			rsp, err := self.client.Do(req)
			if err != nil {
				return nil, err
			}
			defer rsp.Body.Close()

			if rsp.StatusCode > 200 {
				return nil, fmt.Errorf("Portainer API error: %s %s %s", req.Method, req.URL.String(), rsp.Status)
			}

			data, err := ioutil.ReadAll(rsp.Body)
			if err != nil {
				return nil, err
			}

			var swarm struct {
				ID string
			}

			err = json.Unmarshal(data, &swarm)
			if err != nil {
				return nil, err
			}

			e.SwarmID = swarm.ID
			return e, nil
		}
	}

	return nil, fmt.Errorf("Endpoint \"%s\" not found", endpoint)
}

func (self *Portainer) GetStackByName(name string) (*Stack, error) {
	if name == "" {
		return nil, fmt.Errorf("Stack name not defined")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/stacks", self.address), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", self.auth))

	rsp, err := self.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode > 200 {
		return nil, fmt.Errorf("Portainer API error: %s %s %s", req.Method, req.URL.String(), rsp.Status)
	}

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	var stacks []*Stack

	err = json.Unmarshal(data, &stacks)
	if err != nil {
		return nil, err
	}

	for _, stack := range stacks {
		if stack.Name == name {
			return stack, nil
		}
	}

	return nil, nil
}

func (self *Portainer) DeployStackFromGit(endpoint *Endpoint, name string, repo string, path string, user string, pass string, env ...*Env) error {
	args, err := json.Marshal(&struct {
		Name                        string `json:"Name"`
		SwarmID                     string `json:"SwarmID"`
		RepositoryURL               string `json:"RepositoryURL"`
		ComposeFilePathInRepository string `json:"ComposeFilePathInRepository"`
		RepositoryAuthentication    bool   `json:"RepositoryAuthentication"`
		RepositoryUsername          string `json:"RepositoryUsername"`
		RepositoryPassword          string `json:"RepositoryPassword"`
		Env                         []*Env `json:"Env"`
	}{
		Name:                        name,
		SwarmID:                     endpoint.SwarmID,
		RepositoryURL:               repo,
		ComposeFilePathInRepository: path,
		RepositoryAuthentication:    (len(user) > 0 && len(pass) > 0),
		RepositoryUsername:          user,
		RepositoryPassword:          pass,
		Env:                         env,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/stacks?type=1&method=repository&endpointId=%d", self.address, endpoint.Id), bytes.NewBuffer(args))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", self.auth))

	rsp, err := self.client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode > 200 {
		return fmt.Errorf("Portainer API error: %s %s %s", req.Method, req.URL.String(), rsp.Status)
	}

	return nil
}

func (self *Portainer) DeployStackFromString(endpoint *Endpoint, name string, config string, env ...*Env) error {
	args, err := json.Marshal(&struct {
		Name             string `json:"Name"`
		SwarmID          string `json:"SwarmID"`
		StackFileContent string `json:"StackFileContent"`
		Env              []*Env `json:"Env"`
	}{
		Name:             name,
		SwarmID:          endpoint.SwarmID,
		StackFileContent: config,
		Env:              env,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/stacks?type=1&method=string&endpointId=%d", self.address, endpoint.Id), bytes.NewBuffer(args))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", self.auth))

	rsp, err := self.client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode > 200 {
		return fmt.Errorf("Portainer API error: %s %s %s", req.Method, req.URL.String(), rsp.Status)
	}

	return nil
}

func (self *Portainer) DeployStackFromFile(endpoint *Endpoint, name string, path string, env ...*Env) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return self.DeployStackFromString(endpoint, name, string(data), env...)
}

func (self *Portainer) UpdateStackFromString(stack *Stack, config string, prune bool, env ...*Env) error {
	args, err := json.Marshal(&struct {
		StackFileContent string `json:"StackFileContent"`
		Prune            bool   `json:"Prune"`
		Env              []*Env `json:"Env"`
	}{
		StackFileContent: config,
		Prune:            prune,
		Env:              env,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/stacks/%d?endpointId=%d", self.address, stack.Id, stack.EndpointID), bytes.NewBuffer(args))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", self.auth))

	rsp, err := self.client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode > 200 {
		return fmt.Errorf("Portainer API error: %s %s %s", req.Method, req.URL.String(), rsp.Status)
	}

	return nil
}

func (self *Portainer) UpdateStackFromFile(stack *Stack, path string, prune bool, env ...*Env) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return self.UpdateStackFromString(stack, string(data), prune, env...)
}
